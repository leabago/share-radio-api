package chronos

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/goodsign/monday"
)

// Error messages.

var errUnsupportedType = errors.New("unsupported type for Chronos")
var errInvalidTime = errors.New("parameter has invalid time format")
var errUnableParse = errors.New("unable to parse time")
var errEmptyString = errors.New("received an empty string instead of time")

const (
	// Values used in date and time calculations.

	OneDayHours = 24
	Monday      = 0
	Tuesday     = 1
	Wednesday   = 2
	Thursday    = 3
	Friday      = 4
	Saturday    = 5
	Sunday      = 6
	DaysInWeek  = 7

	MinutesInHours = 60
)

// Chronos wraps time.Time and provides extended helper methods for parsing,
// formatting, and manipulating dates and times.
// RFC3339 - used as the main format.
//

//nolint:recvcheck
type Chronos struct {
	time.Time
}

// New creates new *Chronos from time.
func New(t time.Time) *Chronos {
	return &Chronos{t}
}

// FromTime creates new Chronos from time.
func FromTime(t time.Time) Chronos {
	return Chronos{t}
}

// FromString - parses time as string format and return Chronos,
// the function can parse the layouts specified in the formats array.
func FromString(timeStr string) (Chronos, error) {
	if timeStr == "" {
		return Chronos{}, errEmptyString
	}

	formats := []Format{
		// Formats that Chronos can parse.
		PacificDaylight,
		JavaScriptUTC,
		JavaScriptUTC7,
		JavaScriptLocal,
		ISO8601MilliTZ,
		ISO8601MilliUTC,
		UTCDateTime,
		DateTime,
		DateTimeNoSeconds,
		Standard,
		DateOnly,
		HumanDate,
		Exchange,
		DateTimeT,
		ISO8601Nano,
		ISO8601NanoZ,
		ISO8601NanoTZ,
	}

	for _, format := range formats {
		t, err := time.Parse(format.String(), timeStr)
		if err == nil {
			return Chronos{t}, nil
		}
	}

	return Chronos{}, fmt.Errorf("%w: %s", errUnableParse, timeStr)
}

// Parse parses time as string format and return Chronos.
func Parse(format Format, timeStr string) (Chronos, error) {
	if timeStr == "" {
		return Chronos{}, errEmptyString
	}

	t, err := time.Parse(format.String(), timeStr)
	if err != nil {
		return Chronos{}, fmt.Errorf("%w: %s", errUnableParse, timeStr)
	}

	return FromTime(t), nil
}

// MarshalJSON implements the [json.Marshaler] interface.
func (chr Chronos) MarshalJSON() ([]byte, error) {
	if chr.IsZero() {
		return []byte(`""`), nil
	}

	formatted := chr.StandardFormat()

	return json.Marshal(formatted)
}

// UnmarshalJSON implements the [json.Unmarshaler] interface.
func (chr *Chronos) UnmarshalJSON(b []byte) error {
	trim := strings.Trim(string(b), `"`)
	if trim == "" {
		chr.Time = time.Time{}

		return nil
	}

	t, err := FromString(trim)
	if err != nil {
		return err
	}

	chr.Time = t.Time

	return nil
}

// Scanner is an interface used by [Rows.Scan].
var _ sql.Scanner = (*Chronos)(nil)

// Scan implements the database/sql Scanner interface.
func (chr *Chronos) Scan(src any) error {
	if src == nil {
		return nil
	}

	switch srcType := src.(type) {
	case time.Time:
		chr.Time = srcType

		return nil
	case []byte:
		t, err := FromString(string(srcType))
		if err != nil {
			return err
		}

		chr.Time = t.Time

		return nil
	case string:
		t, err := FromString(srcType)
		if err != nil {
			return err
		}

		chr.Time = t.Time

		return nil
	default:
		return fmt.Errorf("%w: %T", errUnsupportedType, src)
	}
}

// Value implements the database/sql/driver Valuer interface.
func (chr Chronos) Value() (driver.Value, error) {
	if chr.IsZero() {
		//nolint:nilnil
		return nil, nil
	}

	return chr.Time, nil
}

// Format returns time in the specified format or empty string if time is zero.
func (chr Chronos) Format(format Format) string {
	// Return empty string if time is zero
	if chr.IsZero() {
		return ""
	}

	return chr.Time.Format(string(format))
}

// StandardFormat returns time in the RFC3339 format or empty string if time is zero.
func (chr Chronos) StandardFormat() string {
	return chr.Format(Standard)
}

// FormatLoc formats the time in the specified location and layout.
func (chr Chronos) FormatLoc(format Format, loc *time.Location) string {
	if chr.IsZero() {
		return ""
	}

	return chr.Time.In(loc).Format(string(format))
}

// FormatRuLocale returns time in specified format with russian locale.
func (chr Chronos) FormatRuLocale(format Format) string {
	if chr.IsZero() {
		return ""
	}

	return monday.Format(chr.Time, string(format), monday.LocaleRuRU)
}

// String returns time as string or empty string if time is zero.
func (chr Chronos) String() string {
	return chr.StandardFormat()
}

// Now returns current time as Chronos.
func Now() Chronos {
	return Chronos{time.Now()}
}

// Today returns current date at 00:00:00 time.
func Today() Chronos {
	now := time.Now()

	return Chronos{time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())}
}

// IsSameDateParse parse and compares two dates, return true if equal.
func (chr Chronos) IsSameDateParse(other string) (bool, error) {
	otherDate, err := FromString(other)
	if err != nil {
		return false, err
	}

	// Compare dates
	return chr.IsSameDate(otherDate), nil
}

// IsSameDate compares two dates, return true if equal.
func (chr Chronos) IsSameDate(other Chronos) bool {
	t1 := chr.Time.UTC()
	t2 := other.Time.UTC()

	return t1.Year() == t2.Year() &&
		t1.Month() == t2.Month() &&
		t1.Day() == t2.Day()
}

// IsSameTime compares two times, return true if equal.
func (chr Chronos) IsSameTime(other Chronos) bool {
	t1 := chr.Time.UTC()
	t2 := other.Time.UTC()

	return t1.Hour() == t2.Hour() &&
		t1.Minute() == t2.Minute()
}

// Equal reports whether chronos and other chronos represent the same time instant.
func (chr Chronos) Equal(other Chronos) bool {
	return chr.Time.Equal(other.Time)
}

// Before reports whether the time instant chronos is before others.
func (chr Chronos) Before(other Chronos) bool {
	return chr.Time.Before(other.Time)
}

// After reports whether the time instant chronos is after others.
func (chr Chronos) After(other Chronos) bool {
	return chr.Time.After(other.Time)
}

// StartOfDay returns the start of the day (00:00:00).
func (chr Chronos) StartOfDay() Chronos {
	if chr.IsZero() {
		return chr
	}

	return Chronos{time.Date(chr.Year(), chr.Month(), chr.Day(), 0, 0, 0, 0, chr.Location())}
}

// StartOfDayUTC возвращает начало дня в UTC (00:00:00).
func (chr Chronos) StartOfDayUTC() Chronos {
	if chr.IsZero() {
		return chr
	}

	return Chronos{time.Date(chr.Year(), chr.Month(), chr.Day(), 0, 0, 0, 0, time.UTC)}
}

// EndOfDay returns the end of the day (23:59:59.999999999).
func (chr Chronos) EndOfDay() Chronos {
	if chr.IsZero() {
		return chr
	}

	return Chronos{time.Date(chr.Year(), chr.Month(), chr.Day(), 23, 59, 59, 999999999, chr.Location())}
}

// EndOfDayUTC возвращает конец дня в UTC(23:59:59.999999999)..
func (chr Chronos) EndOfDayUTC() Chronos {
	if chr.IsZero() {
		return chr
	}

	return Chronos{time.Date(chr.Year(), chr.Month(), chr.Day(), 23, 59, 59, 999999999, time.UTC)}
}

// FirstDayOfMonth returns the first day of the month.
func (chr Chronos) FirstDayOfMonth() Chronos {
	if chr.IsZero() {
		return chr
	}

	return Chronos{time.Date(chr.Year(), chr.Month(), 1, 0, 0, 0, 0, chr.Location())}
}

// LastDayOfMonth returns the last day of the month.
func (chr Chronos) LastDayOfMonth() Chronos {
	if chr.IsZero() {
		return chr
	}

	return Chronos{time.Date(chr.Year(), chr.Month()+1, 0, 0, 0, 0, 0, chr.Location())}
}

// WeekNumber returns the ISO 8601 week number within the year.
func (chr Chronos) WeekNumber() int {
	_, week := chr.ISOWeek()

	return week
}

// WeekOfMonth returns the week number within the month. The first week of the month is 1.
func (chr Chronos) WeekOfMonth() int {
	firstOfMonth := chr.FirstDayOfMonth()
	day := chr.Day()

	// Offset
	weekdayOffset := (int(firstOfMonth.Weekday()) + Sunday) % DaysInWeek
	week := (weekdayOffset + day + Sunday) / DaysInWeek

	return week
}

// DayName returns weekday as string.
func (chr Chronos) DayName() string {
	return chr.Weekday().String()
}

// DayNumber returns the day of the month (1–31).
func (chr Chronos) DayNumber() int {
	return chr.Day()
}

// WeekDayISONumber returns the ISO weekday number. Monday = 1, Sunday = 7.
func (chr Chronos) WeekDayISONumber() int {
	wd := int(chr.Weekday())
	if wd == 0 {
		return DaysInWeek // Sunday
	}

	return wd
}

// WeekDayNumber returns the weekday number. Monday = 0, Sunday = 6.
func (chr Chronos) WeekDayNumber() int {
	weekday := chr.Weekday()

	if weekday == time.Sunday {
		return Sunday
	}

	return int(weekday - 1)
}

// Add adds duration to time.
func (chr Chronos) Add(dur time.Duration) Chronos {
	if chr.IsZero() {
		return chr
	}

	return FromTime(chr.Time.Add(dur))
}

// AddHoursFloat adds float hours to time.
func (chr Chronos) AddHoursFloat(hours float64) Chronos {
	if chr.IsZero() {
		return chr
	}

	return Chronos{Time: chr.Time.Add(time.Duration(hours * float64(time.Hour)))}
}

// AddHours adds hours to time.
func (chr Chronos) AddHours(hours int) Chronos {
	if chr.IsZero() {
		return chr
	}

	return Chronos{Time: chr.Time.Add(time.Duration(hours) * time.Hour)}
}

// AddMinutes adds minutes to time.
func (chr Chronos) AddMinutes(minutes int) Chronos {
	if chr.IsZero() {
		return chr
	}

	return Chronos{Time: chr.Time.Add(time.Duration(minutes) * time.Minute)}
}

// AddHoursInt32 adds hours to time.
func (chr Chronos) AddHoursInt32(hours int32) Chronos {
	if chr.IsZero() {
		return chr
	}

	return Chronos{Time: chr.Time.Add(time.Duration(hours) * time.Hour)}
}

// AddDays adds specified number of days.
func (chr Chronos) AddDays(days int) Chronos {
	if chr.IsZero() {
		return chr
	}

	return Chronos{chr.Time.AddDate(0, 0, days)}
}

// AddMonths adds specified number of months.
func (chr Chronos) AddMonths(months int) Chronos {
	if chr.IsZero() {
		return chr
	}

	return Chronos{chr.Time.AddDate(0, months, 0)}
}

// AddYears adds specified number of years.
func (chr Chronos) AddYears(years int) Chronos {
	if chr.IsZero() {
		return chr
	}

	return Chronos{chr.Time.AddDate(years, 0, 0)}
}

// IsToday checks if the date is today.
func (chr Chronos) IsToday() bool {
	if chr.IsZero() {
		return false
	}

	now := time.Now()

	return chr.Year() == now.Year() && chr.Month() == now.Month() && chr.Day() == now.Day()
}

// IsPast checks if the date is in the past.
func (chr Chronos) IsPast() bool {
	if chr.IsZero() {
		return false
	}

	return chr.Time.Before(time.Now())
}

// IsFuture checks if the date is in the future.
func (chr Chronos) IsFuture() bool {
	if chr.IsZero() {
		return false
	}

	return chr.Time.After(time.Now())
}

// DaysBetween returns the number of days between two dates (absolute value).
func (chr Chronos) DaysBetween(other Chronos) int {
	if chr.IsZero() || other.IsZero() {
		return 0
	}

	// Normalize both dates
	cStart := chr.StartOfDay()
	otherStart := other.StartOfDay()
	duration := otherStart.Time.Sub(cStart.Time)

	return int(math.Abs(duration.Hours()) / OneDayHours)
}

// HoursBetween returns the number of hours between two dates and throws away minutes.
func (chr Chronos) HoursBetween(other Chronos) int {
	if chr.IsZero() || other.IsZero() {
		return 0
	}

	duration := other.Time.Sub(chr.Time)

	return int(math.Abs(duration.Hours()))
}

// Age calculates years from the day of birth.
func (chr Chronos) Age() int {
	if chr.IsZero() {
		return 0
	}

	now := time.Now()
	birth := chr.Time

	if birth.After(now) {
		return 0
	}

	age := now.Year() - birth.Year()
	birthdayThisYear := time.Date(now.Year(), chr.Month(), chr.Day(), 0, 0, 0, 0, now.Location())

	if now.Before(birthdayThisYear) {
		age--
	}

	return age
}

// StartOfWeek returns the first day of the week (Monday).
func (chr Chronos) StartOfWeek() Chronos {
	if chr.IsZero() {
		return chr
	}

	weekday := chr.Weekday()
	// Convert (Monday = 0, Sunday = 6)
	offset := int(weekday) - 1
	if weekday == time.Sunday {
		offset = Sunday
	}

	return chr.AddDays(-offset).StartOfDay()
}

// EndOfWeek returns the last day of the week (Sunday).
func (chr Chronos) EndOfWeek() Chronos {
	return chr.StartOfWeek().AddDays(Sunday).EndOfDay()
}

// StartOfMonth returns the first day of the month.
func (chr Chronos) StartOfMonth() Chronos {
	return Chronos{time.Date(chr.Year(), chr.Month(), 1, 0, 0, 0, 0, chr.Location())}
}

// EndOfMonth возвращает последний день месяца.
func (chr Chronos) EndOfMonth() Chronos {
	return chr.AddMonths(1).Add(-time.Nanosecond)
}

// IsLeapYear checks if the year is a leap year.
func (chr Chronos) IsLeapYear() bool {
	if chr.IsZero() {
		return false
	}

	year := chr.Year()

	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// WeekdayRu return week day in ru.
func (chr Chronos) WeekdayRu() string {
	switch chr.Weekday() {
	case time.Monday:
		return "пн"
	case time.Tuesday:
		return "вт"
	case time.Wednesday:
		return "ср"
	case time.Thursday:
		return "чт"
	case time.Friday:
		return "пт"
	case time.Saturday:
		return "сб"
	case time.Sunday:
		return "вс"
	default:
		return ""
	}
}

// UTC return chronos in UTC.
func (chr Chronos) UTC() Chronos {
	return Chronos{chr.Time.UTC()}
}

// Truncate returns the result of rounding t down to a multiple of d (since the zero time).
func (chr Chronos) Truncate(d time.Duration) Chronos {
	return Chronos{chr.Time.Truncate(d)}
}

// Sub returns the duration t-u. If the result exceeds the maximum (or minimum).
func (chr Chronos) Sub(other Chronos) time.Duration {
	return chr.Time.Sub(other.Time)
}

// AddDate returns the time corresponding to adding the
// given number of years, months, and days to t.
func (chr Chronos) AddDate(years int, months int, days int) Chronos {
	return Chronos{chr.Time.AddDate(years, months, days)}
}

// Local returns t with the location set to local time.
func (chr Chronos) Local() Chronos {
	//nolint:gosmopolitan
	return Chronos{chr.Time.Local()}
}

// Location returns the time zone information associated with t.
func (chr Chronos) Location() *time.Location {
	return chr.Time.Location()
}

// Zone computes the time zone in effect at time t, returning the abbreviated
// name of the zone (such as "CET") and its offset in seconds east of UTC.
func (chr Chronos) Zone() (string, int) {
	return chr.Time.Zone()
}

// ZoneBounds returns the bounds of the time zone in effect at time t.
func (chr Chronos) ZoneBounds() (Chronos, Chronos) {
	start, end := chr.Time.ZoneBounds()

	return Chronos{start}, Chronos{end}
}

// Unix returns t as a Unix time, the number of seconds elapsed
// since January 1, 1970 UTC. The result does not depend on the
// location associated with t.
func (chr Chronos) Unix() int64 {
	return chr.Time.Unix()
}

// UnixMilli returns t as a Unix time, the number of milliseconds elapsed since
// January 1, 1970 UTC. The result is undefined if the Unix time in
// milliseconds cannot be represented by an int64 (a date more than 292 million
// years before or after 1970). The result does not depend on the
// location associated with t.
func (chr Chronos) UnixMilli() int64 {
	return chr.Time.UnixMilli()
}

// UnixMicro returns t as a Unix time, the number of microseconds elapsed since
// January 1, 1970 UTC. The result is undefined if the Unix time in
// microseconds cannot be represented by an int64 (a date before year -290307 or
// after year 294246). The result does not depend on the location associated
// with t.
func (chr Chronos) UnixMicro() int64 {
	return chr.Time.UnixMicro()
}

// UnixNano returns t as a Unix time, the number of nanoseconds elapsed
// since January 1, 1970 UTC. The result is undefined if the Unix time
// in nanoseconds cannot be represented by an int64 (a date before the year
// 1678 or after 2262). Note that this means the result of calling UnixNano
// on the zero Time is undefined. The result does not depend on the
// location associated with t.
func (chr Chronos) UnixNano() int64 {
	return chr.Time.UnixNano()
}

// Round returns the result of rounding t to the nearest multiple of d (since the zero time).
func (chr Chronos) Round(d time.Duration) Chronos {
	return Chronos{chr.Time.Round(d)}
}

// SetLocation - return Chronos in with zone created from offset.
func (chr Chronos) SetLocation(offset int) Chronos {
	loc := CreateZoneFromOffset(offset)

	return Chronos{chr.In(loc)}
}

// InLoc converts time to the specified location.
func (t Chronos) InLoc(loc *time.Location) Chronos {
	return FromTime(t.In(loc))
}

// InLocOffset converts time to a timezone with the given offset.
func (t Chronos) InLocOffset(offset int) Chronos {
	zone := CreateZoneFromOffset(offset)

	return t.InLoc(zone)
}

// InitFiberDecoder - позволяет декодировать query параметр в запросе.
func InitFiberDecoder() {
	fiber.SetParserDecoder(fiber.ParserConfig{
		IgnoreUnknownKeys: true,
		ZeroEmpty:         true,
		ParserType: []fiber.ParserType{
			{
				Customtype: Chronos{},
				Converter: func(value string) reflect.Value {
					if value == "" {
						return reflect.ValueOf(Chronos{})
					}

					t, err := FromString(value)
					if err != nil {
						return reflect.ValueOf(Chronos{})
					}

					return reflect.ValueOf(t)
				},
			},
		},
	})
}
