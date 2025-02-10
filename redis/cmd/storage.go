package cmd

import "fmt"

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
		return strVal + "\n"
	}
	return fmt.Sprintf("%v", value) 
}

func (s *Storage) del (key string) string {
	delete(s.store, key)
	return "1"
}

func getStorage() *Storage {
	return &Storage{store: make(map[string]interface{})}
}