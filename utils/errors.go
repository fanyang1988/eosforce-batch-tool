package utils

import (
	"fmt"
	"github.com/cihub/seelog"
	"os"
)

// ErrorCheck if err is not nil exit and print err
func ErrorCheck(err error, prefix_format string, a ...interface{}) {
	if err != nil {
		seelog.Flush()
		prefix := fmt.Sprintf(prefix_format, a...)
		fmt.Printf("ERROR: %s: %s\n", prefix, err)
		os.Exit(1)
	}
}
