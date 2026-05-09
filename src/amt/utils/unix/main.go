package unix

import (
	"errors"
	"syscall"
)

import (
	"amt/utils/print"
)

func IncreaseUlimit(batchSize uint64) {
	minimumUlimitValue := uint64(1024)

	if batchSize < minimumUlimitValue {
		print.Panic(errors.New("Batch size is too low"))
	}

	var rLimit syscall.Rlimit

	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)

	if err != nil {
		print.Panic(err)
	}

	if rLimit.Max >= batchSize {
		return
	}

	rLimit.Cur = minimumUlimitValue

	rLimit.Max = batchSize

	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)

	if err != nil {
		print.Panic(err)
	}
}