package datamodel

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Message a struct
type Message struct {
	Topic string   `json:"topic"`
	Data  KeyValue `json:"data"`
}
