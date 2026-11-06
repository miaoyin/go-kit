package confutil

type StringReader interface {
	GetString(key string) string
}

type StringReaderFunc func(key string) string

func (f StringReaderFunc) GetString(key string) string {
	return f(key)
}
