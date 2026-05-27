package steam

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type OpenIDValidator struct {
	httpClient *http.Client
}

func NewOpenIDValidator() *OpenIDValidator {
	return &OpenIDValidator{
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// Validate performs the check_authentication request against the Steam OpenID endpoint.
func (v *OpenIDValidator) Validate(ctx context.Context, params map[string]string) (string, error) {
	if params["openid.mode"] != "id_res" {
		return "", errors.New("invalid openid.mode")
	}

	// Prepare check_authentication payload
	payload := url.Values{}
	for key, val := range params {
		payload.Set(key, val)
	}
	payload.Set("openid.mode", "check_authentication")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://steamcommunity.com/openid/login", strings.NewReader(payload.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create check_authentication request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("check_authentication request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("steam openid returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read check_authentication response: %w", err)
	}

	if !strings.Contains(string(body), "is_valid:true") {
		return "", errors.New("steam openid validation failed: signature invalid")
	}

	claimedID := params["openid.claimed_id"]
	if claimedID == "" {
		return "", errors.New("missing openid.claimed_id")
	}

	// Extract SteamID64 from the claimed_id
	// Format is typically: https://steamcommunity.com/openid/id/76561198000000000
	re := regexp.MustCompile(`^https?://steamcommunity\.com/openid/id/(\d+)$`)
	matches := re.FindStringSubmatch(claimedID)
	if len(matches) < 2 {
		return "", errors.New("invalid openid.claimed_id format")
	}

	return matches[1], nil
}
