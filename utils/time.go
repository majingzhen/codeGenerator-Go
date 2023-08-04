package utils

import "time"

func GetCurTimeStr() string {
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	return formattedTime
}

func Time2Str(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func Str2Time(str string) time.Time {
	layout := "2006-01-02 15:04:05"
	t, _ := time.Parse(layout, str)
	return t
}
