package value

import "encoding/json"

func JsonMarshal(obj interface{}) string {
	if marshal, err := json.Marshal(obj); err != nil {
		return ""
	} else {
		rtn := string(marshal)
		return rtn
	}
}
