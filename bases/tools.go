package bases

import (
	"encoding/json"
	"fmt"
)

func Recover() {
	if r := recover(); r != nil {
		if err, ok := r.(error); !ok {
			fmt.Print(err)
		}
	}
}

func Entity2Map(e interface{}) (map[string]interface{}, error) {
	f, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	var mp map[string]interface{}
	if err1 := json.Unmarshal(f, &mp); err1 != nil {
		return nil, err1
	}
	return mp, nil
}
