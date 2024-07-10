package memory

// Memory represents a list of maps with string keys and interface{} values
type Memory struct {
	data []map[string]interface{}
}

//create a new memory list
var memory Memory = Memory{data: []map[string]interface{}{}}

// Append adds a new key-value pair to the memory list and returns the index of the new pair. return bool
func Append(key string, value interface{}) bool {
	//if key already exists, Update()
	if Get(key) != nil {
		return Update(key, value)
	} else {
		memory.data = append(memory.data, map[string]interface{}{key: value})
		return true
	}
}

// Get retrieves the value associated with a key from the memory list. return only value
func Get(key string) interface{} {
	for _, item := range memory.data {
		if val, ok := item[key]; ok {
			return val
		}
	}
	return nil
}

// Delete removes a key-value pair from the memory list, return bool
func Delete(key string) bool {
	for i, item := range memory.data {
		if _, ok := item[key]; ok {
			memory.data = append(memory.data[:i], memory.data[i+1:]...)
			return true
		}
	}
	return false
}

//Update updates the value associated with a key in the memory list. return bool
func Update(key string, value interface{}) bool {
	for i, item := range memory.data {
		if _, ok := item[key]; ok {
			memory.data[i][key] = value
			return true
		}
	}
	return false
}
