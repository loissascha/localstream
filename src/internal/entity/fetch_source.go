package entity

type FetchSource string

// None = Not tried to fetch yet -> try it next time the matcher runs
// Empty = Tried to fetch data but couldn't find anything -> will not try to fetch when matcher runs
// Multiple = Has multiple results, admin needs to select the correct one!

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
