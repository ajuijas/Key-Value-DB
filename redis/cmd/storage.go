package cmd

import (
	"fmt"
	"strconv"
)

type Storage struct {
	store map[string]interface{}
}

func (s *Storage) set (args []string) string {
	switch len(args){
	case 1 : return "(error) ERR wrong number of arguments for 'set' command\n"
	case 2 : break
	default : return "(error) ERR syntax error\n"
	}
	s.store[args[0]] = args[1]
	return "OK\n"
}

func (s *Storage) get (args []string) string {
	if len(args) != 1 {
		return "(error) ERR wrong number of arguments for 'get' command\n"
	}
	value, exists := s.store[args[0]]
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