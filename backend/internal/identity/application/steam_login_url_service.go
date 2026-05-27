package application

import "net/url"

type SteamLoginURLConfig struct {
	FrontendBaseURL string
	OpenIDRealm     string
	OpenIDReturnURL string
}

type SteamLoginURLService struct {
	cfg SteamLoginURLConfig
}

func NewSteamLoginURLService(cfg SteamLoginURLConfig) *SteamLoginURLService {
	return &SteamLoginURLService{cfg: cfg}
}

func (s *SteamLoginURLService) Build() string {
	params := url.Values{}
	params.Set("openid.mode", "checkid_setup")
	params.Set("openid.ns", "http://specs.openid.net/auth/2.0")
	params.Set("openid.identity", "http://specs.openid.net/auth/2.0/identifier_select")
	params.Set("openid.claimed_id", "http://specs.openid.net/auth/2.0/identifier_select")
	params.Set("openid.return_to", s.cfg.OpenIDReturnURL)
	params.Set("openid.realm", s.cfg.OpenIDRealm)

	return "https://steamcommunity.com/openid/login?" + params.Encode()
}
