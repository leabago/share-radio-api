package chronos

import (
	"fmt"
	"strconv"
	"time"
)

//---------------------------------------- Helper functions ----------------------------------------

// Since wrap time Since.
func Since(t Chronos) time.Duration {
	if t.IsZero() {
		return 0
	}

	return time.Since(t.Time)
}

// Second returns a time.Duration equal to the number of seconds.
func Second(seconds int) time.Duration {
	return time.Duration(seconds) * time.Second
}

// Minute returns a time.Duration equal to the number of minutes.
func Minute(minutes int) time.Duration {
	return time.Duration(minutes) * time.Minute
}

// Hour returns a time.Duration equal to the number of hours.
func Hour(hours int) time.Duration {
	return time.Duration(hours) * time.Hour
}

// HourInt32 returns a time.Duration equal to the number of hours.
func HourInt32(hours int32) time.Duration {
	return time.Duration(hours) * time.Hour
}

// UnixNanoString return string time.Now in UnixNano format.
func UnixNanoString() string {
	return strconv.FormatInt(Now().UnixNano(), 10)
}

// UnixNano return int64 time.Now in UnixNano format.
func UnixNano() int64 {
	return Now().UnixNano()
}

// NewTickerSecond wrap time.NewTicker.
func NewTickerSecond(seconds int) *time.Ticker {
	ticker := time.NewTicker(Second(seconds))

	return ticker
}

// ParseTimeWithFormats try convert to proposed formats.
func ParseTimeWithFormats(key, value string, formats []Format) (Chronos, error) {
	for _, f := range formats {
		t, err := time.Parse(string(f), value)
		if err == nil {
			return Chronos{t}, nil
		}
	}

	return Chronos{}, fmt.Errorf("%w: parameter: %s, expected: %q, got: %q", errInvalidTime, key, formats, value)
}

// ParseAndFormat parses a time string and return formated string.
func ParseAndFormat(timeString string, format Format) (string, error) {
	if timeString == "" {
		return "", nil
	}

	ch, err := FromString(timeString)
	if err != nil {
		return "", fmt.Errorf("parse time %q: %w", timeString, err)
	}

	return ch.Format(format), nil
}

// ParseAndFormatLoc parses a time string and formats it into the specified format.
func ParseAndFormatLoc(timeString string, format Format, loc *time.Location) string {
	if timeString == "" {
		return ""
	}

	dateFrom, err := FromString(timeString)
	if err != nil {
		return ""
	}

	return dateFrom.FormatLoc(format, loc)
}

// HoursBetween - returns the number of hours between two dates and round up it.
func HoursBetween(start, end Chronos) int {
	return start.HoursBetween(end)
}

// HoursBetweenFloat32 - returns the number of hours between two dates and round up it.
func HoursBetweenFloat32(start, end Chronos) float32 {
	return float32(HoursBetween(start, end))
}

// DaysBetween returns days between two Chronos.
func DaysBetween(start, end Chronos) int {
	return start.DaysBetween(end)
}

// ParseSameFormat parses a time string in specified format and return formated string.
// Ignores parse errors, returning an empty string if there is an error.
func ParseSameFormat(format Format, timeString string) string {
	ch, _ := Parse(format, timeString)

	return ch.Format(format)
}

// IsTodayParse parses a time string and compare with current time Now().
func IsTodayParse(timeStr string) bool {
	ch, err := FromString(timeStr)
	if err != nil {
		return false
	}

	return Now().IsSameDate(ch)
}

// GetRussianMonths returns months in Russian.
func GetRussianMonths() [12]string {
	return [12]string{
		"января", "февраля", "марта", "апреля", "мая", "июня",
		"июля", "августа", "сентября", "октября", "ноября", "декабря",
	}
}

// FormatDateRu - return date in format like "1 января".
func FormatDateRu(t Chronos) string {
	months := GetRussianMonths()
	day := t.Day()
	month := months[int(t.Month())-1]

	return fmt.Sprintf("%02d %s", day, month)
}

// ChronosToStringPtr return pointer to string.
func ChronosToStringPtr(t *Chronos) *string {
	if t == nil {
		return nil
	}

	s := t.String()

	return &s
}

// Date returns the Chronos corresponding to yyyy-mm-dd hh:mm:ss + nsec nanoseconds.
func Date(year int, month time.Month, day, hour, minutes, sec, nsec int, loc *time.Location) Chronos {
	return Chronos{time.Date(year, month, day, hour, minutes, sec, nsec, loc)}
}

// Until returns the duration until t.
func Until(t Chronos) time.Duration {
	return time.Until(t.Time)
}

// StartOfDay returns the start of the day (00:00:00).
func StartOfDay(c Chronos) Chronos {
	if c.IsZero() {
		return c
	}

	return Chronos{time.Date(c.Year(), c.Month(), c.Day(), 0, 0, 0, 0, c.Location())}
}

// StartOfMonth returns the first day of the month.
func StartOfMonth(c Chronos) Chronos {
	return Chronos{time.Date(c.Year(), c.Month(), 1, 0, 0, 0, 0, c.Location())}
}

// EndOfDay returns the end of the day (23:59:59.999999999).
func EndOfDay(c Chronos) Chronos {
	if c.IsZero() {
		return c
	}

	return Chronos{time.Date(c.Year(), c.Month(), c.Day(), 23, 59, 59, 999999999, c.Location())}
}

// AddHours adds hours to time.
func AddHours(c Chronos, hours int) Chronos {
	if c.IsZero() {
		return c
	}

	return Chronos{Time: c.Time.Add(time.Duration(hours) * time.Hour)}
}

// IsSameDateParse parse and compares two dates, return true if equal.
func IsSameDateParse(chr Chronos, other string) (bool, error) {
	otherDate, err := FromString(other)
	if err != nil {
		return false, err
	}

	// Compare dates
	return chr.IsSameDate(otherDate), nil
}

func WeekdayRuIndex(day string) int {
	switch day {
	case "пн":
		return Monday
	case "вт":
		return Tuesday
	case "ср":
		return Wednesday
	case "чт":
		return Thursday
	case "пт":
		return Friday
	case "сб":
		return Saturday
	case "вс":
		return Sunday
	}

	return -1
}

// CreateZoneFromOffset - create timezone from offset.
func CreateZoneFromOffset(offset int) *time.Location {
	const secondsInHour = 3600

	loc := time.FixedZone(
		fmt.Sprintf("UTC%+d", offset),
		offset*secondsInHour,
	)

	return loc
}

// SameChronosInUTC - заменяет часовой пояс на UTC (Z) без изменения отображаемого времени.
// Оставляет те же значения года, месяца, дня, часа, минуты, секунды и наносекунды,
// но устанавливает зону UTC вместо исходной.
func SameChronosInUTC(ch Chronos) Chronos {
	utcTime := time.Date(
		ch.Year(), ch.Month(), ch.Day(),
		ch.Hour(), ch.Minute(), ch.Second(), ch.Nanosecond(),
		time.UTC,
	)

	return FromTime(utcTime)
}

// SameChronosInLoc - изменяет часовой пояс/смещение БЕЗ изменения времени.
// Пример: 2026-06-24 18:00:00 UTC → 2026-06-24 18:00:00 +03:00.
func SameChronosInLoc(chr Chronos, offset int) Chronos {
	loc := CreateZoneFromOffset(offset)
	newTime := time.Date(
		chr.Year(), chr.Month(), chr.Day(),
		chr.Hour(), chr.Minute(), chr.Second(), chr.Nanosecond(),
		loc,
	)

	return FromTime(newTime)
}
