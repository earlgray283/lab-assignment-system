package lib

import (
	"os"
	"time"
)

func IsDebug() bool {
	return os.Getenv("DEBUG") == "true"
}

func YMD(y, m, d int) time.Time {
	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local)
}
