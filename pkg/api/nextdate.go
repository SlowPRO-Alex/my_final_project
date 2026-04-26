package api

import (
	"strconv"
	"time"
	"strings"
	"errors"
)

var fError = errors.New("указан неверный формат repeat")

func dFunc(date, now time.Time, days string) (string, error) {
	if days == "" {
		return "", fError
	} 
	interval, err := strconv.Atoi(days)
		if err != nil {
			return "",  err
		}
		if interval > 0 && interval <= 400 {
			for {
    			date = date.AddDate(0, 0, interval)
    			if afterNow(date, now) {
        			break
    			}
			} 
			return date.Format(DFormat), err
		} else {
				return "", fError
		}
}

func yFunc(date, now time.Time) string {
	for {
    		date = date.AddDate(1, 0, 0)
    		if afterNow(date, now) {
        		break
    		}
	} 
	return date.Format(DFormat)
}

func afterNow(date, now time.Time) bool {
	if  date.After(now) {
		return true
	} else {
		return false
	}
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	date, err := time.Parse(DFormat, dstart)
	if err != nil {
		return "",  err
	}
	rules := strings.Split(repeat, " ")
	
	switch(rules[0]) {
		case "":
			return "",  fError
		case "d":
			if len(rules) > 1 {
			return dFunc(date, now, rules[1])
			} else {
				return "",  fError
			}
		case "y":
			return yFunc(date, now), nil
		default:
			return "", fError
	}
	return date.Format(DFormat),  err
}