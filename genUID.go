package main

import (
	"crypto/rand"
	"encoding/base64"
)

func genUID(length int) string {
	v := make([]byte, length)
	rand.Read(v)
	return base64.URLEncoding.EncodeToString(v)[:length]
}
