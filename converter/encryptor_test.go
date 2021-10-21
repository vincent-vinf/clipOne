package converter

import (
	"clipOne/util"
	"encoding/base64"
	"fmt"
	"log"
)

import "testing"

func TestEncryptor_Encode(t *testing.T) {
	var c Converter = &BaseConverter{}
	e := &Encryptor{
		Converter: c,
	}
	keyBytes := []byte("2021clipOne")
	log.Println("keyBytes: ", keyBytes)

	keyList := util.MD5(keyBytes)
	log.Println("keyList: ", keyList)
	e.SetKey(keyList)
	c = e
	data, err := c.Encode([]byte("123hello"))
	if err != nil {
		return
	}
	log.Println(data)
	log.Println(len(data))
	sEnc := base64.StdEncoding.EncodeToString([]byte(data))
	fmt.Println(sEnc)
}
