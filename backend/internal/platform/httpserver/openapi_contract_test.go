package httpserver

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestOpenAPIContractContainsCoreRoutesAndStatuses(t *testing.T) {
	specPath := filepath.Clean("../../../api/openapi/openapi.yaml")
	content, err := os.ReadFile(specPath)
	if err != nil {
		t.Fatalf("read openapi spec failed: %v", err)
	}

	spec := string(content)
	requiredSnippets := []string{
		"openapi: 3.1.0",
		"/v1/auth/steam/exchange:",
		"/v1/auth/me:",
		"/v1/auth/logout:",
		"/v1/users/{userId}/games:",
		"/v1/users/{userId}/games/{gameId}:",
		"/v1/integrations/steam/users/{steamId}/games:",
		"/v1/catalog/steam/apps/{appId}:",
		"/readyz:",
		"'503':",
		"'504':",
		"'502':",
		"'400':",
	}

	for _, snippet := range requiredSnippets {
		if !strings.Contains(spec, snippet) {
			t.Fatalf("openapi spec missing required snippet: %s", snippet)
		}
	}
}
