package provider

type ShowSearchResult struct {
	Score float64
	Show  ShowMetadata
}

type ShowMetadata struct {
	ID        int
	URL       string
	Name      string
	Genres    []string
	Premiered *string
	Image     *ShowImage
	Summary   *string
}

type ShowImage struct {
	Medium   string
	Original string
}

type MovieSearchResult struct {
	Page    int           `json:"page"`
	Results []MovieResult `json:"results"`
}

type MovieResult struct {
	Adult         bool   `json:"adult"`
	OriginalTitle string `json:"original_title"`
	Overview      string `json:"overview"`
	ReleaseDate   string `json:"release_date"`
	BackdropPath  string `json:"backdrop_path"`
	PosterPath    string `json:"poster_path"`
}

type TVMetadataProvider interface {
	SearchShow(name string, year int) ([]ShowSearchResult, error)
}
