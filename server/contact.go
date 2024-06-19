package server

import (
	"context"
	"net/http"

	spec "github.com/idaaser/syncspecv1"
	"github.com/labstack/echo/v4"
)

// WithContactStore 设置通讯录的存储
func WithContactStore(store ContactStore) Option {
	return func(srv *Server) {
		srv.contacts = store
	}
}

// WithContactFileStore 通讯录文件格式的存储
func WithContactFileStore(dept, user string) Option {
	return WithContactStore(&contactsFileStore{
		dept: newJsonFileStore[*spec.Department](dept),
		user: newJsonFileStore[*spec.User](user),
	})
}

// ContactStore 部门&用户存储
type ContactStore interface {
	// 分页返回部门列表
	ListDepartments(context.Context, spec.PagingRequest) (*spec.PagingDepartments, error)
	// 根据关键字模糊查询部门
	SearchDepartment(context.Context, string) (*spec.PagingDepartments, error)

	// 分页返回指定部门下的直属用户列表, 不包括子孙部门下的用户
	ListUsersInDepartment(context.Context, string, spec.PagingRequest) (*spec.PagingUsers, error)
	// 根据关键字模糊查询用户
	SearchUser(context.Context, string) (*spec.PagingUsers, error)
}

func (s *Server) listDepts(c echo.Context) error {
	req := spec.PagingRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			spec.ErrResponse{Error: "invalid_request", ErrorMessage: err.Error()},
		)
	}

	return c.JSON(http.StatusOK,
		spec.PagingResult[spec.Department]{
			HasNext: false, Cursor: "",
			Data: []spec.Department{
				{ID: "1", Name: "dept 1", Parent: ""},
				{ID: "1-1", Name: "dept 1-1", Parent: "1"},
			},
		},
	)
}

func (s *Server) searchDept(c echo.Context) error {
	return nil
}

func (s *Server) listUsersInDept(c echo.Context) error {
	return nil
}

func (s *Server) serarchUser(c echo.Context) error {
	return nil
}

type contactsFileStore struct {
	dept *jsonFileStore[*spec.Department]
	user *jsonFileStore[*spec.User]
}

// ListDepartments implements ContactStore.
func (c *contactsFileStore) ListDepartments(context.Context, spec.PagingRequest) (*spec.PagingResult[*spec.Department], error) {
	panic("unimplemented")
}

// ListUsersInDepartment implements ContactStore.
func (c *contactsFileStore) ListUsersInDepartment(context.Context, string, spec.PagingRequest) (*spec.PagingResult[*spec.User], error) {
	panic("unimplemented")
}

// SearchDepartment implements ContactStore.
func (c *contactsFileStore) SearchDepartment(context.Context, string) (*spec.PagingResult[*spec.Department], error) {
	panic("unimplemented")
}

// SearchUser implements ContactStore.
func (c *contactsFileStore) SearchUser(context.Context, string) (*spec.PagingResult[*spec.User], error) {
	panic("unimplemented")
}
