package router

import "errors"

var (
	ErrLocalWriter  = errors.New("client conn writer error")
	ErrLocalFlush   = errors.New("client conn flush error")
	ErrCmdTypeWrong = errors.New("Err command type wrong")
)

const (
	ErrUnknownCommand = "ERR command resp type not support`%s`"
	ErrArguments      = "ERR wrong number of arguments for '%s' command"
)
