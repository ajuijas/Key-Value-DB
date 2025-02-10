package cmd

import (
	"fmt"
	"strconv"
)

type Storage struct {
	store map[string]interface{}
}

func (s *Storage) set (key string, value interface{}) string {
	s.store[key] = value
	return "OK\n"
}

func (s *Storage) get (key string) string {
	value, exists := s.store[key]
	if !exists {
		return "(nil)\n"
	}
	if strVal, ok := value.(string); ok {
		return wrapString(strVal) + "\n"
	}
	return fmt.Sprintf("%v", value) 
}

func (s *Storage) del (keys []string) string {
	var i int = 0
	for _, key := range keys{
		if s.store[key]!=nil{
			i = i +1
			delete(s.store, key)
		}
	}
	return "(integer) " + strconv.Itoa(i) + "\n"
}

func getStorage() *Storage {
	return &Storage{store: make(map[string]interface{})}
}

func wrapString(str string) string {
	return "\"" + str + "\""
}