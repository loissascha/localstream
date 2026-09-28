package media

import (
	"encoding/json"
	"os/exec"
	"strconv"
)

type FFProbeParse struct {
	Streams []MediaStream `json:"streams"`
	Format  FFProbeFormat `json:"format"`
}

type MediaStream struct {
	Index         int               `json:"index"`
	Profile       string            `json:"profile"`
	Level         int               `json:"level"`
	CodecName      string            `json:"codec_name"`
	CodecLongName string            `json:"codec_long_name"`
	CodecType     string            `json:"codec_type"`
	Duration      string            `json:"duration"`
	DurationTS    int64             `json:"duration_ts"`
	Width         int               `json:"width"`
	Height        int               `json:"height"`
	Tags          map[string]string `json:"tags"`
	BitRate       string            `json:"bit_rate"`
}

type FFProbeFormat struct {
	Filename       string            `json:"filename"`
	FormatName     string            `json:"format_name"`
	FormatLongName string            `json:"format_long_name"`
	StartTime      string            `json:"start_time"`
	Duration       string            `json:"duration"`
	Tags           map[string]string `json:"tags"`
	BitRate        string            `json:"bit_rate"`
	Size           string            `json:"size"`
}

func ProbeFile(path string) (FFProbeParse, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_format", "-show_streams", path)
	raw, err := cmd.CombinedOutput()
	if err != nil {
		return FFProbeParse{}, err
	}

	var result FFProbeParse
	err = json.Unmarshal(raw, &result)
	if err != nil {
		return FFProbeParse{}, err
	}

	return result, nil
}

func (f *FFProbeParse) ToMediaFile() (MediaFile, error) {
	dur, err := strconv.ParseFloat(f.Format.Duration, 64)
	if err != nil {
		return MediaFile{}, err
	}
	bitRate, err := strconv.ParseInt(f.Format.BitRate, 10, 64)
	if err != nil {
		return MediaFile{}, err
	}
	size, err := strconv.ParseInt(f.Format.Size, 10, 64)
	if err != nil {
		return MediaFile{}, err
	}
	res := MediaFile{
		Filename:       f.Format.Filename,
		FormatName:     f.Format.FormatName,
		FormatLongName: f.Format.FormatLongName,
		StartTime:      f.Format.StartTime,
		Duration:       dur,
		BitRate:        bitRate,
		Size:           size,
		VideoStreams:   []VideoStream{},
		AudioStreams:   []AudioStream{},
	}

	for _, s := range f.Streams {
		dur, err := strconv.ParseFloat(s.Duration, 64)
		if err != nil {
			return MediaFile{}, err
		}
		bitRate, err := strconv.ParseInt(s.BitRate, 10, 64)
		if err != nil {
			return MediaFile{}, err
		}
		switch s.CodecType {
		case "video":
			res.VideoStreams = append(res.VideoStreams, VideoStream{
				Index:         s.Index,
				Profile:       s.Profile,
				Level:         s.Level,
				CodecName:      s.CodecName,
				CodecLongName: s.CodecLongName,
				Duration:      dur,
				DurationTS:    s.DurationTS,
				Width:         s.Width,
				Height:        s.Height,
				Tags:          s.Tags,
				BitRate:       bitRate,
			})
		case "audio":
			res.AudioStreams = append(res.AudioStreams, AudioStream{
				Index:         s.Index,
				Profile:       s.Profile,
				CodecName:      s.CodecName,
				CodecLongName: s.CodecLongName,
				Duration:      dur,
				DurationTS:    s.DurationTS,
				Tags:          s.Tags,
				BitRate:       bitRate,
			})
		}
	}

	return res, nil
}
