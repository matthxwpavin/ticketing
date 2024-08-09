package iferr

import (
	"context"
	"fmt"

	"github.com/matthxwpavin/ticketing/logging/sugar"
)

func Log(ctx context.Context, err error, msgAndArgs ...any) {
	if err != nil {
		sugar.FromContext(ctx).Errorln(messageFromMsgAndArgs(msgAndArgs...))
	}
}

func messageFromMsgAndArgs(msgAndArgs ...interface{}) string {
	if len(msgAndArgs) == 0 || msgAndArgs == nil {
		return ""
	}
	if len(msgAndArgs) == 1 {
		msg := msgAndArgs[0]
		if msgAsStr, ok := msg.(string); ok {
			return msgAsStr
		}
		return fmt.Sprintf("%+v", msg)
	}
	if len(msgAndArgs) > 1 {
		return fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	}
	return ""
}
