package neyo

import (
	"testing"
)

func TestPayLoad(t *testing.T) {
	_, err := MakePayLoad("/tmp/test_blogs")
	if err != nil {
		t.Error(err)
	}
}
