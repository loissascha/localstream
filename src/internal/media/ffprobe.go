package media

import (
	"encoding/json"
	"os/exec"
)

type FFProbeParse struct {
	Streams []MediaStream `json:"streams"`
	Format  FFProbeFormat `json:"format"`
}

type MediaStream struct {
	Index         int               `json:"index"`
	CodeName      string            `json:"codec_name"`
	CodecLongName string            `json:"codec_long_name"`
	CodecType     string            `json:"codec_type"`
	Duration      string            `json:"duration"`
	Width         int               `json:"width"`
	Height        int               `json:"height"`
	Tags          map[string]string `json:"tags"`
}

type FFProbeFormat struct {
	Filename       string            `json:"filename"`
	FormatName     string            `json:"format_name"`
	FormatLongName string            `json:"format_long_name"`
	StartTime      string            `json:"start_time"`
	Duration       string            `json:"duration"`
	Tags           map[string]string `json:"tags"`
}

func ProbeFile(path string) (*FFProbeParse, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_format", "-show_streams", path)
	raw, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var result FFProbeParse
	err = json.Unmarshal(raw, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
