package global

var Cleanups = make([]func() error, 0)
