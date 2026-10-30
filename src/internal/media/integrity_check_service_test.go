package media

import (
	"context"
	"fmt"
	"testing"
)

func TestReencodeBrokenFile(t *testing.T) {
	s := NewIntegrityCheckService()
	ctx := context.Background()

	err := s.ReencodeFile(ctx, "BrokenFile.mp4", "NotBrokenFile.mp4")
	if err != nil {
		t.Fatal("Reencoding of broken file failed", err)
	}

	result, err := s.Scan(ctx, "NotBrokenFile.mp4")
	if err != nil {
		t.Fatal("Can't scan reencoded file")
	}

	fmt.Println("Scan result: ", result)

	if result.ExitCode == -1 {
		t.Fatal("exit code -1, scan not started")
	}
	if result.Errors != 0 {
		t.Fatal("reencoded file still has errors.")
	}
}

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
