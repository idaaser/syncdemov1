// Package main 程序入口
package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/idaaser/syncdemov1/server"
	"github.com/lestrrat-go/jwx/v3/jwk"
)

func main() {
	srv := server.New(
		8001,
		/**
		server.WithJWTAuthnStore(
			newECKey(),
			120*time.Minute,
			"client_id_1", "client_secret_1",
			"client_id_2", "client_secret_2",
		),
		*/
		server.WithContactFileStore(
			"./server/testdata/departments.json",
			"./server/testdata/users.json",
			"./server/testdata/groups.json",
			"./server/testdata/group-users.json",
		),
	)

	srv.Start()
}

func generateRSAKey() *rsa.PrivateKey {
	k, _ := rsa.GenerateKey(rand.Reader, 3072)
	return k
}

func newECKey() jwk.Key {
	raw, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	k, err := jwk.Import(raw)
	if err != nil {
		panic(err)
	}

	if err := k.Set(jwk.AlgorithmKey, "ES256"); err != nil {
		panic(err)
	}
	if err := k.Set(jwk.KeyIDKey, "kid-"+strconv.Itoa(time.Now().Nanosecond())); err != nil {
		panic(err)
	}
	if err := k.Set(jwk.KeyUsageKey, "sig"); err != nil {
		panic(err)
	}

	b, _ := json.MarshalIndent(k, "", " ")
	fmt.Println(string(b))

	return k
}
