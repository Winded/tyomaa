package util

import (
	"strconv"
	"time"
)

var (
	defaultDateFormat = "02.01.2006 15:04:05"
)

func LeftPad(s string, pad string, plength int) string {
	for i := len(s); i < plength; i++ {
		s = pad + s
	}
	return s
}

func FormatDate(date time.Time) string {
	return date.Local().Format(defaultDateFormat)
}

func FormatDuration(duration time.Duration) string {
	duration = duration.Truncate(time.Second)
	numHours := int(duration.Hours())
	numMinutes := int(duration.Minutes()) % 60
	numSeconds := int(duration.Seconds()) % 60
	hour := LeftPad(strconv.Itoa(numHours), "0", 2)
	min := LeftPad(strconv.Itoa(numMinutes), "0", 2)
	sec := LeftPad(strconv.Itoa(numSeconds), "0", 2)
	return hour + ":" + min + ":" + sec
}
