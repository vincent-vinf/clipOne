package converter

type Converter interface {
	Encode([]byte) ([]byte, error)
	Decode([]byte) ([]byte, error)
}

type BaseConverter struct {
}

func (c *BaseConverter) Encode(data []byte) ([]byte, error) {
	return data, nil
}

func (c *BaseConverter) Decode(data []byte) ([]byte, error) {
	return data, nil

}
