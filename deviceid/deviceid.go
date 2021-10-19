package deviceid

import (
	"crypto/md5"
	"github.com/denisbrodbeck/machineid"
)

const AppName = "clipOne"

type DeviceID []byte

func GetDeviceID() (id DeviceID, err error) {
	source, err := machineid.ID()
	if err != nil {
		return nil, err
	}
	return MD5([]byte(source + AppName)), nil
}

func MD5(in []byte) []byte {
	md5Ctx := md5.New()
	md5Ctx.Write(in)
	return md5Ctx.Sum(nil)
}
