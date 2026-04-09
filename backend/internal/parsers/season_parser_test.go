package parsers

import "testing"

func TestParseSeasonFromName(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantOK     bool
		wantSeason int
		wantType   SeasonType
	}{
		{
			name:       "season word",
			input:      "Season 1",
			wantOK:     true,
			wantSeason: 1,
			wantType:   SeasonTypeRegular,
		},
		{
			name:       "season word with leading zero",
			input:      "Season 01",
			wantOK:     true,
			wantSeason: 1,
			wantType:   SeasonTypeRegular,
		},
		{
			name:       "short s format",
			input:      "S2",
			wantOK:     true,
			wantSeason: 2,
			wantType:   SeasonTypeRegular,
		},
		{
			name:       "short s format with separator",
			input:      "S_03",
			wantOK:     true,
			wantSeason: 3,
			wantType:   SeasonTypeRegular,
		},
		{
			name:       "number only",
			input:      "04",
			wantOK:     true,
			wantSeason: 4,
			wantType:   SeasonTypeRegular,
		},
		{
			name:       "path input uses last segment",
			input:      "Show/Season 2",
			wantOK:     true,
			wantSeason: 2,
			wantType:   SeasonTypeRegular,
		},
		{
			name:       "extracts season word from noisy name",
			input:      "The.Show.Season.1.1080p.WEBRip.x264",
			wantOK:     true,
			wantSeason: 1,
			wantType:   SeasonTypeRegular,
		},
		{
			name:       "extracts short season marker from noisy name",
			input:      "The Show - S03 - 720p - BluRay",
			wantOK:     true,
			wantSeason: 3,
			wantType:   SeasonTypeRegular,
		},
		{
			name:       "specials",
			input:      "Specials",
			wantOK:     true,
			wantSeason: 0,
			wantType:   SeasonTypeSpecials,
		},
		{
			name:       "specials in noisy name",
			input:      "The Show Specials 1080p",
			wantOK:     true,
			wantSeason: 0,
			wantType:   SeasonTypeSpecials,
		},
		{
			name:       "sp alias",
			input:      "SP",
			wantOK:     true,
			wantSeason: 0,
			wantType:   SeasonTypeSpecials,
		},
		{
			name:       "s00 maps to specials",
			input:      "S00",
			wantOK:     true,
			wantSeason: 0,
			wantType:   SeasonTypeSpecials,
		},
		{
			name:       "season 0 maps to specials",
			input:      "Season 0",
			wantOK:     true,
			wantSeason: 0,
			wantType:   SeasonTypeSpecials,
		},
		{
			name:   "episode token is invalid",
			input:  "S01E02",
			wantOK: false,
		},
		{
			name:   "out of range is invalid",
			input:  "Season 100",
			wantOK: false,
		},
		{
			name:   "no season marker is invalid",
			input:  "Show Name 1080p x264",
			wantOK: false,
		},
		{
			name:       "number zero maps to specials",
			input:      "0",
			wantOK:     true,
			wantSeason: 0,
			wantType:   SeasonTypeSpecials,
		},
		{
			name:   "blank is invalid",
			input:  "   ",
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ParseSeasonFromName(tt.input)
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if !tt.wantOK {
				return
			}

			if got == nil {
				t.Fatal("got nil info")
			}

			if got.Season != tt.wantSeason {
				t.Fatalf("season = %d, want %d", got.Season, tt.wantSeason)
			}

			if got.Type != tt.wantType {
				t.Fatalf("type = %q, want %q", got.Type, tt.wantType)
			}
		})
	}
}
