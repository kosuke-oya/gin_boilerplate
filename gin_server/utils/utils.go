package utils

import (
	rand "crypto/rand"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	Init()
}

func Init() {
	if os.Getenv("ENV") == "prod" {
		// 本番環境は.envファイルを読み込む
		err := godotenv.Load("/app/.env")
		if err != nil {
			panic(err)
		}
	} else {
		// local開発環境
		err := godotenv.Load("/app/dev.env")
		if err != nil {
			panic(err)
		}

	}
}

const alphanum = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func UniqueID(digit int) string {
	b := make([]byte, digit)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Sprintf("crypto/rand.Read error: %v", err))
	}
	for i, byt := range b {
		b[i] = alphanum[int(byt)%len(alphanum)]
	}
	return string(b)
}
