package provider

type MovieResult struct {
	Adult        bool   `json:"adult"`
	Title        string `json:"original_title"`
	Description  string `json:"overview"`
	ReleaseDate  string `json:"release_date"`
	BackdropPath string `json:"backdrop_path"`
	PosterPath   string `json:"poster_path"`
}

type MovieMetadataProvider interface {
	SearchMovie(name string, year int) ([]MovieResult, error)
}
