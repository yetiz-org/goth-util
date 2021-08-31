package structs

type Queue interface {
	Pop() interface{}
	Push(v interface{})
}
