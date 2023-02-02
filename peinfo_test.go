package peinfo

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestPeinfo(t *testing.T) {
	peinfo, err := New("testdata/putty.exe")
	if err != nil {
		t.Fail()
	}

	data, _ := json.Marshal(peinfo)
	fmt.Println(string(data))
}

func TestPeinfo2(t *testing.T) {
	peinfo, err := New("testdata/7587b1e66394558c8322c404941ca5a9.bin")
	if err != nil {
		t.Fail()
	}

	data, _ := json.Marshal(peinfo)
	fmt.Println(string(data))
}
