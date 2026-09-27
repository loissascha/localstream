package media

import (
	"fmt"
	"testing"
)

func TestParse2Fast(t *testing.T) {

	fmt.Println("Test first file...")
	parse, err := ProbeFile("Movie1.mp4")
	if err != nil {
		fmt.Println("error", err)
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", parse)

	mediaF := parse.ToMediaFile()
	fmt.Printf("\nas mediaFile: %+v\n", mediaF)

	fmt.Println("")
	fmt.Println("Test second file...")
	parse, err = ProbeFile("Movie2.mp4")
	if err != nil {
		fmt.Println("error", err)
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", parse)

	mediaF = parse.ToMediaFile()
	fmt.Printf("\nas mediaFile: %+v\n", mediaF)

	fmt.Println("")
	fmt.Println("Test third file...")
	parse, err = ProbeFile("Show1.mp4")
	if err != nil {
		fmt.Println("error", err)
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", parse)

	mediaF = parse.ToMediaFile()
	fmt.Printf("\nas mediaFile: %+v\n", mediaF)
}
