package server

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	spec "github.com/idaaser/syncspecv1"
)

// 请求时自动生成通讯录数据的store, 用于mock测试

func newJITContactStore(prefix string, dept, user int) *jitStore {
	return &jitStore{
		dept:   dept,
		user:   user,
		prefix: prefix,
	}
}

type (
	jitStore struct {
		// 部门总数, 深度为1
		dept int

		// 每个部门下的用户数量
		user int

		// 用户部门的前缀
		prefix string
	}
)

func (s *jitStore) ListDepartments(_ context.Context, req spec.ListDepatmentRequest) (*spec.PagingDepartments, error) {
	paging := asindexBasedPaging(req)
	start, end := paging.start(), paging.end()
	if start >= s.dept {
		return &spec.PagingDepartments{HasNext: false}, nil
	}
	end = min(end, s.dept-1)

	data := []*spec.Department{}
	for i := start; i <= end; i++ {
		data = append(data, s.newDepartment(i))
	}

	hasMore := end < (s.dept - 1)
	next := ""
	if hasMore {
		next = strconv.Itoa(end + 1)
	}

	return &spec.PagingDepartments{
		Data:    data,
		HasNext: hasMore,
		Cursor:  next,
	}, nil
}

func (s *jitStore) SearchDepartment(context.Context, string) ([]*spec.Department, error) {
	return nil, errors.New("unsupported")
}

// 根据关键字模糊查询用户
func (s *jitStore) SearchUser(context.Context, string) ([]*spec.User, error) {
	return nil, errors.New("unsupported")
}

// 分页返回指定部门下的直属用户列表, 不包括子孙部门下的用户
func (s *jitStore) ListUsersInDepartment(_ context.Context, req spec.ListUsersInDepatmentRequest) (*spec.PagingUsers, error) {
	paging := asindexBasedPaging(req.PagingParam)
	start, end := paging.start(), paging.end()
	if start >= s.user {
		return &spec.PagingUsers{HasNext: false}, nil
	}
	end = min(end, s.user-1)

	data := []*spec.User{}
	for i := start; i <= end; i++ {
		data = append(data, s.newUser(req.DepartmentID, i))
	}

	hasMore := end < (s.user - 1)
	next := ""
	if hasMore {
		next = strconv.Itoa(end + 1)
	}

	return &spec.PagingUsers{
		Data:    data,
		HasNext: hasMore,
		Cursor:  next,
	}, nil
}

func (s *jitStore) newDepartment(index int) *spec.Department {
	width := fmt.Sprintf("%d", len(strconv.Itoa(s.dept)))
	format := "%s-%0" + width + "d"

	name := fmt.Sprintf(format, s.prefix, index+1)
	return &spec.Department{
		Name: name, ID: name,
		Parent: "",
		Order:  index + 1,
	}
}

func (s *jitStore) newUser(deptid string, index int) *spec.User {
	u := &spec.User{}
	id := fmt.Sprintf("%s-u-%d", deptid, index)
	u.ID = id
	u.Username = spec.Pointer(id)
	u.Email = spec.Pointer(id + "@mailinator.com")

	u.Name = id
	u.Position = spec.Pointer("mock")
	u.EmployeeNumber = spec.Pointer(id)
	u.Status = spec.UserStatusInitialized
	u.Order = index
	u.MainDepartmentID = deptid

	return u
}

// 用户总数
func (s *jitStore) totalUsers() int {
	return s.user * s.dept
}

type indexBasedPaging struct {
	// 数组下标, 0开始
	idx int

	// 单页数量
	size int
}

func (p indexBasedPaging) start() int {
	return p.idx
}

func (p indexBasedPaging) end() int {
	return p.start() + p.size - 1
}

func asindexBasedPaging(p spec.PagingParam) indexBasedPaging {
	idx := 0
	if p.Cursor != "" {
		idx, _ = strconv.Atoi(p.Cursor)
	}
	if idx < 0 {
		idx = 0
	}

	return indexBasedPaging{
		idx:  idx,
		size: p.GetSize(),
	}
}
