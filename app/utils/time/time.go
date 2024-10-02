package time

import (
	"fmt"
	utils "framework/utils/common"
	"strconv"
	"strings"
	"time"
)

var tz = utils.GetEnv("DB_TIMEZONE", "Asia/Kuala_Lumpur")
var loc, _ = time.LoadLocation(tz)

func ConvertStringToTime(s string) time.Time {
	t, err := time.ParseInLocation("2006-01-02 15:04:05.999", s, loc)
	if err != nil {
		fmt.Printf("time conversion failed, s: %s\n", s)
		return time.Now()
	}
	return t
}

func GetStartOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, loc)
}

func GetEndOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 23, 59, 59, 999999999, loc)
}

// getStartDate determines the start date for a given time range based on the current time,
// optional query parameter, and specified path type.
//
// If a queryStart parameter is provided, it is used as the start date. If not, the start
// date is calculated based on the pathType:
// - For "weekly" path type, it returns the beginning of the current week.
// - For other path types, it returns the first day of the current month.
func GetStartDate(now time.Time, queryStart string, pathType string) string {
	if queryStart != "" {
		return queryStart
	}

	if pathType == "weekly" {
		return now.AddDate(0, 0, -int(now.Weekday())).Format("2006-01-02")
	}

	return now.AddDate(0, 0, -int(now.Day())+1).Format("2006-01-02")
}

// getEndDate determines the end date for a given time range based on the start date,
// optional query parameter, and specified path type.
//
// If a queryEnd parameter is provided, it is used as the end date. If not, the end date
// is calculated based on the pathType:
// - For "weekly" path type, it returns the end of the current week.
// - For other path types, it returns the last day of the current month.
func GetEndDate(now time.Time, start string, queryEnd string, pathType string) string {
	if queryEnd != "" {
		return queryEnd
	}

	startTime, _ := time.Parse("2006-01-02", start)

	if pathType == "weekly" {
		return startTime.AddDate(0, 0, 6).Format("2006-01-02")
	}

	return startTime.AddDate(0, 1, -1).Format("2006-01-02")
}

// parseDurationToBigInt takes a duration string in the format "DDD HH:mm:ss.SSS" and
// converts it to a bigint value representing the duration in milliseconds. The function
// is designed to parse duration strings that may include days, hours, minutes, seconds,
// and milliseconds.
//
// Parameters:
//   - input (string): The duration string to be parsed. Example: "000 00:00:00.010"
//
// Returns:
//   - int64: The duration represented as a bigint value in milliseconds.
//   - error: An error is returned if the input duration string is in an invalid format
//            or if there are parsing errors.
//
// Usage:
//   input := "000 00:00:00.010"
//   durationAsBigInt, err := parseDurationToBigInt(input)
//   if err != nil {
//       log.Fatal(err)
//   }
//   fmt.Printf("Duration as BigInt: %d\n", durationAsBigInt)
func ParseDurationToBigInt(input string) (int64, error) {
	parts := strings.Fields(input)

	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid duration format")
	}

	days, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("failed to parse days: %v", err)
	}

	// Format the time part by removing the leading zero
	timePart := strings.TrimLeft(parts[1]+parts[2], "0")
	duration, err := time.ParseDuration(timePart)
	if err != nil {
		return 0, fmt.Errorf("failed to parse time duration: %v", err)
	}

	// Convert days to hours and add to the time duration
	duration += time.Duration(days) * 24 * time.Hour

	// Convert the duration to milliseconds
	durationMilliseconds := duration.Milliseconds()

	return durationMilliseconds, nil
}
