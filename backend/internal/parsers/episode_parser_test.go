package parsers

import "testing"

func TestParseEpisodeFromFilename(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantOK      bool
		wantSeries  string
		wantSeason  int
		wantEpisode int
	}{
		{
			name:        "sxx epxx format",
			input:       "My.Show.S02.EP04.1080p.mkv",
			wantOK:      true,
			wantSeries:  "My Show",
			wantSeason:  2,
			wantEpisode: 4,
		},
		{
			name:        "compact sxxepxx format",
			input:       "My Show S02EP04",
			wantOK:      true,
			wantSeries:  "My Show",
			wantSeason:  2,
			wantEpisode: 4,
		},
		{
			name:        "lowercase with spaces",
			input:       "my show s02 ep04.mp4",
			wantOK:      true,
			wantSeries:  "my show",
			wantSeason:  2,
			wantEpisode: 4,
		},
		{
			name:        "existing sxxexx format still works",
			input:       "My.Show.S02E04.mkv",
			wantOK:      true,
			wantSeries:  "My Show",
			wantSeason:  2,
			wantEpisode: 4,
		},
		{
			name:        "existing x separator format still works",
			input:       "My Show 2x04",
			wantOK:      true,
			wantSeries:  "My Show",
			wantSeason:  2,
			wantEpisode: 4,
		},
		{
			name:   "missing episode marker is invalid",
			input:  "My Show Season 2 Episode 4",
			wantOK: false,
		},
		{
			name:   "zero values are invalid",
			input:  "My Show S00 EP00",
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ParseEpisodeFromFilename(tt.input)
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if !tt.wantOK {
				return
			}

			if got == nil {
				t.Fatal("got nil info")
			}

			if got.Series != tt.wantSeries {
				t.Fatalf("series = %q, want %q", got.Series, tt.wantSeries)
			}

			if got.Season != tt.wantSeason {
				t.Fatalf("season = %d, want %d", got.Season, tt.wantSeason)
			}

			if got.Episode != tt.wantEpisode {
				t.Fatalf("episode = %d, want %d", got.Episode, tt.wantEpisode)
			}
		})
	}
}
