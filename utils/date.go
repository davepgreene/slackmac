package utils

import (
	"strconv"
	"time"
)

func EpochStringToTime(epochTime string) (time.Time, error) {
	msInt, err := strconv.ParseInt(epochTime, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(msInt, 0), nil
}