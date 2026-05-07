package api

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var errF = errors.New("incorrect repeat format")

func mFunc(dstart, now time.Time, rules []string) (string, error) {
	var mounthRules []string
	mounthDayRules := strings.Split(rules[1], ",")
	if len(rules) > 2 {
		mounthRules = strings.Split(rules[2], ",")
	} else {
		mounthRules = strings.Split("1,2,3,4,5,6,7,8,9,10,11,12", ",")
	}
	dates := make([]string, 0, (len(mounthRules)*len(mounthDayRules))*2)
	for y := 0; y < 2; y++ {
		for i := 0; i <= len(mounthRules)-1; i++ {
			mounth, err := strconv.Atoi(mounthRules[i])
			if err != nil {
				return "", err
			}
			if mounth > 0 && mounth < 13 {
				for j := 0; j <= len(mounthDayRules)-1; j++ {
					day, err := strconv.Atoi(mounthDayRules[j])
					if err != nil {
						return "", err
					}
					if (day > -3 && day < 0) || (day > 0 && day < 32) {
						if day > 0 {
							t := time.Date(now.Year()+y, time.Month(mounth), day, 0, 0, 0, 0, time.UTC)
							if afterNow(t, dstart) && afterNow(t, now) && t.Day() == day {
								dates = append(dates, t.Format(DFormat))
							}
						} else {
							t := time.Date(now.Year()+y, time.Month(mounth+1), 1, 0, 0, 0, 0, time.UTC)
							res := t.AddDate(0, 0, day)
							if afterNow(res, dstart) && afterNow(res, now) {
								dates = append(dates, res.Format(DFormat))
							}
						}
					} else {
						return "", errF
					}
				}
			} else {
				return "", errF
			}
		}
	}
	d, err := time.Parse(DFormat, dates[0])
	if err != nil {
		return "", err
	}
	for _, v := range dates { // вычисляем близжайшую дату
		j, err := time.Parse(DFormat, v)
		if err != nil {
			return "", err
		}
		if j.Before(d) && afterNow(j, dstart) {
			d = j
		}

	}
	return d.Format(DFormat), nil
}

func wFunc(dstart, now time.Time, weekDays string) (string, error) {
	week := [7]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	if weekDays == "" {
		return "", errF
	}
	weekRules := strings.Split(weekDays, ",")
	date := dstart
	d := now
	for i := 0; i <= len(weekRules)-1; i++ { //вычисляем подходящие даты для каждого правила
		day, err := strconv.Atoi(weekRules[i])
		if err != nil {
			return "", err
		}
		if day > 0 && day <= 7 {
			date = dstart
			for {
				date = date.AddDate(0, 0, 1)
				if afterNow(date, now) && date.Weekday().String() == week[day-1] {
					weekRules[i] = date.Format(DFormat)
					break
				}
			}
		} else {
			return "", errF
		}
	}
	d, err := time.Parse(DFormat, weekRules[0])
	if err != nil {
		return "", err
	}
	for _, v := range weekRules { // вычисляем близжайшую дату
		j, err := time.Parse(DFormat, v)
		if err != nil {
			return "", err
		}
		if j.Before(d) && afterNow(j, dstart) {
			d = j
		}

	}
	return d.Format(DFormat), err
}

func dFunc(date, now time.Time, days string) (string, error) {
	if days == "" {
		return "", errF
	}
	interval, err := strconv.Atoi(days)
	if err != nil {
		return "", err
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
		return "", errF
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
	return date.After(now)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	date, err := time.Parse(DFormat, dstart)
	if err != nil {
		return "", err
	}
	rules := strings.Split(repeat, " ")

	switch rules[0] {
	case "d":
		if len(rules) > 1 {
			return dFunc(date, now, rules[1])
		} else {
			return "", errF
		}
	case "y":
		return yFunc(date, now), nil
	case "":
		return "", nil
	case "w":
		if len(rules) > 1 {
			return wFunc(date, now, rules[1])
		} else {
			return "", errF
		}
	case "m":
		if len(rules) > 1 {
			return mFunc(date, now, rules)
		} else {
			return "", errF
		}
	default:
		log.Println(errF)
		return "", errF
	}
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJson(w, map[string]string{"error": "Method Not Allowed"}, http.StatusMethodNotAllowed)
		return
	}
	now := time.Now()
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	sNow := r.URL.Query().Get("now")
	if len(sNow) > 0 {
		pNow, err := time.Parse(DFormat, sNow)
		now = pNow

		if err != nil {
			log.Println(err)
		}
	}
	dstart := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")
	answer, err := NextDate(now, dstart, repeat)
	if err != nil {
		log.Println(err)
	}
	_, err = w.Write([]byte(answer))
	if err != nil {
		log.Println(err)
	}
}
