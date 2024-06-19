package server

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net/http"
	"strings"
	"time"

	spec "github.com/idaaser/syncspecv1"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

// WithTokenStore 使用自定义的TokenStore
func WithTokenStore(store TokenStore) Option {
	return func(srv *Server) {
		srv.tokens = store
	}
}

// WithRS256JWTTokenStore 使用JWT格式的token(使用RSA格式的私钥, 私钥长度建议>=2048)
func WithRS256JWTTokenStore(key *rsa.PrivateKey, ttl time.Duration) Option {
	return WithTokenStore(&jwtTokenStore{key: key, ttl: ttl})
}

// TokenStore 定义token相关的接口
type TokenStore interface {
	Create(context.Context) (*spec.Token, error)
	Verify(context.Context, string) error
}

type anyTokenStore struct{}

// Create 实现了TokenStore接口
func (t *anyTokenStore) Create(context.Context) (*spec.Token, error) {
	return &spec.Token{AccessToken: "access_token", ExpiresIn: 300}, nil
}

// Verify 实现了TokenStore接口
func (t *anyTokenStore) Verify(context.Context, string) error {
	return nil
}

type jwtTokenStore struct {
	key any
	ttl time.Duration
}

// Create 实现了TokenStore接口, 生成一个token
func (t *jwtTokenStore) Create(context.Context) (*spec.Token, error) {
	token := jwt.New()
	_ = token.Set(jwt.IssuedAtKey, time.Now().Unix())
	_ = token.Set(jwt.ExpirationKey, time.Now().Add(t.ttl).Unix())

	_ = token.Set("spec", "v1")

	tok, err := jwt.Sign(token,
		jwt.WithKey(jwa.RS256, t.key),
	)
	if err != nil {
		return nil, err
	}

	return &spec.Token{
		AccessToken: string(tok),
		ExpiresIn:   int32(t.ttl.Seconds()),
	}, nil
}

// Verify 实现了TokenStore接口, 验证token
func (t *jwtTokenStore) Verify(_ context.Context, tok string) error {
	_, err := jwt.Parse([]byte(tok),
		jwt.WithKey(jwa.RS256, t.key),
		jwt.WithAcceptableSkew(2*time.Minute),
		jwt.WithClaimValue("spec", "v1"),
	)
	if err != nil {
		return err
	}

	return nil
}

type tokenRequest struct {
	ClientID     string `form:"client_id"`
	ClientSecret string `form:"client_secret"`
}

func (req tokenRequest) validate() error {
	if req.ClientID == "" {
		return fmt.Errorf("client id MUST NOT be empty")
	}

	if req.ClientSecret == "" {
		return fmt.Errorf("client secret MUST NOT be empty")
	}

	return nil
}

func (s *Server) token(c echo.Context) error {
	req := tokenRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			spec.ErrResponse{Error: "invalid_request", ErrorMessage: err.Error()})
	}

	if clientid, clientsecret, ok := c.Request().BasicAuth(); ok {
		req.ClientID = clientid
		req.ClientSecret = clientsecret
	}

	if err := req.validate(); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			spec.ErrResponse{Error: "invalid_request", ErrorMessage: err.Error()})
	}

	if err := s.clients.Verify(req.ClientID, req.ClientSecret); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			spec.ErrResponse{Error: "invalid_request", ErrorMessage: err.Error()})
	}

	tok, err := s.tokens.Create(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			spec.ErrResponse{Error: "internal_error", ErrorMessage: err.Error()})
	}
	return c.JSON(201, tok)
}

const (
	bearer = "Bearer"
)

// authentication middleware
func (s *Server) authn() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get(echo.HeaderAuthorization)
			l := len(bearer)
			if len(auth) > l+1 && strings.EqualFold(auth[:l], bearer) {
				if tok := strings.TrimSpace(auth[l+1:]); tok != "" {
					if err := s.tokens.Verify(c.Request().Context(), tok); err != nil {
						return c.JSON(
							http.StatusUnauthorized,
							spec.ErrResponse{
								Error:        "invalid_token",
								ErrorMessage: err.Error(),
							},
						)
					}

					return next(c)
				}
			}

			return c.JSON(
				http.StatusUnauthorized,
				spec.ErrResponse{
					Error:        "invalid_token",
					ErrorMessage: "missing access_token in header Authorization: Bearer <your_access_token>",
				},
			)
		}
	}
}
