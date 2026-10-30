package media

import "testing"

func TestIntegrityFile(t *testing.T) {
	s := NewIntegrityCheckService()
	s.CheckFFMpegErrors("BrokenFile.mp4")
}
