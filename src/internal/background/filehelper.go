package background

import (
	"os"
	"path/filepath"
	"strings"
)

type fResult struct {
	Name string
	Path string
}

func getAllFilesWithPath(startPoint string, extensions []string) ([]fResult, error) {
	result := []fResult{}
	err := filepath.WalkDir(startPoint, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			// logger.Debug(nil, "DIR: {Dir}", path)
		} else {
			for _, extension := range extensions {
				if strings.HasSuffix(path, extension) {
					result = append(result, fResult{
						Path: path,
						Name: d.Name(),
					})
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
