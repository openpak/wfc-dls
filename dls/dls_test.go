package dls

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"openpak/wfc-common/common"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func dwcForm(t *testing.T, values map[string]string) string {
	t.Helper()
	form := url.Values{}
	for key, value := range values {
		form.Set(key, common.Base64DwcEncoding.EncodeToString([]byte(value)))
	}
	return form.Encode()
}

func authRequest(t *testing.T, form string) *http.Request {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "http://nas.nintendowifi.net/ac", strings.NewReader(form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Host", "nas.nintendowifi.net")
	if err := req.ParseForm(); err != nil {
		t.Fatal(err)
	}
	return req
}

func dlsRequest(t *testing.T, host, action, rhgamecd string) *http.Request {
	t.Helper()
	form := map[string]string{"action": action}
	if rhgamecd != "" {
		form["rhgamecd"] = rhgamecd
	}
	req := httptest.NewRequest(http.MethodPost, "http://"+host+"/download",
		strings.NewReader(dwcForm(t, form)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Host", host)
	return req
}

func TestDLS1CountsDLCFiles(t *testing.T) {
	oldDir := dlcDir
	t.Cleanup(func() { dlcDir = oldDir })
	dlcDir = t.TempDir()

	for game, files := range map[string]int{"RMCE": 2, "ABCD": 0} {
		dir := filepath.Join(dlcDir, game)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		for i := 0; i < files; i++ {
			if err := os.WriteFile(filepath.Join(dir, fmt.Sprintf("dlc%d.bin", i)), []byte{1}, 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}

	rec := httptest.NewRecorder()
	Handler(rec, dlsRequest(t, "dls1.nintendowifi.net", "count", "RMCE"))
	if rec.Code != http.StatusOK || rec.Body.String() != "2\x00" {
		t.Fatalf("count = %d %q, want 200 \"2\\x00\"", rec.Code, rec.Body.String())
	}

	// A game with no DLC on the server still gets an answer, never an error:
	// the console treats a failed count as failed DLC.
	rec = httptest.NewRecorder()
	Handler(rec, dlsRequest(t, "dls1.nintendowifi.net", "count", "NOPE"))
	if rec.Code != http.StatusOK || rec.Body.String() != "0\x00" {
		t.Fatalf("missing game = %d %q, want 200 \"0\\x00\"", rec.Code, rec.Body.String())
	}
}

// TestDLS1RejectsBadGameCodes pins the input validation: the game code lands
// in a filesystem path, so anything non-alphanumeric must be refused.
func TestDLS1RejectsBadGameCodes(t *testing.T) {
	for _, bad := range []string{"rmce", "RMC", "RMCE1", "RMC-E"} {
		rec := httptest.NewRecorder()
		Handler(rec, dlsRequest(t, "dls1.nintendowifi.net", "count", bad))
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("rhgamecd %q: code = %d, want 400", bad, rec.Code)
		}
	}
}
