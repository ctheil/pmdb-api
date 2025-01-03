package services

import (
	"fmt"
	"io"
	"strings"

	"github.com/ctheil/pmdb-api/internal/config"
	"github.com/ctheil/pmdb-api/internal/model"
	"github.com/ctheil/pmdb-api/pkg/utils"
)

type TitleService struct {
	Config  *config.Config
	BaseURL string
}

func NewTitleService() (*TitleService, error) {
	cfg, err := config.FetchConfig()
	if err != nil {
		return nil, err
	}

	return &TitleService{Config: cfg, BaseURL: "https://api.themoviedb.org"}, nil
}

func (ts *TitleService) GetDomain(domain string) (string, error) {
	switch domain {
	case "movie":
		return domain, nil
	case "tv":
		return domain, nil
	case "person":
		return domain, nil
	default:
		return "", fmt.Errorf("invalid domain: %s", domain)
	}
}

func (ts *TitleService) ExtractIncludes(includes string) string {
	if includes == "" {
		return ""
	}
	chars := strings.Split(includes, "")
	out := ""

	for _, c := range chars {
		switch c {
		case "c":
			out += "credits%2C"
		case "v":
			out += "videos%2C"
		case "p":
			out += "watch%2Fproviders%2C"
		}
	}
	// trim last '%2C' from end
	out = out[:len(out)-3]
	return out
}

func (ts *TitleService) TMDBRequest(endpoint, method string, reader io.Reader, out interface{}) error {
	req, err := utils.BuildReq(ts.BaseURL+endpoint, method, reader)
	if err != nil {
		return fmt.Errorf("error building request: %e", err)
	}
	resp, err := utils.ExecReq(req)
	if err != nil {
		return fmt.Errorf("error executing request: %e", err)
	}

	if err := utils.ReqToJSON(resp.Body, out); err != nil {
		return fmt.Errorf("error decoding body: %e", err)
	}
	return nil
}

func (ts *TitleService) BuildPosterPaths(titles []model.Title, size int) {
	len := len(ts.Config.Images.PosterSizes)
	if size > len {
		size = len - 1
	}
	for i, title := range titles {
		titles[i].PosterPath = ts.Config.Images.BaseUrl + ts.Config.Images.PosterSizes[size] + title.PosterPath
	}
}

func (ts *TitleService) BuildBackdropPaths(titles []model.Title, size int) {
	len := len(ts.Config.Images.BackdropSizes)
	if size > len {
		size = len - 1
	}
	for i, title := range titles {
		path := ts.Config.Images.BaseUrl + ts.Config.Images.BackdropSizes[size] + title.BackdropPath
		titles[i].BackdropPath = path
	}
}

func (ts *TitleService) BuildLogoPaths(companies []model.ProductionCompany, size int) {
	len := len(ts.Config.Images.LogoSizes)
	if size > len {
		size = len - 1
	}
	for i, c := range companies {
		if c.LogoPath == "" {
			continue
		}
		path := ts.Config.Images.BaseUrl + ts.Config.Images.LogoSizes[size] + c.LogoPath
		companies[i].LogoPath = path
	}
}

func (ts *TitleService) BuiltProfilePaths(credits []model.Credit, size int) {
	len := len(ts.Config.Images.ProfileSizes)
	if size > len {
		size = len - 1
	}
	for i, c := range credits {
		if c.ProfilePath == "" {
			credits[i].ProfilePath = ""
			continue
		}
		path := ts.Config.Images.BaseUrl + ts.Config.Images.ProfileSizes[size] + c.ProfilePath
		credits[i].ProfilePath = path
	}
}

func (ts *TitleService) BuildProviderLogoPaths(p model.WatchProviders, size int) {
	lenSizes := len(ts.Config.Images.LogoSizes)
	if size > lenSizes {
		size = lenSizes - 1
	}

	for regionCode, region := range p.Results {
		if len(region.Buy) > 0 {
			for i, provider := range region.Buy {
				path := ts.Config.Images.BaseUrl + ts.Config.Images.LogoSizes[size] + provider.LogoPath
				p.Results[regionCode].Buy[i].LogoPath = path
			}
		}
		if len(region.Rent) > 0 {
			for i, provider := range region.Rent {
				path := ts.Config.Images.BaseUrl + ts.Config.Images.LogoSizes[size] + provider.LogoPath
				p.Results[regionCode].Rent[i].LogoPath = path
			}
		}
		if len(region.Flatrate) > 0 {
			for i, provider := range region.Flatrate {
				path := ts.Config.Images.BaseUrl + ts.Config.Images.LogoSizes[size] + provider.LogoPath
				p.Results[regionCode].Flatrate[i].LogoPath = path
			}
		}
	}
}
