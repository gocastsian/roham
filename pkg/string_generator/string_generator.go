package stringgenerator

import (
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomString(n int64) string {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[r.Intn(len(letterBytes))]
    }

    return string(b)
}