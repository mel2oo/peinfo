package resource

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReader(t *testing.T) {
	r := NewReader()
	notepad := "../testdata/ae70fc1fbbab5f1efb8a4e1fcafd6b3a856a400485f4092431cf1706e9962274"

	pe, err := r.Read(notepad)
	t.Log(pe)
	require.NoError(t, err)
	require.NotNil(t, pe)
}
