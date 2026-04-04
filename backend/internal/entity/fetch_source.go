package entity

type FetchSource string

const (
	FetchSourceNone   FetchSource = "none"
	FetchSourceError  FetchSource = "error"
	FetchSourceTMDB   FetchSource = "tmdb"
	FetchSourceTVMaze FetchSource = "tvmaze"
)

func (t FetchSource) IsNone() bool {
	return t == FetchSourceNone
}

func (t FetchSource) IsError() bool {
	return t == FetchSourceError
}
