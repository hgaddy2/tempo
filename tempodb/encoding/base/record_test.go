package base

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/grafana/tempo/tempodb/encoding/common"
	"github.com/stretchr/testify/assert"
)

func TestEncodeDecodeRecord(t *testing.T) {
	expected, err := makeRecord(t)
	assert.NoError(t, err, "unexpected error making trace record")

	buff := make([]byte, RecordLength)

	marshalRecord(expected, buff)
	actual := UnmarshalRecord(buff)

	assert.Equal(t, expected, actual)
}

func TestMarshalUnmarshalRecords(t *testing.T) {
	numRecords := 10
	expected := make([]*common.Record, 0, numRecords)

	for i := 0; i < numRecords; i++ {
		r, err := makeRecord(t)
		if err != nil {
			assert.NoError(t, err, "unexpected error making trace record")
		}
		expected = append(expected, r)
	}

	recordBytes, err := MarshalRecords(expected)
	assert.NoError(t, err, "unexpected error encoding records")
	assert.Equal(t, len(expected)*28, len(recordBytes))

	actual, err := unmarshalRecords(recordBytes)
	assert.NoError(t, err, "unexpected error decoding records")

	assert.Equal(t, expected, actual)
}

func TestSortRecord(t *testing.T) {
	numRecords := 10
	expected := make([]*common.Record, 0, numRecords)

	for i := 0; i < numRecords; i++ {
		r, err := makeRecord(t)
		if err != nil {
			assert.NoError(t, err, "unexpected error making trace record")
		}
		expected = append(expected, r)
	}

	SortRecords(expected)

	for i := range expected {
		if i == 0 {
			continue
		}

		idSmaller := expected[i-1].ID
		idLarger := expected[i].ID

		assert.NotEqual(t, 1, bytes.Compare(idSmaller, idLarger))
	}
}

// todo: belongs in util/test?
func makeRecord(t *testing.T) (*common.Record, error) {
	t.Helper()

	r := newRecord()
	r.Start = rand.Uint64()
	r.Length = rand.Uint32()

	_, err := rand.Read(r.ID)
	if err != nil {
		return nil, err
	}

	return r, nil
}
