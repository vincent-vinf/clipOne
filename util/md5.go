package util

import (
	"crypto/md5"
)

func MD5(in []byte) []byte {
	md5Ctx := md5.New()
	md5Ctx.Write(in)
	return md5Ctx.Sum(nil)
}
