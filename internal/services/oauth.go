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
}

func NewOAuth() (*OAuth, error) {
	cfg, err := GetCFG()
	if err != nil {
		return &OAuth{}, err
	}
	return &OAuth{&cfg}, nil
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

// const getTokenParams = (code) =>
//
//	queryString.stringify({
//	  client_id: config.clientId,
//	  client_secret: config.clientSecret,
//	  code,
//	  grant_type: 'authorization_code',
//	  redirect_uri: config.redirectUrl,
//	})
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

type OAuthResp struct {
	ID_Token string `json:"id_token"`
}

func (a *OAuth) FetchAuth() (*OAuthResp, error) {
	req_url := fmt.Sprintf("%s?%s", a.CFG.Auth_URL, a.GetOAuthParams())
	req, err := http.NewRequest("POST", req_url, nil)
	if err != nil {
		return nil, err
	}

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	fmt.Printf("\n\n    OAuth fetch request: %+v\n\n    Body: %s\n\n    Original Endpoint: %s\n\n", resp, resp.Body, req_url)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch oauth data: %s", resp.Status)
	}

	oauth_resp := OAuthResp{}
	result := map[string]string{}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("Error decoding oauth into json: %e\n\n", err)
		return nil, err
	}
	fmt.Printf("JSON result: %+v", result)

	return &oauth_resp, nil
}

func (a *OAuth) VerifyToken(c *gin.Context) (newToken string, error error) {
	accessToken, err := c.Cookie("access_token")
	if err != nil {
		return "", fmt.Errorf("no access token found in request")
	}
	user_claims := ParseOAuthUserToken(accessToken)
	return NewOAuthUserToken(user_claims.UserData, a.CFG.Token_Secret)
}
