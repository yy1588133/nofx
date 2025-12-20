package coinank

import (
	"os"
	"testing"
)

var TestApikey = os.Getenv("COINANK_API_KEY")

func requireTestApikey(t *testing.T) string {
	t.Helper()
	if TestApikey == "" {
		t.Skip("COINANK_API_KEY is not set")
	}
	return TestApikey
}
