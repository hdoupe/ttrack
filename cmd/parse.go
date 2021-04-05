package cmd

import (
	"errors"
	"fmt"
	"time"
)

// ParseTimeArg converts a string into a time.Time object.
func ParseTimeArg(arg string) (time.Time, error) {
	dateFmts := []string{
		"2006-01-02 3:04:05 PM MST",
		"01-02 3:04:05 PM MST",
		"2006-01-02 3:04 PM MST",
		"01-02 3:04 PM MST",
		time.Kitchen,
		"2006-01-02",
	}
	var res time.Time
	for _, dateFmt := range dateFmts {
		// Special update for kitchen time format to include year, month,
		// and day.
		if dateFmt == time.Kitchen {
			loc, _ := time.LoadLocation("Local")
			t, err := time.ParseInLocation(dateFmt, arg, loc)
			if err != nil {
				continue
			}
			res = t
			n := time.Now()
			// Month and day default values are already 1!
			res = res.AddDate(n.Year(), int(n.Month())-1, n.Day()-1)
			break
		} else {
			t, err := time.Parse(dateFmt, arg)
			if err == nil {
				res = t
				break
			}
		}
	}
	if res.IsZero() {
		msg := fmt.Sprintf("Unable to parse %s using formats %v", arg, dateFmts)
		return res, errors.New(msg)
	}
	if res.Year() == 0 {
		res = res.AddDate(time.Now().In(res.Location()).Year(), 0, 0)
	}
	return res.UTC(), nil
}
