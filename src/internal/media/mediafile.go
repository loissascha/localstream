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

func (f FFProbeParse) ToMediaFile() MediaFile {
	res := MediaFile{
		Filename:       f.Format.Filename,
		FormatName:     f.Format.FormatName,
		FormatLongName: f.Format.FormatLongName,
		StartTime:      f.Format.StartTime,
		Duration:       f.Format.Duration,
		VideoStreams:   []VideoStream{},
		AudioStreams:   []AudioStream{},
	}

	for _, s := range f.Streams {
		switch s.CodecType {
		case "video":
			res.VideoStreams = append(res.VideoStreams, VideoStream{
				Index:         s.Index,
				CodeName:      s.CodeName,
				CodecLongName: s.CodecLongName,
				Duration:      s.Duration,
				Width:         s.Width,
				Height:        s.Height,
				Tags:          s.Tags,
			})
		case "audio":
			res.AudioStreams = append(res.AudioStreams, AudioStream{
				Index:         s.Index,
				CodeName:      s.CodeName,
				CodecLongName: s.CodecLongName,
				Duration:      s.Duration,
				Tags:          s.Tags,
			})
		}
	}

	return res
}
