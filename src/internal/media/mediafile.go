package media

type MediaFile struct {
	Filename       string
	FormatName     string
	FormatLongName string
	StartTime      string
	Duration       string
	VideoStreams   []VideoStream
	AudioStreams   []AudioStream
}

type VideoStream struct {
	Index         int
	CodeName      string
	CodecLongName string
	Duration      float64 // in sekunden
	DurationTS    int64
	Width         int
	Height        int
	Tags          map[string]string
}

type AudioStream struct {
	Index         int
	CodeName      string
	CodecLongName string
	Duration      float64 // in sekunden
	DurationTS    int64
	Tags          map[string]string
}
