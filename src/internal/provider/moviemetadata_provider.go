package provider

type MovieResult struct {
	FetchID      int    `json:"id"`
	Adult        bool   `json:"adult"`
	Title        string `json:"original_title"`
	Description  string `json:"overview"`
	ReleaseYear  int    `json:"release_year"`
	BackdropPath string `json:"backdrop_path"`
	PosterPath   string `json:"poster_path"`
}

type MovieMetadataProvider interface {
	GetMovieByID(id int) (*MovieResult, error)
	SearchMovie(name string, year int) ([]MovieResult, error)
}
