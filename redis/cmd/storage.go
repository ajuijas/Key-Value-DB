package cmd

import (
	"fmt"
	"strconv"
	"sync"
)

type Storage struct {
	store map[string]interface{}
	mutex sync.Mutex
}

func (s *Storage) set (args []string) string {
	switch len(args){
	case 1 : return "(error) ERR wrong number of arguments for 'set' command\n"  //TODO: Proper error handling mechanism
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

func (s *Storage) incr (args []string) string {

	if len(args) != 1 {
		return "(error) ERR wrong number of arguments for 'incr' command\n"
	}

	key := args[0]
	value, exists := s.store[key]
	if !exists {
		s.store[key] = 1
		return "(integer) 1\n"
	}
	valueInt, err := strconv.Atoi(fmt.Sprintf("%v", value))  // TODO: check if this handle all cases
	if err!=nil {
		return "(error) ERR value is not an integer or out of range\n"
	}else{
		s.store[key] = valueInt + 1
		return fmt.Sprintf("(integer) %v\n", s.store[key])
	}
}

func (s *Storage) incrby (args []string) string {
	if len(args) !=2 {
		return "(error) ERR wrong number of arguments for 'incrby' command\n"
	}
	key := args[0]
	incrby, err := strconv.Atoi(args[1])

	if err!=nil {
		return "(error) ERR value is not an integer or out of range\n"
	}

	value, exists := s.store[key]
	if !exists {
		s.store[key] = incrby
		return fmt.Sprintf("(integer) %v\n", s.store[key])
	}
	valueInt, err := strconv.Atoi(fmt.Sprintf("%v", value))  // TODO: check if this handle all cases
	if err!=nil {
		return "(error) ERR value is not an integer or out of range\n"
	}else{
		s.store[key] = valueInt + incrby
		return fmt.Sprintf("(integer) %v\n", s.store[key])
	}

}

func getStorage() *Storage {
	return &Storage{store: make(map[string]interface{})}
}

func wrapString(str string) string {
	return "\"" + str + "\""
}