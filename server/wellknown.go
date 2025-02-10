package server

import (
	"strings"

	spec "github.com/idaaser/syncspecv1"
	"github.com/labstack/echo/v4"
)

func (s *Server) wellknown(c echo.Context) error {
	u := strings.TrimSuffix(c.Request().URL.String(), ".well-known")
	w := spec.Wellknown{
		TokenEndpoint:            s.absoluteURL(c, u, "token"),
		ListUsersInDeptEndpoint:  s.absoluteURL(c, u, "users"),
		SearchUserEndpoint:       s.absoluteURL(c, u, "users/search"),
		ListDepartmentsEndpoint:  s.absoluteURL(c, u, "depts"),
		SearchDepartmentEndpoint: s.absoluteURL(c, u, "depts/search"),
		ListGroupsEndpoint:       s.absoluteURL(c, u, "groups"),
		SearchGroupEndpoint:      s.absoluteURL(c, u, "groups/search"),
		ListUsersInGroupEndpoint: s.absoluteURL(c, u, "groups/users"),
	}

	return c.JSON(200, w)
}
