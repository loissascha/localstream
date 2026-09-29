package media

import "testing"

func TestTranscodeService(t *testing.T) {
	ts := NewTranscodeService("./test")

	err := ts.StartTranscode("NotExistingFile.mp4")
	if err == nil {
		t.Error("No error on not existing file!")
	}

	err = ts.StartTranscode("Movie1.mp4")
	if err != nil {
		t.Fatal(err)
	}

}
