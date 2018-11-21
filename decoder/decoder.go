package decoder

type Decoder interface {
	Decode(interface{}) (interface{}, error)
}
