package entity

type KeyValue struct {
	Key   string
	Value string
	TTL   int
}

type HashField struct {
	Key   string
	Field string
	Value string
}
