package steamapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Client struct {
	apiKey     string
	httpClient *http.Client
}

type PlayerSummary struct {
	SteamID     string `json:"steamid"`
	PersonaName string `json:"personaname"`
	AvatarFull  string `json:"avatarfull"`
}

type OwnedGame struct {
	AppID           int64  `json:"appid"`
	Name            string `json:"name"`
	PlaytimeForever int64  `json:"playtime_forever"`
}

type StoreGameDetails struct {
	Name                string   `json:"name"`
	DetailedDescription string   `json:"detailed_description"`
	ShortDescription    string   `json:"short_description"`
	Developers          []string `json:"developers"`
	Genres              []struct {
		Description string `json:"description"`
	} `json:"genres"`
	ReleaseDate struct {
		Date string `json:"date"`
	} `json:"release_date"`
	Background  string `json:"background"`
	HeaderImage string `json:"header_image"`
}

func NewClient(apiKey string, timeout time.Duration) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) GetPlayerSummary(ctx context.Context, steamID string) (*PlayerSummary, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("steam api key is not configured")
	}

	u := url.URL{
		Scheme: "https",
		Host:   "api.steampowered.com",
		Path:   "/ISteamUser/GetPlayerSummaries/v0002/",
	}
	q := u.Query()
	q.Set("key", c.apiKey)
	q.Set("steamids", steamID)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("steam player summary status: %d", res.StatusCode)
	}

	var payload struct {
		Response struct {
			Players []PlayerSummary `json:"players"`
		} `json:"response"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
	}

	if len(payload.Response.Players) == 0 {
		return nil, nil
	}

	return &payload.Response.Players[0], nil
}

func (c *Client) GetOwnedGames(ctx context.Context, steamID string) ([]OwnedGame, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("steam api key is not configured")
	}


	u := url.URL{
		Scheme: "https",
		Host:   "api.steampowered.com",
		Path:   "/IPlayerService/GetOwnedGames/v1/",
	}
	q := u.Query()
	q.Set("key", c.apiKey)
	q.Set("steamid", steamID)
	q.Set("include_appinfo", "1")
	q.Set("include_played_free_games", "1")
	q.Set("format", "json")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("steam owned games status: %d", res.StatusCode)
	}

	var payload struct {
		Response struct {
			Games []OwnedGame `json:"games"`
		} `json:"response"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return payload.Response.Games, nil
}

func (c *Client) GetStoreDetails(ctx context.Context, appID string) (*StoreGameDetails, error) {
	u := url.URL{
		Scheme: "https",
		Host:   "store.steampowered.com",
		Path:   "/api/appdetails",
	}
	q := u.Query()
	q.Set("appids", appID)
	q.Set("l", "brazilian")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("steam store details status: %d", res.StatusCode)
	}

	var payload map[string]struct {
		Success bool             `json:"success"`
		Data    StoreGameDetails `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
	}

	entry, ok := payload[appID]
	if !ok || !entry.Success {
		return nil, nil
	}

	return &entry.Data, nil
}

func AppIDToString(appID int64) string {
	return strconv.FormatInt(appID, 10)
}
