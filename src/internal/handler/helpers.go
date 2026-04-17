package handler

import (
	"errors"
	"strconv"
	"strings"
)

func parseSingleRange(rangeHeader string, fileSize int64) (int64, int64, error) {
	if fileSize <= 0 {
		return 0, 0, errors.New("invalid file size")
	}

	if !strings.HasPrefix(rangeHeader, "bytes=") {
		return 0, 0, errors.New("invalid range unit")
	}

	rawRange := strings.TrimSpace(strings.TrimPrefix(rangeHeader, "bytes="))
	if rawRange == "" || strings.Contains(rawRange, ",") {
		return 0, 0, errors.New("multiple or empty ranges are not supported")
	}

	parts := strings.SplitN(rawRange, "-", 2)
	if len(parts) != 2 {
		return 0, 0, errors.New("malformed range")
	}

	startPart := strings.TrimSpace(parts[0])
	endPart := strings.TrimSpace(parts[1])

	if startPart == "" {
		suffixLength, err := strconv.ParseInt(endPart, 10, 64)
		if err != nil || suffixLength <= 0 {
			return 0, 0, errors.New("invalid suffix range")
		}
		if suffixLength > fileSize {
			suffixLength = fileSize
		}
		start := fileSize - suffixLength
		return start, fileSize - 1, nil
	}

	start, err := strconv.ParseInt(startPart, 10, 64)
	if err != nil || start < 0 || start >= fileSize {
		return 0, 0, errors.New("invalid start range")
	}

	if endPart == "" {
		return start, fileSize - 1, nil
	}

	end, err := strconv.ParseInt(endPart, 10, 64)
	if err != nil || end < start {
		return 0, 0, errors.New("invalid end range")
	}
	if end >= fileSize {
		end = fileSize - 1
	}

	return start, end, nil
}
