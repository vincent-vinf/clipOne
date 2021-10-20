package converter

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"log"
)

type Compressor struct {
	Converter Converter
}

func (c *Compressor) Encode(in []byte) ([]byte, error) {
	log.Println("source size:", len(in))

	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err := zw.Write(in)
	if err != nil {
		return nil, err
	}

	err = zw.Close()
	if err != nil {
		return nil, err
	}
	log.Println("compress size:", buf.Len())
	return buf.Bytes(), nil
}

func (c *Compressor) Decode(in []byte) ([]byte, error) {
	zr, err := gzip.NewReader(bytes.NewBuffer(in))
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(zr)
	if err != nil {
		return nil, err
	}
	return data, nil
}
