package parsers

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type ShowInfo struct {
	RawName string
	Series  string
	Year    *int
}

var (
	reShowTrailingYear = regexp.MustCompile(`^(.+?)\s*(?:\((\d{4})\)|\[(\d{4})\]|\{(\d{4})\}|(\d{4}))$`)
	reYearOnly         = regexp.MustCompile(`^(?:\((\d{4})\)|\[(\d{4})\]|\{(\d{4})\]|(\d{4}))$`)
)

// ParseShowFromName extracts a show/series name and optional trailing year.
func ParseShowFromName(name string) (*ShowInfo, bool) {
	base := filepath.Base(strings.TrimSpace(name))
	if base == "" || base == "." || base == string(filepath.Separator) {
		return nil, false
	}

	base = strings.NewReplacer(
		".", " ",
		"_", " ",
	).Replace(base)
	base = strings.Join(strings.Fields(base), " ")

	if base == "" || reYearOnly.MatchString(base) {
		return nil, false
	}

	series := base
	var year *int

	if m := reShowTrailingYear.FindStringSubmatch(base); m != nil {
		series = strings.TrimSpace(m[1])
		for i := 2; i <= 5; i++ {
			if m[i] == "" {
				continue
			}
			y, err := strconv.Atoi(m[i])
			if err != nil {
				return nil, false
			}
			year = &y
			break
		}
	}

	if series == "" {
		return nil, false
	}

	return &ShowInfo{
		RawName: name,
		Series:  series,
		Year:    year,
	}, true
}
