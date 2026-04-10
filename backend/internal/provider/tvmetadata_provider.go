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

type TVMetadataProvider interface {
	SearchShow(name string, year int) ([]ShowSearchResult, error)
	SearchSeasons(showId int) error
}
