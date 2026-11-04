package modules

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"
)

type platformDef struct {
	Name  string
	Check func(username string) PlatformCheck

	// Reliable is false for platforms whose "not found" page is
	// indistinguishable from a real profile without executing JavaScript or
	// an authenticated API (X and Reddit both serve an identical, generic
	// client-rendered shell either way -- confirmed by diffing real vs fake
	// profile pages byte-for-byte). They're still included in the main
	// single-username check (with the existing caveat), but excluded from
	// the variations sweep, which only surfaces "found" hits with no
	// surrounding context -- a guess there would be actively misleading
	// rather than just a documented caveat.
	Reliable bool
}

// platformDefs lists every platform CheckUsername probes. Most sites return
// a clean 404 for an unknown user, so statusCodeCheck (found == HTTP 200)
// covers them; a few need a smarter check because they always return 200.
var platformDefs = []platformDef{
	{"GitHub", statusCodeCheck("https://github.com/%s"), true},
	{"X", statusCodeCheck("https://x.com/%s"), false},
	{"Reddit", statusCodeCheck("https://reddit.com/user/%s"), false},
	{"GitLab", statusCodeCheck("https://gitlab.com/%s"), true},
	{"Keybase", keybaseCheck, true},
	{"Hacker News", hackerNewsCheck, true},
	{"Telegram", telegramCheck, true},
}

// CheckUsername probes every platform in platformDefs concurrently and
// reports whether a profile appears to exist on each.
//
// Note: some platforms render profiles client-side and return HTTP 200 even
// for pages that don't exist, so results should be treated as a lead, not
// proof.
func CheckUsername(username string, result *UsernameResult) {
	result.Platforms = checkAllPlatforms(username, platformDefs)
}

// reliablePlatformDefs is platformDefs filtered down to the platforms whose
// checks can actually tell "found" from "not found".
func reliablePlatformDefs() []platformDef {
	var defs []platformDef
	for _, p := range platformDefs {
		if p.Reliable {
			defs = append(defs, p)
		}
	}
	return defs
}

// checkAllPlatforms runs every check in defs concurrently for one username
// and returns the results sorted by platform name.
func checkAllPlatforms(username string, defs []platformDef) []PlatformCheck {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var checks []PlatformCheck

	for _, p := range defs {
		wg.Add(1)
		go func(p platformDef) {
			defer wg.Done()
			check := p.Check(username)
			check.Platform = p.Name

			mu.Lock()
			checks = append(checks, check)
			mu.Unlock()
		}(p)
	}

	wg.Wait()

	sort.Slice(checks, func(i, j int) bool {
		return checks[i].Platform < checks[j].Platform
	})
	return checks
}

// profileURL fills a "%s" URL template with the username, escaped so that
// characters like "/", "?" or "#" can't point the check at a different page
// (e.g. "torvalds/linux" would otherwise hit a repo and come back as found).
func profileURL(urlTemplate, username string) string {
	return fmt.Sprintf(urlTemplate, url.PathEscape(username))
}

// statusCodeCheck builds a checker that treats an HTTP 200 as "found". Good
// enough for sites that return a real 404 for unknown users.
func statusCodeCheck(urlTemplate string) func(string) PlatformCheck {
	return func(username string) PlatformCheck {
		u := profileURL(urlTemplate, username)
		check := PlatformCheck{URL: u}

		resp, err := fetch(u)
		if err != nil {
			check.Error = err.Error()
			return check
		}
		defer resp.Body.Close()

		check.Found = resp.StatusCode == http.StatusOK
		return check
	}
}

type keybaseResponse struct {
	Status struct {
		Code int `json:"code"`
	} `json:"status"`
}

// keybaseCheck uses Keybase's real lookup API instead of guessing from a
// status code, so it's the most reliable check here.
func keybaseCheck(username string) PlatformCheck {
	check := PlatformCheck{URL: profileURL("https://keybase.io/%s", username)}

	resp, err := fetch("https://keybase.io/_/api/1.0/user/lookup.json?username=" + url.QueryEscape(username))
	if err != nil {
		check.Error = err.Error()
		return check
	}
	defer resp.Body.Close()

	var data keybaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		check.Error = err.Error()
		return check
	}

	check.Found = data.Status.Code == 0
	return check
}

// hackerNewsCheck always gets a 200, so existence has to be read from the
// page content. HN renders "No such user." for an unknown id -- but under
// load it can also serve a generic block/rate-limit page that has neither
// that text nor any real profile content. Treating "found" as "anything
// that isn't the 'No such user.' page" is what caused false positives
// during testing, so this instead requires a positive marker ("karma:",
// which every real profile page shows) before calling it found, and reports
// anything else as inconclusive rather than guessing.
func hackerNewsCheck(username string) PlatformCheck {
	u := "https://news.ycombinator.com/user?id=" + url.QueryEscape(username)
	check := PlatformCheck{URL: u}

	throttle("Hacker News", 500*time.Millisecond)
	body, err := fetchBody(u)
	if err != nil {
		check.Error = err.Error()
		return check
	}

	switch {
	case strings.Contains(body, "No such user."):
		check.Found = false
	case strings.Contains(body, "karma:"):
		check.Found = true
	default:
		check.Error = "could not determine (Hacker News returned an unexpected page, possibly rate-limiting)"
	}
	return check
}

// telegramCheck also always gets a 200 for both real and unclaimed
// usernames. Telegram falls back to an og:title of "Telegram: Contact
// @username" when there's no real profile behind it, so that's the marker
// used to tell the two apart.
func telegramCheck(username string) PlatformCheck {
	u := profileURL("https://t.me/%s", username)
	check := PlatformCheck{URL: u}

	throttle("Telegram", 300*time.Millisecond)
	body, err := fetchBody(u)
	if err != nil {
		check.Error = err.Error()
		return check
	}

	check.Found = !strings.Contains(body, "Telegram: Contact @")
	return check
}
