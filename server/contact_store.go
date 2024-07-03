package server

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"

	spec "github.com/idaaser/syncspecv1"
)

// WithContactStore 设置通讯录的存储
func WithContactStore(store ContactStore) Option {
	return func(srv *Server) {
		srv.contacts = store
	}
}

// WithContactFileStore 通讯录文件格式的存储
func WithContactFileStore(dept, user string) Option {
	return WithContactStore(&contactsFS{
		dept: newJSONFileStore[*spec.Department](dept),
		user: newJSONFileStore[*spec.User](user),
	})
}

// ContactStore 部门&用户存储
type ContactStore interface {
	// 分页返回部门列表
	ListDepartments(context.Context, spec.ListDepatmentRequest) (*spec.PagingDepartments, error)
	// 根据关键字模糊查询部门
	SearchDepartment(context.Context, string) ([]*spec.Department, error)

	// 分页返回指定部门下的直属用户列表, 不包括子孙部门下的用户
	ListUsersInDepartment(context.Context, spec.ListUsersInDepatmentRequest) (*spec.PagingUsers, error)
	// 根据关键字模糊查询用户
	SearchUser(context.Context, string) ([]*spec.User, error)
}

// nopcs 不返回任何数据的ContactStore. 注: 仅用于测试
type nopcs struct{}

// ListDepartments 实现ContactStore接口
func (c *nopcs) ListDepartments(context.Context, spec.ListDepatmentRequest) (
	*spec.PagingDepartments, error,
) {
	return &spec.PagingDepartments{}, nil
}

// ListUsersInDepartment implements ContactStore.
func (c *nopcs) ListUsersInDepartment(context.Context, spec.ListUsersInDepatmentRequest) (
	*spec.PagingUsers, error,
) {
	return &spec.PagingUsers{}, nil
}

func (c *nopcs) SearchDepartment(context.Context, string) ([]*spec.Department, error) {
	return []*spec.Department{}, nil
}

// SearchUser implements ContactStore.
func (c *nopcs) SearchUser(context.Context, string) ([]*spec.User, error) {
	return []*spec.User{}, nil
}

type contactsFS struct {
	dept *jsonFS[*spec.Department]
	user *jsonFS[*spec.User]
}

// ListDepartments implements ContactStore.
func (c *contactsFS) ListDepartments(ctx context.Context, req spec.ListDepatmentRequest) (
	*spec.PagingDepartments, error,
) {
	cursor, err := intCursor(req.Cursor).int()
	if err != nil {
		return nil, fmt.Errorf("invalid cursor %q", req.Cursor)
	}

	data, next := c.dept.sublist(cursor, req.GetSize())
	return &spec.PagingDepartments{
		HasNext: next != -1,
		Cursor: func() string {
			if next == -1 {
				return ""
			}
			return strconv.Itoa(next)
		}(),
		Data: data,
	}, nil
}

// ListUsersInDepartment implements ContactStore.
func (c *contactsFS) ListUsersInDepartment(ctx context.Context, req spec.ListUsersInDepatmentRequest) (
	*spec.PagingUsers, error,
) {
	deptid := req.DepartmentID
	cursor, err := intCursor(req.Cursor).int()
	if err != nil {
		return nil, fmt.Errorf("invalid cursor %q", req.Cursor)
	}

	all := []*spec.User{}
	for _, user := range c.user.load() {
		if deptid == user.MainDepartmentID || slices.Contains(user.OtherDepartmentsID, deptid) {
			all = append(all, user)
		}
	}
	data, next := sublist(all, cursor, req.GetSize())

	return &spec.PagingUsers{
		HasNext: next != -1,
		Cursor: func() string {
			if next == -1 {
				return ""
			}
			return strconv.Itoa(next)
		}(),
		Data: data,
	}, nil
}

// SearchDepartment implements ContactStore.
func (c *contactsFS) SearchDepartment(ctx context.Context, kw string) ([]*spec.Department, error) {
	if kw = strings.TrimSpace(kw); kw == "" {
		return []*spec.Department{}, nil
	}

	data := []*spec.Department{}
	for _, item := range c.dept.load() {
		if strings.EqualFold(item.Name, kw) || strings.EqualFold(item.ID, kw) {
			data = append(data, item)
		}
		// 返回前10个
		if len(data) >= 10 {
			break
		}
	}
	return data, nil
}

// SearchUser implements ContactStore.
func (c *contactsFS) SearchUser(_ context.Context, kw string) ([]*spec.User, error) {
	if kw = strings.TrimSpace(kw); kw == "" {
		return []*spec.User{}, nil
	}

	data := []*spec.User{}
	lower := strings.ToLower
	eq := func(s, chars string) bool {
		return strings.Contains(lower(s), lower(chars))
	}
	eqp := func(p *string, chars string) bool {
		return strings.Contains(lower(safes(p)), lower(chars))
	}
	for _, item := range c.user.load() {
		if eq(item.Name, kw) || eq(item.ID, kw) ||
			eqp(item.Username, kw) || eqp(item.Email, kw) ||
			eqp(item.Mobile, kw) || eqp(item.EmployeeNumber, kw) {
			data = append(data, item)
		}
		// 返回前10个
		if len(data) >= 10 {
			break
		}
	}
	return data, nil
}

func safes(sp *string) string {
	if sp == nil {
		return ""
	}
	return *sp
}

type intCursor string

func (c intCursor) int() (int, error) {
	if c == "" {
		return 0, nil
	}
	return strconv.Atoi(string(c))
}
