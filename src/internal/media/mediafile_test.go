package media

import (
	"fmt"
	"testing"
)

func TestParse2Fast(t *testing.T) {

	fmt.Println("Test first file...")
	parse, err := ProbeFile("./2.Fast.2.Furious.mp4")
	if err != nil {
		fmt.Println("error", err)
		t.Fatal(err)
	}
	fmt.Println(parse)

	fmt.Println("Test second file...")
}
