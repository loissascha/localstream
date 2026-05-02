package provider

type ShowSearchResult struct {
	Score float64      `json:"score"`
	Show  ShowMetadata `json:"show"`
}

type ShowMetadata struct {
	ID        int        `json:"id"`
	URL       string     `json:"url"`
	Name      string     `json:"name"`
	Genres    []string   `json:"genres"`
	Premiered *string    `json:"premiered"`
	Image     *ShowImage `json:"image"`
	Summary   *string    `json:"summary"`
}

type SeasonMetadata struct {
	ID           int        `json:"id"`
	Url          string     `json:"url"`
	Number       int        `json:"number"`
	Summary      string     `json:"summary"`
	PremiereDate string     `json:"premiere_date"`
	Image        *ShowImage `json:"Image"`
}

type EpisodeMetadata struct {
	ID      int        `json:"id"`
	Url     string     `json:"url"`
	Name    string     `json:"name"`
	Number  int        `json:"number"`
	Summary string     `json:"summary"`
	Image   *ShowImage `json:"image"`
}

type ShowImage struct {
	Medium   string `json:"medium"`
	Original string `json:"original"`
}

type TVMetadataProvider interface {
	SearchShow(name string, year int) ([]ShowSearchResult, error)
	SearchSeasons(showId int) ([]SeasonMetadata, error)
	SearchEpisodes(seasonId int) ([]EpisodeMetadata, error)
}
