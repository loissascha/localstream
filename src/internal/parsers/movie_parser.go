package parsers

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type MovieInfo struct {
	RawName string
	Title   string
	Year    *int
}

var (
	reMovieYearToken = regexp.MustCompile(`^(?:19|20)\d{2}$`)
	reMovieBrackets  = regexp.MustCompile(`[\(\[\{][^\)\]\}]+[\)\]\}]`)
	reMovieQuality   = regexp.MustCompile(`^\d{3,4}p$`)

	movieReleaseJunk = map[string]struct{}{
		"2160p":    {},
		"1080p":    {},
		"720p":     {},
		"480p":     {},
		"4k":       {},
		"8k":       {},
		"web":      {},
		"webdl":    {},
		"webrip":   {},
		"bluray":   {},
		"bdrip":    {},
		"brrip":    {},
		"dvdrip":   {},
		"hdrip":    {},
		"remux":    {},
		"x264":     {},
		"x265":     {},
		"h264":     {},
		"h265":     {},
		"hevc":     {},
		"av1":      {},
		"aac":      {},
		"ac3":      {},
		"dts":      {},
		"atmos":    {},
		"hdr":      {},
		"proper":   {},
		"repack":   {},
		"extended": {},
		"unrated":  {},
		"internal": {},
		"limited":  {},
		"multi":    {},
		"subbed":   {},
		"dubbed":   {},
		"dl":       {},
	}
)

// ParseMovieFromFilename extracts a movie title and optional year from a filename.
func ParseMovieFromFilename(name string) (*MovieInfo, bool) {
	base := filepath.Base(strings.TrimSpace(name))
	if base == "" || base == "." || base == string(filepath.Separator) {
		return nil, false
	}

	base = strings.TrimSuffix(base, filepath.Ext(base))
	if base == "" {
		return nil, false
	}

	base = reMovieBrackets.ReplaceAllStringFunc(base, func(seg string) string {
		inner := strings.TrimSpace(seg[1 : len(seg)-1])
		if reMovieYearToken.MatchString(inner) {
			return " " + inner + " "
		}
		return " "
	})

	base = strings.NewReplacer(
		".", " ",
		"_", " ",
		"-", " ",
	).Replace(base)
	base = strings.Join(strings.Fields(base), " ")
	if base == "" {
		return nil, false
	}

	tokens := strings.Fields(base)
	if len(tokens) == 0 {
		return nil, false
	}

	cutIndex := len(tokens)
	for i := 0; i < len(tokens); i++ {
		if isMovieMetadataToken(tokens[i]) {
			cutIndex = i
			break
		}
	}

	titleTokens := append([]string(nil), tokens[:cutIndex]...)
	if len(titleTokens) == 0 {
		return nil, false
	}

	var year *int

	lastYearToken := -1
	for i := len(titleTokens) - 1; i >= 0; i-- {
		t := normalizeMovieToken(titleTokens[i])
		if !reMovieYearToken.MatchString(t) {
			continue
		}
		if len(titleTokens) <= 1 {
			break
		}
		lastYearToken = i
		if year == nil {
			y, err := strconv.Atoi(t)
			if err != nil {
				return nil, false
			}
			year = &y
		}
		break
	}

	if lastYearToken >= 0 && len(titleTokens) > 1 {
		titleTokens = append(titleTokens[:lastYearToken], titleTokens[lastYearToken+1:]...)
	}

	title := strings.Join(titleTokens, " ")
	title = strings.TrimSpace(strings.Trim(title, "-_."))
	title = strings.Join(strings.Fields(title), " ")
	if title == "" {
		return nil, false
	}

	return &MovieInfo{
		RawName: name,
		Title:   title,
		Year:    year,
	}, true
}

func normalizeMovieToken(token string) string {
	token = strings.ToLower(token)
	token = strings.Trim(token, "[](){}.,")
	return token
}

func isMovieMetadataToken(token string) bool {
	normalized := normalizeMovieToken(token)
	if normalized == "" {
		return false
	}
	if _, ok := movieReleaseJunk[normalized]; ok {
		return true
	}
	return reMovieQuality.MatchString(normalized)
}
