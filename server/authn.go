package server

import (
	"fmt"
	"strings"

	spec "github.com/idaaser/syncspecv1"

	"github.com/labstack/echo/v4"
)

func (s *Server) token(c echo.Context) error {
	req := spec.GetTokenRequest{}
	if err := c.Bind(&req); err != nil {
		return s.returnBadRequest(c, err)
	}

	if clientid, clientsecret, ok := c.Request().BasicAuth(); ok {
		req.ClientID = clientid
		req.ClientSecret = clientsecret
	}

	if err := req.Validate(); err != nil {
		return s.returnBadRequest(c, err)
	}

	tok, err := s.clients.Auth(c.Request().Context(), req.ClientID, req.ClientSecret)
	if err != nil {
		return s.returnJSONError(c, 401, spec.ErrInvalidClient, err)
	}
	return c.JSON(200, tok)
}

const (
	bearer = "Bearer"

	// 当鉴权成功时, 将会把请求者的client_id添加至context中对应的key
	contextClientIDKey = "authn.clientid"
)

// 鉴权middleware,
// 当校验成功时, 将会把请求者的client_id添加至context中
// 当校验失败时, 返回401
func (s *Server) authn() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get(echo.HeaderAuthorization)
			l := len(bearer)
			if len(auth) > l+1 && strings.EqualFold(auth[:l], bearer) {
				if tok := strings.TrimSpace(auth[l+1:]); tok != "" {
					clientid, err := s.clients.Verify(c.Request().Context(), tok)
					if err != nil {
						return s.returnJSONError(c, 401, spec.ErrInvalidToken, err)
					}

					// 把请求者的client_id添加至context中
					c.Set(contextClientIDKey, clientid)

					return next(c)
				}
			}

			return s.returnJSONError(c, 401, spec.ErrInvalidToken,
				fmt.Errorf("missing access_token in header Authorization: Bearer <your_access_token>"),
			)
		}
	}
}
