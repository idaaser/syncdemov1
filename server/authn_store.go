package server

import (
	"context"
	"crypto/rsa"
	"fmt"
	"time"

	spec "github.com/idaaser/syncspecv1"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

// WithAuthnStore 使用自定义的AuthnStore
func WithAuthnStore(store AuthnStore) Option {
	return func(srv *Server) {
		srv.clients = store
	}
}

// WithJWTAuthnStore 使用JWT token的AuthnStore, 用内存来管理client_id/client_secret, 以及使用JWT的token来鉴权
// 注: 使用RSA格式的私钥来签发鉴权token, 私钥长度建议>=2048
func WithJWTAuthnStore(key *rsa.PrivateKey, exp time.Duration, clientIDAndSecrets ...string) Option {
	store := &jwtAuthnStore{
		clients: map[string]string{},
		key:     key, exp: exp,
	}
	store.addClient(clientIDAndSecrets...)

	return WithAuthnStore(store)
}

// AuthnStore 定义鉴权相关接口, 包括生成token以及校验token
type AuthnStore interface {
	// Auth 给发起请求的client_id/client_secret, 颁发access_token
	Auth(ctx context.Context, clientid, clientsecret string) (*spec.Token, error)

	// Verify 校验token的合法性, 若校验成功返回token对应的clientid
	Verify(ctx context.Context, tok string) (string, error)
}

// 允许任何access_token. 注: 仅用于测试
type allowAnyAs struct{}

func (s *allowAnyAs) Auth(ctx context.Context, clientid, clientsecret string) (*spec.Token, error) {
	return &spec.Token{AccessToken: "any token", ExpiresIn: 7200}, nil
}

func (s *allowAnyAs) Verify(ctx context.Context, tok string) (string, error) {
	return "any_client", nil
}

// 1. 使用map来维护client_id/client_secret
// 2. 使用JWT token(RS256签名算法)
type jwtAuthnStore struct {
	clients map[string]string

	key any
	exp time.Duration
}

// Auth 实现了AuthnStore接口, 生成一个token
func (s *jwtAuthnStore) Auth(ctx context.Context, clientid, clientsecret string) (*spec.Token, error) {
	if err := s.verifyClient(ctx, clientid, clientsecret); err != nil {
		return nil, err
	}

	return s.issueToken(ctx, clientid)
}

func (s *jwtAuthnStore) Verify(ctx context.Context, tok string) (string, error) {
	token, err := jwt.Parse([]byte(tok),
		jwt.WithKey(jwa.RS256, s.key),
		jwt.WithAcceptableSkew(2*time.Minute),
		jwt.WithClaimValue("spec", "v1"),
	)
	if err != nil {
		return "", err
	}

	sub := token.Subject()
	if _, found := s.clients[sub]; !found {
		return "", fmt.Errorf("invalid client_id %q", sub)
	}

	return sub, nil
}

// verifyClient 校验 client_id/client_secret的合法性
func (s *jwtAuthnStore) verifyClient(_ context.Context, clientid, clientsecret string) error {
	found, ok := s.clients[clientid]
	if !ok {
		return fmt.Errorf("invalid client id or client secret")
	}

	if found != clientsecret {
		return fmt.Errorf("invalid client id or client secret")
	}

	return nil
}

func (s *jwtAuthnStore) issueToken(_ context.Context, clientid string) (*spec.Token, error) {
	token := jwt.New()
	_ = token.Set(jwt.SubjectKey, clientid)
	_ = token.Set("spec", "v1")
	_ = token.Set(jwt.IssuedAtKey, time.Now().Unix())
	_ = token.Set(jwt.NotBeforeKey, time.Now().Unix())
	_ = token.Set(jwt.ExpirationKey, time.Now().Add(s.exp).Unix())

	tok, err := jwt.Sign(token,
		jwt.WithKey(jwa.RS256, s.key),
	)
	if err != nil {
		return nil, err
	}

	return &spec.Token{
		AccessToken: string(tok),
		ExpiresIn:   int32(s.exp.Seconds()),
	}, nil
}

func (s *jwtAuthnStore) addClient(clientIDAndSecrets ...string) {
	l := len(clientIDAndSecrets)
	if l%2 == 0 {
		for i := 0; i < l; i += 2 {
			id, secret := clientIDAndSecrets[i], clientIDAndSecrets[i+1]
			s.clients[id] = secret
		}
	}
}
