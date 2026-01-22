package gins

type CustomFields map[string]string

func (c CustomFields) Get(key string) string {
	return c[key]
}
