package peinfo

import (
	"fmt"
	"testing"
)

func TestPeinfo(t *testing.T) {
	peinfo, err := New("testdata/putty.exe")
	if err != nil {
		t.Fail()
	}

	fmt.Println(peinfo)
}

func TestPeinfo2(t *testing.T) {
	peinfo, err := New("testdata/9edc89c4489fb29a61a1e1ff66d9d49b4fb2c230d0efb323dbf1829895497ba8")
	if err != nil {
		t.Fail()
	}

	fmt.Println(peinfo)
}
