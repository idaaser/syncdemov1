// Package main 程序入口
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"time"

	"github.com/idaaser/syncdemov1/server"
)

func main() {
	srv := server.New(8000,
		server.WithMemoryClient("client_id_1", "client_secret_1"),
		server.WithRS256JWTTokenStore(
			generateRSAKey(),
			2*time.Hour,
		),
	)
	srv.Start()
}

func generateRSAKey() *rsa.PrivateKey {
	k, _ := rsa.GenerateKey(rand.Reader, 3092)
	return k
}
