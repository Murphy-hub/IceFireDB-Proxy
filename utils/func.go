package utils

import (
	"log"
	"os"
	"runtime/debug"
	"strconv"
	"time"
)

func GoWithRecover(handler func(), recoverHandler func(r interface{})) {
	go func() {
		defer func() {
			if r := recover(); r != nil {

				log.Println("%s goroutine panic: %v\n%s\n", time.Now().Format("2006-01-02 15:04:05"), r, string(debug.Stack()))
				if recoverHandler != nil {
					go func() {
						defer func() {
							if p := recover(); p != nil {
								log.Println("recover goroutine panic:%v\n%s\n", p, string(debug.Stack()))
							}
						}()
						recoverHandler(r)
					}()
				}
			}
		}()
		handler()
	}()
}

func GetInterfaceString(param interface{}) string {
	switch param := param.(type) {
	case []byte:
		return string(param)
	case string:
		return param
	case int:
		return strconv.Itoa(param)
	case float64:
		return strconv.Itoa(int(param))
	}
	return ""
}

func InArray(in string, array []string) bool {
	for k := range array {
		if in == array[k] {
			return true
		}
	}
	return false
}

var _hostname string

func GetHostname() string {
	if _hostname != "" {
		return _hostname
	}

	hostname, err := os.Hostname()
	if err != nil {
		return ""
	}
	_hostname = hostname
	return _hostname
}
