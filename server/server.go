// Package server "企业通讯录数据同步接口v1"的参考实现
package server

import (
	"net/url"
	"strconv"

	spec "github.com/idaaser/syncspecv1"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// New 创建一个新的服务
func New(port int, opts ...Option) *Server {
	srv := &Server{
		port:     port,
		clients:  &allowAnyAs{},
		contacts: &nopcs{},
	}

	for _, opt := range opts {
		opt(srv)
	}

	return srv
}

type (
	// Server "企业通讯录数据同步接口v1"的参考实现
	Server struct {
		port int

		clients  AuthnStore
		contacts ContactStore
	}

	// Option Server可接受的配置选项
	Option func(srv *Server)
)

// Start 启动服务
func (s *Server) Start() {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.RequestIDWithConfig(
		middleware.RequestIDConfig{
			RequestIDHandler: func(c echo.Context, reqid string) {
				// 把X-Request-ID写到context中
				c.Set("reqid", reqid)
			},
		},
	))

	v1 := e.Group("/v1")
	v1.GET("/.well-known", s.wellknown)
	// 生成access_token
	v1.POST("/token", s.token)

	withAuth := v1.Group("", s.authn())
	// 分页获取部门详情
	withAuth.GET("/depts", s.listDepts)
	// 根据关键字, 搜索部门
	withAuth.GET("/depts/search", s.searchDept)
	// 分页获取指定部门下的用户详情
	withAuth.GET("/users", s.listUsersInDept)
	// 根据关键字, 搜索用户
	withAuth.GET("/users/search", s.serarchUser)

	// 分页获取group详情
	withAuth.GET("/groups", s.listGroups)
	// 根据关键字, 搜索group
	withAuth.GET("/groups/search", s.searchGroup)
	// 分页获取指定group下的用户id列表
	withAuth.GET("/groups/users", s.listUsersInGroup)

	// jit mock, for test only
	jit := v1.Group("/jit/:prefix/:count", s.jit())
	jit.GET("/.well-known", s.wellknown)

	// 生成access_token
	jit.POST("/token", s.token)
	jitAuth := jit.Group("", s.authn())
	// 分页获取部门详情
	jitAuth.GET("/depts", s.listDepts)
	// 分页获取指定部门下的用户详情
	jitAuth.GET("/users", s.listUsersInDept)

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(s.port)))
}

func (s *Server) absoluteURL(c echo.Context, paths ...string) string {
	u, _ := url.JoinPath(s.rootURL(c), paths...)
	return u
}

func (s *Server) rootURL(c echo.Context) string {
	return s.scheme(c) + "://" + s.host(c)
}

func (s *Server) scheme(c echo.Context) string {
	if xfp := c.Request().Header.Get("X-Forwarded-Proto"); xfp != "" {
		return xfp
	}

	return c.Scheme()
}

func (s *Server) host(c echo.Context) string {
	if xfh := c.Request().Header.Get("X-Forwarded-Host"); xfh != "" {
		return xfh
	}

	return c.Request().Host
}

func (s *Server) returnJSONError(c echo.Context, status int, code string, err error) error {
	resp := spec.ErrResponse{Code: code, Msg: err.Error()}
	if reqid, ok := c.Get("reqid").(string); ok {
		resp.RequestID = reqid
	}

	return c.JSON(status, resp)
}

func (s *Server) returnBadRequest(c echo.Context, err error) error {
	return s.returnJSONError(c, 400, spec.ErrInvalidRequest, err)
}
