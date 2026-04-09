package parsers

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type SeasonType string

const (
	SeasonTypeRegular  SeasonType = "regular"
	SeasonTypeSpecials SeasonType = "specials"
)

type SeasonInfo struct {
	RawName string
	Season  int
	Type    SeasonType
}

var (
	reSeasonEpisodeMarker = regexp.MustCompile(`(?i)(\bS\d{1,2}\s*E\d{1,3}\b|\b\d{1,2}\s*x\s*\d{1,3}\b)`)
	reSeasonSpecials      = regexp.MustCompile(`(?i)\b(specials?|sp)\b`)
	reSeasonWord          = regexp.MustCompile(`(?i)\bseason\s*(\d{1,2})\b`)
	reSeasonShort         = regexp.MustCompile(`(?i)\bs\s*(\d{1,2})\b`)
	reSeasonNumber        = regexp.MustCompile(`^(\d{1,2})$`)
)

// ParseSeasonFromName extracts a season number from a season directory name.
func ParseSeasonFromName(name string) (*SeasonInfo, bool) {
	base := filepath.Base(strings.TrimSpace(name))
	if base == "" || base == "." || base == string(filepath.Separator) {
		return nil, false
	}

	base = strings.NewReplacer(
		".", " ",
		"_", " ",
		"-", " ",
	).Replace(base)
	base = strings.Join(strings.Fields(base), " ")

	if base == "" || reSeasonEpisodeMarker.MatchString(base) {
		return nil, false
	}

	season := -1
	if m := reSeasonWord.FindStringSubmatch(base); m != nil {
		n, err := strconv.Atoi(m[1])
		if err != nil {
			return nil, false
		}
		season = n
	} else if m := reSeasonShort.FindStringSubmatch(base); m != nil {
		n, err := strconv.Atoi(m[1])
		if err != nil {
			return nil, false
		}
		season = n
	} else if m := reSeasonNumber.FindStringSubmatch(base); m != nil {
		n, err := strconv.Atoi(m[1])
		if err != nil {
			return nil, false
		}
		season = n
	} else if reSeasonSpecials.MatchString(base) {
		return &SeasonInfo{
			RawName: name,
			Season:  0,
			Type:    SeasonTypeSpecials,
		}, true
	} else {
		return nil, false
	}

	if season < 0 || season > 99 {
		return nil, false
	}

	if season == 0 {
		return &SeasonInfo{
			RawName: name,
			Season:  0,
			Type:    SeasonTypeSpecials,
		}, true
	}

	return &SeasonInfo{
		RawName: name,
		Season:  season,
		Type:    SeasonTypeRegular,
	}, true
}
