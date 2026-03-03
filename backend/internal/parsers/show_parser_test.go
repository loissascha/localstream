package parsers

import "testing"

func intPtr(v int) *int {
	return &v
}

func TestParseShowFromName(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantOK     bool
		wantSeries string
		wantYear   *int
	}{
		{
			name:       "plain show name",
			input:      "Breaking Bad",
			wantOK:     true,
			wantSeries: "Breaking Bad",
			wantYear:   nil,
		},
		{
			name:       "show name with bare year",
			input:      "Dark 2017",
			wantOK:     true,
			wantSeries: "Dark",
			wantYear:   intPtr(2017),
		},
		{
			name:       "show name with parenthesized year",
			input:      "The Office (2005)",
			wantOK:     true,
			wantSeries: "The Office",
			wantYear:   intPtr(2005),
		},
		{
			name:       "show name with bracketed year",
			input:      "Shogun [2024]",
			wantOK:     true,
			wantSeries: "Shogun",
			wantYear:   intPtr(2024),
		},
		{
			name:       "path input only uses last segment",
			input:      "Series/The Office (US) (2005)",
			wantOK:     true,
			wantSeries: "The Office (US)",
			wantYear:   intPtr(2005),
		},
		{
			name:       "normalizes dots and underscores",
			input:      "The.Office_US (2005)",
			wantOK:     true,
			wantSeries: "The Office US",
			wantYear:   intPtr(2005),
		},
		{
			name:   "year only is invalid",
			input:  "(2020)",
			wantOK: false,
		},
		{
			name:   "blank input is invalid",
			input:  "   ",
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ParseShowFromName(tt.input)
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
