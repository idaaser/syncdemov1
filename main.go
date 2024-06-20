// Package main 程序入口
package main

import (
	"crypto/rand"
	"crypto/rsa"

	"github.com/idaaser/syncdemov1/server"
)

func main() {
	srv := server.New(8000,
		server.WithContactFileStore(
			"./server/testdata/departments.json",
			"./server/testdata/users.json",
		),
	)
	/**
	server.WithJWTAuthnStore(
		generateRSAKey(),
		30*time.Minute,
		"client_id_1", "client_secret_1",
		"client_id_2", "client_secret_2",
	),
	*/

	srv.Start()
}

func generateRSAKey() *rsa.PrivateKey {
	k, _ := rsa.GenerateKey(rand.Reader, 3072)
	return k
}
