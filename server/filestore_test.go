package server

import (
	"testing"

	spec "github.com/idaaser/syncspecv1"
	"github.com/stretchr/testify/assert"
)

func Test_jsonfile_store(t *testing.T) {
	{
		store := newJSONFileStore[*spec.Department]("./testdata/departments.json")
		data := store.load()
		assert.NotEmpty(t, data)
		assert.Equal(t, "1", data[0].ID)
	}

	{
		store := newJSONFileStore[*spec.User]("./testdata/users.json")
		data := store.load()
		assert.NotEmpty(t, data)
		assert.Equal(t, "uid-1", data[0].ID)
	}
}
