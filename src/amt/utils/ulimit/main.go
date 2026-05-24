package ulimit

import (
	"syscall"
)

func Increase(batchSize uint64) {
	var rLimit syscall.Rlimit

	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)

	if err != nil {
		panic("Could not get ulimit value: " + err.Error())
	}

	if rLimit.Max >= batchSize {
		return
	}

	rLimit.Max = batchSize

	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)

	if err != nil {
		panic("Could not increase ulimit value: " + err.Error())
	}
}