package media

import (
	"context"
	"fmt"
	"testing"
)

func TestIntegrityBrokenFile(t *testing.T) {
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

func TestIntegrityWorkingFile(t *testing.T) {
	s := NewIntegrityCheckService()

	ctx := context.Background()
	result, err := s.Scan(ctx, "Testfile.mp4")
	if err != nil {
		t.Fatal("error scanning Testfile.mp4", err)
	}

	if result.ExitCode == -1 {
		t.Fatal("exit code -1, meaning not started")
	}

	if result.Errors > 0 {
		t.Fatal("errors on Testfile.mp4")
	}

	fmt.Println("Exit Code:", result.ExitCode)

	fmt.Println("Scan result: ", result)
}
