package server

import (
	spec "github.com/idaaser/syncspecv1"
	"github.com/labstack/echo/v4"
)

func (s *Server) wellknown(c echo.Context) error {
	w := spec.Wellknown{
		TokenEndpoint:            s.absoluteURL(c, "v1", "token"),
		ListUsersInDeptEndpoint:  s.absoluteURL(c, "v1", "users"),
		SearchUserEndpoint:       s.absoluteURL(c, "v1", "users:search"),
		ListDepartmentsEndpoint:  s.absoluteURL(c, "v1", "depts"),
		SearchDepartmentEndpoint: s.absoluteURL(c, "v1", "depts:search"),
	}

	return c.JSON(200, w)
}
