package converter

type Converter interface {
	Encode([]byte) ([]byte, error)
	Decode([]byte) ([]byte, error)
}

