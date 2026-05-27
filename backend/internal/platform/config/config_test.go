package config

import (
	"os"
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "memory repository invalid",
			cfg: Config{
				LibraryRepository: "memory",
			},
			wantErr: true,
		},
		{
			name: "postgres repository requires dsn",
			cfg: Config{
				LibraryRepository: "postgres",
				PostgresDSN:       "",
			},
			wantErr: true,
		},
		{
			name: "postgres repository with dsn valid",
			cfg: Config{
				LibraryRepository: "postgres",
				PostgresDSN:       "postgres://example",
			},
			wantErr: false,
		},
		{
			name: "unknown repository invalid",
			cfg: Config{
				LibraryRepository: "mongo",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

func TestLoadEnv(t *testing.T) {
	// Create a temporary .env file
	content := `
# A comment line
TEST_KEY_1=value1
TEST_KEY_2="value2"
TEST_KEY_3 = 'value3'

# Another comment
  TEST_KEY_4 = value4  
`
	tmpfile, err := os.CreateTemp("", ".env")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Clean up environment variables
	keys := []string{"TEST_KEY_1", "TEST_KEY_2", "TEST_KEY_3", "TEST_KEY_4"}
	for _, k := range keys {
		os.Unsetenv(k)
	}
	defer func() {
		for _, k := range keys {
			os.Unsetenv(k)
		}
	}()

	loadEnv(tmpfile.Name())

	expected := map[string]string{
		"TEST_KEY_1": "value1",
		"TEST_KEY_2": "value2",
		"TEST_KEY_3": "value3",
		"TEST_KEY_4": "value4",
	}

	for k, want := range expected {
		got := os.Getenv(k)
		if got != want {
			t.Errorf("expected %s to be %q, got %q", k, want, got)
		}
	}
}

