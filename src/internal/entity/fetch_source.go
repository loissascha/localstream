package entity

type FetchSource string

const (
	FetchSourceNone     FetchSource = "none"
	FetchSourceEmpty    FetchSource = "empty"
	FetchSourceMultiple FetchSource = "multiple"
	FetchSourceTMDB     FetchSource = "tmdb"
	FetchSourceTVMaze   FetchSource = "tvmaze"
)

func (t FetchSource) IsNone() bool {
	return t == FetchSourceNone
}

func (t FetchSource) IsEmpty() bool {
	return t == FetchSourceEmpty
}

func (t FetchSource) IsMultiple() bool {
	return t == FetchSourceMultiple
}
