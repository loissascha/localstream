package parsers

import "testing"

func TestParseMovieFromFilename(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantOK    bool
		wantTitle string
		wantYear  *int
	}{
		{
			name:      "test deadpool wolverine",
			input:     "Deadpool.Wolverine.2024.1080p.BluRay.x264.AAC5.1-[YTS.MX].mp4",
			wantOK:    true,
			wantTitle: "Deadpool Wolverine",
			wantYear:  intPtr(2024),
		},
		{
			name:      "test blade",
			input:     "Blade - Trinity (2004) Unrated (1080p BluRay x265 10bit Tigole).mp4",
			wantOK:    true,
			wantTitle: "Blade Trinity",
			wantYear:  intPtr(2004),
		},
		{
			name:      "test atlantis",
			input:     "Atlantis The Lost Empire (2001) BDrip 1080p ENG-ITA x264 bluray -Shiv@.mp4",
			wantOK:    true,
			wantTitle: "Atlantis The Lost Empire",
			wantYear:  intPtr(2001),
		},
		{
			name:      "title with year and release tokens",
			input:     "The.Matrix.1999.1080p.BluRay.x264.mkv",
			wantOK:    true,
			wantTitle: "The Matrix",
			wantYear:  intPtr(1999),
		},
		{
			name:      "title with web dl tokens",
			input:     "Dune.Part.Two.2024.2160p.WEB-DL.HEVC.mp4",
			wantOK:    true,
			wantTitle: "Dune Part Two",
			wantYear:  intPtr(2024),
		},
		{
			name:      "path input uses last segment",
			input:     "/library/movies/Interstellar (2014).mkv",
			wantOK:    true,
			wantTitle: "Interstellar",
			wantYear:  intPtr(2014),
		},
		{
			name:      "keeps tv-like markers when present",
			input:     "Movie.Title.S01E02.2023.1080p.mkv",
			wantOK:    true,
			wantTitle: "Movie Title S01E02",
			wantYear:  intPtr(2023),
		},
		{
			name:      "without year",
			input:     "No.Country.for.Old.Men.1080p.BluRay.mkv",
			wantOK:    true,
			wantTitle: "No Country for Old Men",
			wantYear:  nil,
		},
		{
			name:      "numeric title without year stays valid",
			input:     "1917.mkv",
			wantOK:    true,
			wantTitle: "1917",
			wantYear:  nil,
		},
		{
			name:   "blank is invalid",
			input:  "   ",
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ParseMovieFromFilename(tt.input)
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if !tt.wantOK {
				return
			}

			if got == nil {
				t.Fatal("got nil info")
			}

			if got.Title != tt.wantTitle {
				t.Fatalf("title = %q, want %q", got.Title, tt.wantTitle)
			}

			switch {
			case tt.wantYear == nil && got.Year != nil:
				t.Fatalf("year = %v, want nil", *got.Year)
			case tt.wantYear != nil && got.Year == nil:
				t.Fatalf("year = nil, want %d", *tt.wantYear)
			case tt.wantYear != nil && got.Year != nil && *tt.wantYear != *got.Year:
				t.Fatalf("year = %d, want %d", *got.Year, *tt.wantYear)
			}
		})
	}
}
