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
	Duration      string
	Width         int
	Height        int
	Tags          map[string]string
}

type AudioStream struct {
	Index         int
	CodeName      string
	CodecLongName string
	Duration      string
	Tags          map[string]string
}

