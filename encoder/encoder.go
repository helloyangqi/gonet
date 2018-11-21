package encoder

type Encoder interface {
	Encode(in interface{}) (interface{}, error)
}
