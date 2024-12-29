package pkg

import (
	"context"
	"time"
)

func TimeContext(timeSecond ...int) (context.Context, context.CancelFunc) {

	var second int

	if len(timeSecond) > 0 {
		second = timeSecond[0]
	}
	if second == 0 {
		second = 10
	}

	c, cFunc := context.WithTimeout(context.Background(), time.Duration(second)*time.Second)

	return c, cFunc
}
