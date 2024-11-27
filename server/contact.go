package server

import (
	"strconv"
	"strings"

	spec "github.com/idaaser/syncspecv1"
	"github.com/labstack/echo/v4"
)

func (s *Server) listDepts(c echo.Context) error {
	req := spec.ListDepatmentRequest{}
	if err := c.Bind(&req); err != nil {
		return s.returnBadRequest(c, err)
	}

	data, err := s.getContactStore(c).
		ListDepartments(c.Request().Context(), req)
	if err != nil {
		return s.returnBadRequest(c, err)
	}
	return c.JSON(200, spec.ListDepartmentResponse{PagingDepartments: *data})
}

func (s *Server) jit() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			store := s.loadJitStore(c)
			c.Set("_store_", store)
			return next(c)
		}
	}
}

func (s *Server) loadJitStore(c echo.Context) *jitStore {
	prefix := c.Param("prefix")
	count := strings.Split(c.Param("count"), ",")
	dept, user := 10, 10
	if len(count) == 2 {
		dept, _ = strconv.Atoi(count[0])
		user, _ = strconv.Atoi(count[1])
	}

	return newJITContactStore(prefix, dept, user)
}

func (s *Server) searchDept(c echo.Context) error {
	req := spec.SearchDepartmentRequest{}
	if err := c.Bind(&req); err != nil {
		return s.returnBadRequest(c, err)
	}
	keyword := strings.TrimSpace(req.Keyword)
	if keyword == "" {
		return c.JSON(200, spec.SearchDepartmentResponse{})
	}

	data, err := s.getContactStore(c).
		SearchDepartment(c.Request().Context(), keyword)
	if err != nil {
		return s.returnBadRequest(c, err)
	}

	return c.JSON(200, &spec.SearchDepartmentResponse{Data: data})
}

func (s *Server) listUsersInDept(c echo.Context) error {
	req := spec.ListUsersInDepatmentRequest{}
	if err := c.Bind(&req); err != nil {
		return s.returnBadRequest(c, err)
	}
	if err := req.Validate(); err != nil {
		return s.returnBadRequest(c, err)
	}

	data, err := s.getContactStore(c).ListUsersInDepartment(c.Request().Context(), req)
	if err != nil {
		return s.returnBadRequest(c, err)
	}
	return c.JSON(200, spec.ListUsersInDepartmentResponse{PagingUsers: *data})
}

func (s *Server) serarchUser(c echo.Context) error {
	req := spec.SearchUserRequest{}
	if err := c.Bind(&req); err != nil {
		return s.returnBadRequest(c, err)
	}
	keyword := strings.TrimSpace(req.Keyword)
	if keyword == "" {
		return c.JSON(200, spec.SearchUserResponse{})
	}

	data, err := s.getContactStore(c).
		SearchUser(c.Request().Context(), keyword)
	if err != nil {
		return s.returnBadRequest(c, err)
	}

	return c.JSON(200, &spec.SearchUserResponse{Data: data})
}

func (s *Server) getContactStore(c echo.Context) ContactStore {
	store := c.Get("_store_")
	if store == nil {
		return s.contacts
	}

	return store.(ContactStore)
}
