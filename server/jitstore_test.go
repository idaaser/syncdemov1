package server

import (
	"context"
	"testing"

	spec "github.com/idaaser/syncspecv1"
	"github.com/stretchr/testify/assert"
)

func Test_jitStore_ListDepartments(t *testing.T) {
	store := newJITContactStore("beijing", 10, 1)

	size, cursor := 3, ""

	// first page, [0:3]
	first, err := store.ListDepartments(context.TODO(),
		spec.PagingParam{Size: size, Cursor: cursor})
	assert.NoError(t, err)
	assert.True(t, first.HasNext)
	assert.Equal(t, "3", first.Cursor)
	assert.Len(t, first.Data, size)

	// second page, [3:6]
	cursor = first.Cursor
	second, err := store.ListDepartments(context.TODO(),
		spec.PagingParam{Size: size, Cursor: cursor})
	assert.NoError(t, err)
	assert.True(t, second.HasNext)
	assert.Equal(t, "6", second.Cursor)
	assert.Len(t, second.Data, size)

	// third page, [6:9]
	cursor = second.Cursor
	third, err := store.ListDepartments(context.TODO(),
		spec.PagingParam{Size: size, Cursor: cursor})
	assert.NoError(t, err)
	assert.True(t, third.HasNext)
	assert.Equal(t, "9", third.Cursor)
	assert.Len(t, third.Data, size)

	// last page, [9:]
	cursor = third.Cursor
	fourth, err := store.ListDepartments(context.TODO(),
		spec.PagingParam{Size: size, Cursor: cursor})
	assert.NoError(t, err)
	assert.False(t, fourth.HasNext)
	assert.Equal(t, "", fourth.Cursor)
	assert.Len(t, fourth.Data, 1)

	// 请求不存在的index(大于总数)
	notexists, err := store.ListDepartments(context.TODO(),
		spec.PagingParam{Size: size, Cursor: "10"})
	assert.NoError(t, err)
	assert.False(t, notexists.HasNext)
	assert.Equal(t, "", notexists.Cursor)
	assert.Len(t, notexists.Data, 0)
}

func Test_jitStore_ListUsersInDepartment(t *testing.T) {
	store := newJITContactStore("beijing", 1, 10)

	size, cursor := 3, ""

	// first page, [0:3]
	first, err := store.ListUsersInDepartment(context.TODO(),
		spec.ListUsersInDepatmentRequest{
			DepartmentID: "beijing-1",
			PagingParam:  spec.PagingParam{Size: size, Cursor: cursor},
		})
	assert.NoError(t, err)
	assert.True(t, first.HasNext)
	assert.Equal(t, "3", first.Cursor)
	assert.Len(t, first.Data, size)

	// second page, [3:6]
	cursor = first.Cursor
	second, err := store.ListUsersInDepartment(context.TODO(),
		spec.ListUsersInDepatmentRequest{
			DepartmentID: "beijing-1",
			PagingParam:  spec.PagingParam{Size: size, Cursor: cursor},
		})
	assert.NoError(t, err)
	assert.True(t, second.HasNext)
	assert.Equal(t, "6", second.Cursor)
	assert.Len(t, second.Data, size)

	// third page, [6:9]
	cursor = second.Cursor
	third, err := store.ListUsersInDepartment(context.TODO(),
		spec.ListUsersInDepatmentRequest{
			DepartmentID: "beijing-1",
			PagingParam:  spec.PagingParam{Size: size, Cursor: cursor},
		})
	assert.NoError(t, err)
	assert.True(t, third.HasNext)
	assert.Equal(t, "9", third.Cursor)
	assert.Len(t, third.Data, size)

	// last page, [9:]
	cursor = third.Cursor
	fourth, err := store.ListUsersInDepartment(context.TODO(),
		spec.ListUsersInDepatmentRequest{
			DepartmentID: "beijing-1",
			PagingParam:  spec.PagingParam{Size: size, Cursor: cursor},
		})
	assert.NoError(t, err)
	assert.False(t, fourth.HasNext)
	assert.Equal(t, "", fourth.Cursor)
	assert.Len(t, fourth.Data, 1)

	// 请求不存在的index(大于总数)
	notexists, err := store.ListUsersInDepartment(context.TODO(),
		spec.ListUsersInDepatmentRequest{
			DepartmentID: "beijing-1",
			PagingParam:  spec.PagingParam{Size: size, Cursor: "10"},
		})
	assert.NoError(t, err)
	assert.False(t, notexists.HasNext)
	assert.Equal(t, "", notexists.Cursor)
	assert.Len(t, notexists.Data, 0)
}
