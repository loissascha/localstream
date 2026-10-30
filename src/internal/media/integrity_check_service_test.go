package media

import (
	"context"
	"fmt"
	"testing"
)

func TestIntegrityFile(t *testing.T) {
	s := NewIntegrityCheckService()

	ctx := context.Background()
	result, err := s.Scan(ctx, "BrokenFile.mp4")
	if err != nil {
		t.Fatal("error scanning BrokenFile.mp4", err)
	}

	if result.ExitCode == -1 {
		t.Fatal("exit code -1, meaning not started")
	}

	if result.Errors == 0 {
		t.Fatal("No errors on BrokenFile.mp4")
	}

	fmt.Println("Exit Code:", result.ExitCode)

	fmt.Println("Scan result: ", result)
}
