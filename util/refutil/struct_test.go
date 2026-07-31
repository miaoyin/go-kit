package refutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type record struct {
	ID int64
}

func TestStructField(t *testing.T) {
	rd := &record{}
	sr := NewStructRef(rd)
	assert.Equal(t, sr.ExistField("ID"), true)
	assert.Equal(t, sr.ExistField("ID1"), false)
	assert.Equal(t, sr.GetFieldValue("ID"), int64(0))
	assert.Equal(t, sr.GetFieldValue("ID1"), int64(0))
}
