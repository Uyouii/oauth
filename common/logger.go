package common

import (
	"context"
	"log"
)

func GetErrorf(ctx context.Context) func(format string, fields ...interface{}) {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds | log.LstdFlags)
	return log.Printf
}

func GetInfof(ctx context.Context) func(format string, fields ...interface{}) {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds | log.LstdFlags)
	return log.Printf
}

func GetLogFuns(ctx context.Context) (func(format string, fields ...interface{}), func(format string, fields ...interface{})) {
	return GetInfof(ctx), GetErrorf(ctx)
}
