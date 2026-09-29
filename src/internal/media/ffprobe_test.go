package media

import (
	"fmt"
	"strings"
	"testing"
)

func TestParse2Fast(t *testing.T) {

	fmt.Println("Test first file...")
	parse, err := ProbeFile("Testfile.mp4")
	if err != nil {
		fmt.Println("error", err)
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", parse)

	mediaF, err := parse.ToMediaFile()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("\nas mediaFile: %+v\n", mediaF)

	splits := strings.Split(parse.Format.Filename, ".")
	container := splits[len(splits)-1]
	fmt.Println("container:", container)
}
