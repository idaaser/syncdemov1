package server

import (
	"strings"

	spec "github.com/idaaser/syncspecv1"
	"github.com/labstack/echo/v4"
)

func (s *Server) listDepts(c echo.Context) error {
	req := spec.ListDepatmentRequest{}
	if err := c.Bind(&req); err != nil {
		return s.returnBadRequest(c, err)
	}

	data, err := s.contacts.ListDepartments(c.Request().Context(), req)
	if err != nil {
		return s.returnBadRequest(c, err)
	}
	return c.JSON(200, spec.ListDepartmentResponse{PagingDepartments: *data})
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

	data, err := s.contacts.SearchDepartment(c.Request().Context(), keyword)
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

	data, err := s.contacts.ListUsersInDepartment(c.Request().Context(), req)
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

	data, err := s.contacts.SearchUser(c.Request().Context(), keyword)
	if err != nil {
		return s.returnBadRequest(c, err)
	}

	return c.JSON(200, &spec.SearchUserResponse{Data: data})
}
