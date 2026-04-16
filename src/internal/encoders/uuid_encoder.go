package encoders

import (
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
)

func EncodeUUID(u uuid.UUID) string {
	return base64.RawURLEncoding.EncodeToString(u[:])
}

func DecodeUUID(s string) (uuid.UUID, error) {
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return uuid.Nil, err
	}
	if len(b) != 16 {
		return uuid.Nil, fmt.Errorf("invalid decoded length: got %d, want 16", len(b))
	}

	var u uuid.UUID
	copy(u[:], b)
	return u, nil
}
