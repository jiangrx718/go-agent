package gins

import jsoniter "github.com/json-iterator/go"

func BytesToAny(bytes []byte) any {
	var a any
	if err := jsoniter.Unmarshal(bytes, &a); err != nil {
		return string(bytes)
	}
	return a
}
