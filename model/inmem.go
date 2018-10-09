package model

import (
	"strconv"
	"time"
)

type TrackId struct {
	Id int
}

type URL struct {
	Url string `json:"url"`
}

type Information struct {
	Uptime  string
	Info    string
	Version string
}

type Track struct {
	ID          int
	HDate       time.Time
	Pilot       string
	Glider      string
	GliderId    string
	TrackLength float64
}

/*
** getUptime updates uptime and formates it in ISO 8601 standard
 */
func GetUptime(t time.Time) (uptime string) {
	now := time.Now()
	newTime := now.Sub(t)
	hours := int(newTime.Hours())
	sek := strconv.Itoa(int(newTime.Seconds()) % 36000 % 60)
	min := strconv.Itoa(int(newTime.Minutes()) % 60)
	y, m, d := "0", "0", "0"

	// Setting the days correct
	if hours > 23 {
		d = strconv.Itoa(hours / 24)
		hours %= 24
	}
	days, _ := strconv.Atoi(d)
	// Setting the month correct
	if days > 31 {
		m = strconv.Itoa(days / 31)
		d = strconv.Itoa(days % 31)

	}
	months, _ := strconv.Atoi(m)
	// Setting the year correct
	if months > 12 {
		y = strconv.Itoa(months / 12)
		m = strconv.Itoa(months % 12)
	}

	hour := strconv.Itoa(hours)
	uptime = "P" + y + "Y" + m + "M" + d + "DT" + hour + "H" + min + "M" + sek + "S"

	return uptime
}
