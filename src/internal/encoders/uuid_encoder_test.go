package encoders

import (
	"testing"

	"github.com/google/uuid"
)

func TestEncoder(t *testing.T) {
	uid, err := uuid.NewV7()
	if err != nil {
		t.Fatal(err)
	}

	encoded := EncodeUUID(uid)
	decoded, err := DecodeUUID(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if decoded != uid {
		t.Error("the decoded uid is not equal to the original uid. Test failed.")
	}
}
