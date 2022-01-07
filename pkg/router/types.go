package router

import "github.com/IceFireDB/IceFireDB-Proxy/pkg/RedSHandle"

type IRoutes interface {
	Use(...HandlerFunc) IRoutes
	AddCommand(string, ...HandlerFunc) IRoutes
	InitCMD()
	Handle(w *RedSHandle.WriterHandle, args []interface{}) error
	Close() error
}

type HandlerFunc func(*Context) error

// HandlersChain defines a HandlerFunc array.
type HandlersChain []HandlerFunc

// Last returns the last handler in the chain. ie. the last handler is the main one.
func (c HandlersChain) Last() HandlerFunc {
	if length := len(c); length > 0 {
		return c[length-1]
	}
	return nil
}
