package lib

type InferType struct {
	isInt    bool
	isBool   bool
	isString bool
}

func NewInfer(i interface{}) *InferType {
	t := InferType{}
	switch i.(type) {
	case int:
		t.isInt = true
	case string:
		t.isString = true
	case bool:
		t.isBool = true
	}
	return &t
}
