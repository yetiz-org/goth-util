package hex

import "encoding/hex"

func EncodeToString(src []byte) string {
	return hex.EncodeToString(src)
}

func DecodeString(src string) []byte {
	if bytes, e := hex.DecodeString(src); e != nil {
		return nil
	} else {
		return bytes
	}
}
