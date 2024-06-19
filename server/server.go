// Package server "企业通讯录数据同步接口v1"的参考实现
package server

import (
	"net/url"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// New 创建一个新的服务
func New(port int, opts ...Option) *Server {
	srv := &Server{
		port:    port,
		clients: &memoryClientStore{clients: map[string]string{}},
		tokens:  &anyTokenStore{},
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

		clients  ClientStore
		tokens   TokenStore
		contacts ContactStore
	}

	// Option Server可接受的参数
	Option func(srv *Server)
)

// Start 启动服务
func (s *Server) Start() {
	e := echo.New()
	e.Use(middleware.Recover())

	v1 := e.Group("/v1")
	v1.GET("/.well-known", s.wellknown)
	v1.POST("/token", s.token)

	withAuth := v1.Group("", s.authn())
	withAuth.GET("/users:search", s.serarchUser)
	withAuth.GET("/users", s.listUsersInDept)

	withAuth.GET("/depts:search", s.searchDept)
	withAuth.GET("/depts", s.listDepts)

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
