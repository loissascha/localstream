package media

import (
	"fmt"
	"testing"
	"time"
)

func TestTranscodeService(t *testing.T) {
	ts := NewTranscodeService("./test")

	_, err := ts.StartTranscode("NotExistingFile.mp4")
	if err == nil {
		t.Error("No error on not existing file!")
	}

	sess, err := ts.StartTranscode("Testfile.mp4")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Transcoding running for id:", sess.ID, "output:", sess.OutputDir)
	fmt.Println("Transcoding running for id:", sess.ID, "output:", sess.OutputDir)
	fmt.Println("Transcoding running for id:", sess.ID, "output:", sess.OutputDir)
	fmt.Println("Transcoding running for id:", sess.ID, "output:", sess.OutputDir)
	fmt.Println("Transcoding running for id:", sess.ID, "output:", sess.OutputDir)
	fmt.Println("Transcoding running for id:", sess.ID, "output:", sess.OutputDir)
	fmt.Println("Transcoding running for id:", sess.ID, "output:", sess.OutputDir)
	fmt.Println("Transcoding running for id:", sess.ID, "output:", sess.OutputDir)
	fmt.Println("Transcoding running for id:", sess.ID, "output:", sess.OutputDir)
	fmt.Println("Transcoding running for id:", sess.ID, "output:", sess.OutputDir)
	fmt.Println("Transcoding running for id:", sess.ID, "output:", sess.OutputDir)
	fmt.Println("Transcoding running for id:", sess.ID, "output:", sess.OutputDir)
	fmt.Println("Transcoding running for id:", sess.ID, "output:", sess.OutputDir)

	time.Sleep(30 * time.Second)
	ts.StopTranscode(sess.ID)

	fmt.Println("Transcoding stopped!!!!!")
	fmt.Println("Transcoding stopped!!!!!")
	fmt.Println("Transcoding stopped!!!!!")
	fmt.Println("Transcoding stopped!!!!!")
	fmt.Println("Transcoding stopped!!!!!")
	fmt.Println("Transcoding stopped!!!!!")
	fmt.Println("Transcoding stopped!!!!!")

}
