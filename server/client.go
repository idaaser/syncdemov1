package server

import "fmt"

// WithClient 使用自定义的ClientStore
func WithClient(store ClientStore) Option {
	return func(srv *Server) {
		srv.clients = store
	}
}

// WithMemoryClient 使用本地存储(内存)的一个ClientStore
func WithMemoryClient(clientid, clientsecret string) Option {
	return func(srv *Server) {
		if srv.clients == nil {
			WithClient(&memoryClientStore{
				clients: map[string]string{clientid: clientsecret},
			})(srv)
		} else {
			if m, ok := srv.clients.(*memoryClientStore); ok {
				m.add(clientid, clientsecret)
			}
		}
	}
}

// ClientStore 定义接口校验应用合法性
type ClientStore interface {
	Verify(clientid, clientsecret string) error
}

type memoryClientStore struct {
	clients map[string]string
}

func (s *memoryClientStore) add(clientid, clientsecret string) {
	s.clients[clientid] = clientsecret
}

func (s *memoryClientStore) Verify(clientid, clientsecret string) error {
	found, ok := s.clients[clientid]
	if !ok {
		return fmt.Errorf("invalid client id or client secret")
	}

	if found != clientsecret {
		return fmt.Errorf("invalid client id or client secret")
	}

	return nil
}
