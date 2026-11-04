package modules

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProfileURLEscapesUsername(t *testing.T) {
	cases := map[string]string{
		"johndoe":        "https://github.com/johndoe",
		"john.doe_99":    "https://github.com/john.doe_99",
		"torvalds/linux": "https://github.com/torvalds%2Flinux",
		"foo?tab=repos":  "https://github.com/foo%3Ftab=repos",
		"foo#bar":        "https://github.com/foo%23bar",
		"a b":            "https://github.com/a%20b",
	}
	for in, want := range cases {
		if got := profileURL("https://github.com/%s", in); got != want {
			t.Errorf("profileURL(%q) = %q, want %q", in, got, want)
		}
	}
}

// A username containing "/" must not reach a different page on the site and
// get reported as found.
func TestStatusCodeCheckDoesNotFollowSlashIntoAnotherPage(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only a real profile path and a repo-style path exist on this fake site.
		switch r.URL.EscapedPath() {
		case "/torvalds", "/torvalds/linux":
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	check := statusCodeCheck(srv.URL + "/%s")

	if got := check("torvalds"); !got.Found || got.Error != "" {
		t.Errorf("real user: got %+v, want found", got)
	}
	if got := check("torvalds/linux"); got.Found {
		t.Errorf("username with a slash was reported as found: %+v", got)
	}
	if got := check("torvalds?tab=repos"); got.Found {
		t.Errorf("username with a query string was reported as found: %+v", got)
	}
}
