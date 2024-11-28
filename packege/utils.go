package utils

import (
	"fmt"
	"runtime"
)

func GetCallerInfo() string {
	pc, _, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	return fmt.Sprintf("%s:%d", fn.Name(), line)
}
