package ulimit

import (
	"fmt"
	"syscall"
)

func IncreaseUlimit(batchSize uint64) {
	minimumUlimitValue := uint64(1024)

	if batchSize < minimumUlimitValue {
		panic(fmt.Sprintf("Batch size is too low, it has to be at least %d", minimumUlimitValue))
	}

	var rLimit syscall.Rlimit

	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)

	if err != nil {
		panic(fmt.Sprintf("Could not get ulimit value: %s", err.Error()))
	}

	if rLimit.Max >= batchSize {
		return
	}

	rLimit.Cur = minimumUlimitValue

	rLimit.Max = batchSize

	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)

	if err != nil {
		panic(fmt.Sprintf("Could not set a ulimit value: %s", err.Error()))
	}
}