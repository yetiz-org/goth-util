package kkutil

import (
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

var zeroTime = time.Unix(0, 0)

func BytesFromUUID(uuidObj uuid.UUID) []byte {
	bytes, _ := uuidObj.MarshalBinary()
	return bytes
}

func ZeroTimestamp() int64 {
	return 0
}

func MaxTimestamp() int64 {
	return 4294967295
}

func MaxTimestamp32u() uint {
	return 4294967295
}

func UnixToTime(unixSecond int64) time.Time {
	return time.Unix(unixSecond, 0)
}

func IsInt(str string) bool {
	if _, err := strconv.ParseInt(str, 10, 64); err != nil {
		return false
	}

	return true
}

func CastString(obj interface{}) *string {
	if str, ok := obj.(string); ok {
		return &str
	}

	return nil
}

func NowUInt() uint {
	return uint(time.Now().Unix())
}

func SplitRemoteAddr(addr string) (ip net.IP, port string) {
	if strings.Count(addr, ":") > 1 {
		if strings.LastIndex(addr, "]") == -1 {
			return net.ParseIP(addr), ""
		}

		ip := net.ParseIP(addr[1 : strings.LastIndex(addr, ":")-1])
		port := addr[strings.LastIndex(addr, ":")+1:]
		return ip, port
	} else {
		if strings.LastIndex(addr, ":") == -1 {
			return net.ParseIP(addr), ""
		}

		ip := net.ParseIP(addr[:strings.LastIndex(addr, ":")])
		port := addr[strings.LastIndex(addr, ":")+1:]
		return ip, port
	}
}
