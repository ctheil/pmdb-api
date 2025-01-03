package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
)

type OAuth struct {
	CFG *CFG
	TS  *TokenService
}

func NewOAuth() (*OAuth, error) {
	cfg, err := GetCFG()
	if err != nil {
		return &OAuth{}, err
	}
	return &OAuth{CFG: &cfg, TS: NewTokenService()}, nil
}

type CFG struct {
	Client_ID        string `json:"clientId"`
	Client_Secret    string `json:"clientSecret"`
	Auth_URL         string `json:"authUrl"`
	Token_URL        string `json:"tokenUrl"`
	Redirect_URL     string `json:"redirectUrl"`
	Client_URL       string `json:"clientUrl"`
	Token_Secret     string `json:"tokenSecret"`
	Token_Expiration int    `json:"tokenExpiration"`
	Post_URL         string `json:"postUrl"`
}

func GetCFG() (CFG, error) {
	env_cfg := map[string]string{
		"Client_ID":     "GG_OAUTH_CLIENT_ID",
		"Client_Secret": "GG_OAUTH_CLIENT_SECRET",
		"Token_Secret":  "JWT_SECRET",
	}
	cfg := &CFG{
		Auth_URL:         "https://accounts.google.com/o/oauth2/v2/auth",
		Token_URL:        "https://oauth2.googleapis.com/token",
		Redirect_URL:     "http://localhost:5173/auth/callback",
		Client_URL:       "http://localhost:5173",
		Token_Expiration: int(time.Millisecond * 36000),
		Post_URL:         "https://jsonplaceholder.typicode.com/posts",
	}

	v := reflect.ValueOf(cfg).Elem()

	for k, envVar := range env_cfg {
		envValue, ok := os.LookupEnv(envVar)
		if !ok {
			return *cfg, fmt.Errorf("environment variable '%s' not set", v)
		}
		field := v.FieldByName(k)
		if !field.IsValid() {
			return *cfg, fmt.Errorf("field '%s' does not exist in OAuth_CFG", k)
		}
		if field.CanSet() {
			field.SetString(envValue)
		} else {
			return *cfg, fmt.Errorf("field '%s' cannot be set", k)
		}
	}
	return *cfg, nil
}

type OAuthParams struct {
	client_id     string
	redirect_uri  string
	response_type string
	scope         string
	access_type   string
	state         string
	prompt        string
}

func (a *OAuth) GetOAuthParams() string {
	params := map[string]string{
		"scope":                  "openid profile email",
		"access_type":            "offline",
		"include_granted_scopes": "true",
		"response_type":          "code",
		"state":                  "standard_oauth",
		"redirect_uri":           a.CFG.Redirect_URL,
		"client_id":              a.CFG.Client_ID,
	}
	// "prompt":        "consent",
	vals := url.Values{}
	for k, v := range params {
		vals.Add(k, v)
	}
	return vals.Encode()
}

func (a *OAuth) GetOAuthTokenParams(code string) string {
	params := map[string]string{
		"code":          code,
		"grant_type":    "authorization_code",
		"client_secret": a.CFG.Client_Secret,
		"redirect_uri":  a.CFG.Redirect_URL,
		"client_id":     a.CFG.Client_ID,
	}
	vals := url.Values{}
	for k, v := range params {
		vals.Add(k, v)
	}
	return vals.Encode()
}

type OAuthResp struct {
	ID_Token string `json:"id_token"`
}

func (a *OAuth) FetchAuthToken(code string) (*OAuthUserClaims, error) {
	req_url := fmt.Sprintf("%s?%s", a.CFG.Token_URL, a.GetOAuthTokenParams(code))
	req, err := http.NewRequest("POST", req_url, nil)
	if err != nil {
		return nil, err
	}

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch oauth data: %s", resp.Status)
	}

	oauth_resp := OAuthResp{}

	if err := json.NewDecoder(resp.Body).Decode(&oauth_resp); err != nil {
		fmt.Printf("Error decoding oauth into json: %e\n\n", err)
		return nil, err
	}

	user := a.TS.ParseOAuthUserToken(oauth_resp.ID_Token)

	return user, nil
}

func (a *OAuth) VerifyToken(c *gin.Context) (newToken string, ok bool) {
	accessToken, err := c.Cookie("access_token")
	if err != nil {
		// fmt.Errorf("no access token found in request")
		return "", false
	}
	user_claims := a.TS.ParseOAuthUserToken(accessToken)
	isValid := user_claims.ExpiresAt.After(time.Now())
	token, err := a.TS.NewOAuthUserToken(user_claims)
	if err != nil {
		return "", false
	}
	fmt.Printf("token: %s, isValid: %t\n", token, isValid)
	return token, isValid
}
