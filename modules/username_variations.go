package modules

import (
	"sort"
	"sync"
)

// variationSuffixes and variationPrefixes are common patterns people use for
// alt accounts, backup accounts, or slightly-different handles when their
// exact username is already taken.
var (
	variationSuffixes = []string{"1", "123", "_", "01", "official"}
	variationPrefixes = []string{"the_", "real_", "its_"}
)

// GenerateVariations returns common variations of a username, e.g. "johndoe"
// becomes ["johndoe1", "johndoe123", ..., "real_johndoe", ...].
func GenerateVariations(username string) []string {
	seen := map[string]bool{username: true}
	var variants []string

	add := func(v string) {
		if !seen[v] {
			seen[v] = true
			variants = append(variants, v)
		}
	}

	for _, suffix := range variationSuffixes {
		add(username + suffix)
	}
	for _, prefix := range variationPrefixes {
		add(prefix + username)
	}

	return variants
}

// CheckUsernameVariations generates common variations of the username and
// checks every one of them concurrently against the reliable platforms only
// (see platformDef.Reliable) -- this only surfaces "found" hits with no
// other context, so a platform that can't tell found from not-found would
// just produce noise here. Only matches are kept; with this many
// combinations, showing every "not found" would be noise too.
func CheckUsernameVariations(username string, result *UsernameResult) {
	variants := GenerateVariations(username)
	result.Variations = variants
	defs := reliablePlatformDefs()

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, variant := range variants {
		wg.Add(1)
		go func(variant string) {
			defer wg.Done()
			for _, check := range checkAllPlatforms(variant, defs) {
				if !check.Found {
					continue
				}
				mu.Lock()
				result.VariationMatches = append(result.VariationMatches, VariationMatch{
					Username: variant,
					Platform: check.Platform,
					URL:      check.URL,
				})
				mu.Unlock()
			}
		}(variant)
	}

	wg.Wait()

	sort.Slice(result.VariationMatches, func(i, j int) bool {
		a, b := result.VariationMatches[i], result.VariationMatches[j]
		if a.Username != b.Username {
			return a.Username < b.Username
		}
		return a.Platform < b.Platform
	})
}
