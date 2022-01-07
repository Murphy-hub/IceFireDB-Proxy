package router

import (
	"math"

	"github.com/IceFireDB/IceFireDB-Proxy/pkg/RedSHandle"
)

const AbortIndex int8 = math.MaxInt8 / 2

type Context struct {
	Writer *RedSHandle.WriterHandle
	// Args     [][]byte
	Args     []interface{}
	Cmd      string
	Handlers HandlersChain // Middleware and final handler functions
	Index    int8
	Op       OpFlag
	Reply    interface{}
}

func (c *Context) Reset() {
	c.Writer = nil
	c.Args = c.Args[0:0]
	c.Handlers = nil
	c.Index = -1
	c.Reply = nil
}

/************************************/
/*********** FLOW CONTROL ***********/
/************************************/

// Next should be used only inside middleware.
// It executes the pending Handlers in the chain inside the calling handler.
// See example in GitHub.
func (c *Context) Next() error {
	c.Index++
	for c.Index < int8(len(c.Handlers)) {
		err := c.Handlers[c.Index](c)
		if err != nil {
			return err
		}
		c.Index++
	}
	return nil
}

// IsAborted returns true if the current context was aborted.
func (c *Context) IsAborted() bool {
	return c.Index >= AbortIndex
}

// Abort prevents pending Handlers from being called. Note that this will not stop the current handler.
// Let's say you have an authorization middleware that validates that the current request is authorized.
// If the authorization fails (ex: the password does not match), call Abort to ensure the remaining Handlers
// for this request are not called.
func (c *Context) Abort() {
	c.Index = AbortIndex
}
