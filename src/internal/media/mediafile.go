package media

type MediaFile struct {
	Filename       string
	FormatName     string
	FormatLongName string
	StartTime      string
	Duration       float64
	BitRate        int64
	Size           int64
	VideoStreams   []VideoStream
	AudioStreams   []AudioStream
}

type VideoStream struct {
	Index         int
	Profile       string
	Level         int
	CodecName      string
	CodecLongName string
	Duration      float64 // in sekunden
	DurationTS    int64
	Width         int
	Height        int
	Tags          map[string]string
	BitRate       int64
}

type AudioStream struct {
	Index         int
	Profile       string
	CodecName      string
	CodecLongName string
	Duration      float64 // in sekunden
	DurationTS    int64
	Tags          map[string]string
	BitRate       int64
}
