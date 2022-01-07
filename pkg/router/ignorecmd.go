package router

import (
	"fmt"

	"github.com/IceFireDB/IceFireDB-Proxy/utils"
)

// Start the middleware according to the configuration
func IgnoreCMDMiddleware(enable bool, cmdList []string) HandlerFunc {
	return func(context *Context) error {
		// Ignore custom commands
		if enable && len(cmdList) > 0 {
			if utils.InArray(context.Cmd, cmdList) {
				// 被限制
				return fmt.Errorf(ErrUnknownCommand, context.Cmd)
			}
		}
		return context.Next()
	}
}
