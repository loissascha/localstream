package media

import "testing"

func TestTranscodeService(t *testing.T) {
	ts := NewTranscodeService("./test")

	err := ts.StartTranscode("Movie1.mp4")
	if err != nil {
		t.Fatal(err)
	}
}
