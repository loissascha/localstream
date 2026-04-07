package provider

type MovieMetadataProvider interface {
	SearchMovie(name string, year int)
}
