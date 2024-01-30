package log

import (
	"fmt"
	"time"
)

func Printf(str string, a ...interface{}) {
	fmt.Printf("%s: %s\n", time.Now().Format("15:04:05"), fmt.Sprintf(str, a...))
}
