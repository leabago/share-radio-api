package chronos

import "time"

const (
	// Formats that Chronos can parse.

	Standard          Format = time.RFC3339 // used as the main format
	DateOnly          Format = "2006-01-02"
	DateTime          Format = "2006-01-02 15:04:05"
	DateTimeNoSeconds Format = "2006-01-02 15:04"
	UTCDateTime       Format = "2006-01-02T15:04:05Z"
	ISO8601MilliTZ    Format = "2006-01-02T15:04:05.000Z07:00"
	ISO8601MilliUTC   Format = "2006-01-02T15:04:05.000Z"
	JavaScriptUTC     Format = "Mon Jan 02 2006 15:04:05 GMT+0000 (Coordinated Universal Time)"
	JavaScriptLocal   Format = "Mon Jan 02 2006 15:04:05 GMT-0700 (MST)"
	JavaScriptUTC7    Format = "Mon Jan 02 2006 15:04:05 GMT-0700 (Coordinated Universal Time)"
	HumanDate         Format = "02.01.2006"
	Exchange          Format = "2006-01-02T15:04:05"
	DateTimeT         Format = "2006-01-02T15:04"

	ISO8601Nano   Format = "2006-01-02T15:04:05.999999999"       // Same as ISO8601Milli but with 9 digits
	ISO8601NanoZ  Format = "2006-01-02T15:04:05.999999999Z"      // With Z timezone
	ISO8601NanoTZ Format = "2006-01-02T15:04:05.999999999Z07:00" // With timezone offset

	// Additional formats used in the project.

	ColonSeparated Format = "2006-01-02:15:04"

	// Deprecated: Use Standard instead.
	Vitasystem       Format = "2006-01-02T15:04:05Z07:00"
	VitasystemZ      Format = "Z07:00"
	PacificDaylight  Format = "2006-01-02 15:04:05 -0700"
	HumanDateShort   Format = "2.1.2006"
	NotifShiftUpdate Format = "2 January 2006 15:04"
	Month            Format = "January"
	VSDateTimeFormat Format = "2006-01-02T15:04:05 Z"
	PharmWtime       Format = "1504"
	HoursMinutes     Format = "15:04"
	// Deprecated: Use HoursMinutes instead.
	PharmWorkTime      Format = "15:04"
	VsNeedPartner      Format = "02.01.2006 15:04:05"
	HumanDateShortYear Format = "02.01.06"
)

// Format define the pattern used for parsing and formatting dates and times.
type Format string

func (f Format) String() string {
	return string(f)
}
