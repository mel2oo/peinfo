package resource

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReader(t *testing.T) {
	r := NewReader()
	notepad := "../testdata/putty.exe"

	pe, err := r.Read(notepad)
	require.NoError(t, err)
	require.NotNil(t, pe)
}
