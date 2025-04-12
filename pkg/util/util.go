package util

import (
	"fmt"
	"strings"
	"time"
)

func ValidateDates(dates string) bool {
	dateArr := strings.Split(dates, ":")
	fmt.Println(dateArr)
	const layout = "20060102"
	for _, date := range dateArr {
		if _, err := time.Parse(layout, date); err != nil {
			return false
		}
	}
	return true
}
