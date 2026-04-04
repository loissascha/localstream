package entity

type FetchSource string

const (
	FetchSourceNone   FetchSource = "none"
	FetchSourceTMDB   FetchSource = "tmdb"
	FetchSourceTVMaze FetchSource = "tvmaze"
)

func (t FetchSource) IsValid() bool {
	switch t {
	case FetchSourceTMDB, FetchSourceTVMaze:
		return true
	default:
		return false
	}
}
