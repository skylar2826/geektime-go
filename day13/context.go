package rpc

import "context"

var oneWayKey struct{}

func OneWayContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, oneWayKey, true)
}

func isOneWay(ctx context.Context) bool {
	oneWay, ok := ctx.Value(oneWayKey).(bool)
	if ok && oneWay == true {
		return true
	}
	return false
}
