package media

import (
	"context"
	"testing"
)

func TestIntegrityFile(t *testing.T) {
	s := NewIntegrityCheckService()

	ctx := context.Background()
	s.Scan(ctx, "BrokenFile.mp4")
}
