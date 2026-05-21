package parsers

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type EpisodeInfo struct {
	RawName string
	Series  string
	Season  int
	Episode int
}

var (
	// S03E03 / s03e03 / S3E3
	reSxxEyy = regexp.MustCompile(`(?i)\bS(\d{1,2})\s*(?:E|EP)\s*(\d{1,3})\b`)
	// 3x03
	reX = regexp.MustCompile(`(?i)\b(\d{1,2})\s*x\s*(\d{1,3})\b`)
	// Episode 12 / Episode 5: A Final Zone
	reEpisodeOnly = regexp.MustCompile(`(?i)^Episode\s+(\d{1,3})\b\s*:?.*`)

	// junk tokens that often appear AFTER the episode marker; we mainly cut before marker,
	// but this helps if someone gives you "Show Name S01E01 1080p" and you parse title wrong.
	reBrackets = regexp.MustCompile(`[\[\(].*?[\]\)]`)
)

// ParseEpisodeFromFilename tries to extract series name + season + episode from a typical scene-ish filename.
func ParseEpisodeFromFilename(name string) (*EpisodeInfo, bool) {
	base := filepath.Base(name)
	ext := filepath.Ext(base)
	base = strings.TrimSuffix(base, ext)

	// Remove bracketed segments (often release group / extra tags)
	base = reBrackets.ReplaceAllString(base, " ")

	// Normalize separators to spaces
	base = strings.NewReplacer(
		".", " ",
		"_", " ",
		"-", " ",
	).Replace(base)
	base = strings.Join(strings.Fields(base), " ") // collapse whitespace

	// Find episode marker
	var (
		loc         []int
		season      int
		ep          int
		episodeOnly bool
	)

	if m := reSxxEyy.FindStringSubmatchIndex(base); m != nil {
		// m[0]:m[1] is full match; groups at m[2]:m[3], m[4]:m[5]
		loc = m[0:2]
		s, _ := strconv.Atoi(base[m[2]:m[3]])
		e, _ := strconv.Atoi(base[m[4]:m[5]])
		season, ep = s, e
	} else if m := reX.FindStringSubmatchIndex(base); m != nil {
		loc = m[0:2]
		s, _ := strconv.Atoi(base[m[2]:m[3]])
		e, _ := strconv.Atoi(base[m[4]:m[5]])
		season, ep = s, e
	} else if m := reEpisodeOnly.FindStringSubmatchIndex(base); m != nil {
		e, _ := strconv.Atoi(base[m[2]:m[3]])
		ep = e
		episodeOnly = true
	} else {
		return nil, false
	}

	titlePart := ""
	if !episodeOnly {
		titlePart = strings.TrimSpace(base[:loc[0]])
		titlePart = strings.Trim(titlePart, " -")
	}

	// Very light cleanup
	titlePart = strings.Join(strings.Fields(titlePart), " ")

	if ep <= 0 || (!episodeOnly && (titlePart == "" || season <= 0)) {
		return nil, false
	}

	return &EpisodeInfo{
		RawName: name,
		Series:  titlePart,
		Season:  season,
		Episode: ep,
	}, true
}
