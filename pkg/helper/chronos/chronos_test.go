package chronos

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

const TestDateTime = "2025-10-17 15:04:05"

type TestDate struct {
	Date Chronos `json:"date"`
}

func TestFromString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"RFC3339", "2025-10-19T15:04:05Z", false},
		{"DateOnly", "2025-10-19", false},
		{"DateTime", "2025-10-19 15:04:05", false},
		{"JSFormat", "Mon Jan 02 2006 15:04:05 GMT+0000 (Invalid Coordinated Universal Time)", true},
		{"Invalid", "not-a-date", true},
		{"Empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromString(tt.input)
			if tt.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if err == nil && got.IsZero() {
				t.Errorf("expected non-zero time")
			}
		})
	}
}

func TestMarshalUnmarshalJSON(t *testing.T) {
	c := Chronos{time.Date(2025, 10, 19, 12, 1, 2, 0, time.UTC)}
	data, err := json.Marshal(c)
	require.NoError(t, err)

	var d Chronos

	err = json.Unmarshal(data, &d)
	require.NoError(t, err)

	if !d.Equal(c) {
		t.Error("expected equal times")
	}
}

func TestCompareDate(t *testing.T) {
	time1 := "2025-01-02"
	time2 := "2025-01-02T15:04:05Z"

	time1Ch, err := FromString(time1)
	require.NoError(t, err)

	equal, err := time1Ch.IsSameDateParse(time2)
	require.NoError(t, err)
	assert.True(t, equal)

	b := Chronos{time.Date(2025, 01, 02, 23, 59, 0, 0, time.UTC)}
	if !time1Ch.IsSameDate(b) {
		t.Error("IsSameDate() should be true")
	}
}

func TestFormat_ZeroTime(t *testing.T) {
	ch := FromTime(time.Time{})
	assert.Empty(t, ch.Format(DateOnly))
	assert.Empty(t, ch.String())
}

func TestMarshalJSON(t *testing.T) {
	ch, err := FromString(TestDateTime)
	require.NoError(t, err)
	chronosMarshal, err := ch.MarshalJSON()
	require.NoError(t, err)
	assert.Equal(t, `"2025-10-17T15:04:05Z"`, string(chronosMarshal))

	c := Chronos{}
	data, err := json.Marshal(c)
	require.NoError(t, err)

	if string(data) != `""` {
		t.Errorf("expected empty string, got %s", data)
	}
}

func TestValueAndScan(t *testing.T) {
	c := Chronos{time.Date(2025, 10, 19, 12, 0, 0, 0, time.UTC)}

	val, err := c.Value()
	if err != nil {
		t.Fatalf("Value() failed: %v", err)
	}

	var scanned Chronos

	timeVal, ok := val.(time.Time)
	if !ok {
		t.Fatalf("expected time.Time, got %T", val)
	}

	err = scanned.Scan(timeVal)
	if err != nil {
		t.Fatalf("Scan() failed: %v", err)
	}

	if !scanned.Equal(c) {
		t.Error("scanned value differs")
	}

	// interface conformance
	var _ driver.Valuer = (*Chronos)(nil)
}

func TestCustomMarshalJSON(t *testing.T) {
	chron1, err := FromString(TestDateTime)
	require.NoError(t, err)

	// Custom type with other MarshalJSON function
	cc := ModCustomTime{chron1}
	customTimeMarshal, err := cc.MarshalJSON()
	require.NoError(t, err)
	assert.Equal(t, `"2025-10-17T15:04:05.000Z"`, string(customTimeMarshal))
}

func TestUnmarshalJSON(t *testing.T) {
	jsonData := `"2025-12-25T15:04:05.123Z"`

	var result ModCustomTime

	err := json.Unmarshal([]byte(jsonData), &result)
	require.NoError(t, err)
	assert.Equal(t, "2025-12-25T15:04:05Z", result.UTC().String())
}

func TestStandardFormat(t *testing.T) {
	// RFC3339
	ch1, err := FromString(TestDateTime)
	require.NoError(t, err)
	assert.Equal(t, "2025-10-17T15:04:05Z", ch1.StandardFormat())

	// Zero time
	ch2 := FromTime(time.Time{})

	require.NoError(t, err)
	assert.Empty(t, ch2.StandardFormat())

	// func String() also RFC3339
	assert.Equal(t, ch1.StandardFormat(), ch1.String())
}

func TestParseAndFormat(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		format    Format
		want      string
		expectErr bool
	}{
		{
			name:   "valid UTC JS date",
			input:  "Mon Jan 02 2025 15:04:05 GMT+0000 (Coordinated Universal Time)",
			format: ISO8601MilliUTC,
			want:   "2025-01-02T15:04:05.000Z",
		},
		{
			name:      "invalid valid UTC JS date",
			input:     "Mon Jan 02 2025 15:04:05",
			format:    ISO8601MilliUTC,
			want:      "",
			expectErr: true,
		},
		{
			name:   "RFC3339 to ISO",
			input:  "2025-10-01T10:00:00Z",
			format: ISO8601MilliUTC,
			want:   "2025-10-01T10:00:00.000Z",
		},
		{
			name:   "DateOnly only",
			input:  "2025-10-01",
			format: DateOnly,
			want:   "2025-10-01",
		},
		{
			name:      "Invalid input",
			input:     "invalid-date",
			format:    ISO8601MilliUTC,
			expectErr: true,
		},
		{
			name:   "Empty input",
			input:  "",
			format: ISO8601MilliUTC,
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAndFormat(tt.input, tt.format)
			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Errorf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestFormatLoc(t *testing.T) {
	// FormatLoc
	date1 := "2025-10-17 15:04:05"
	loc1 := time.FixedZone(fmt.Sprintf("%+03d:00", 1), int(1*3600))
	ch, err := FromString(date1)
	require.NoError(t, err)

	form := ch.FormatLoc(time.RFC3339, loc1)
	assert.Equal(t, "2025-10-17T16:04:05+01:00", form)

	// ParseAndFormatLoc
	loc2, err := time.LoadLocation("Asia/Novosibirsk")
	require.NoError(t, err)

	dateFormated := ParseAndFormatLoc(date1, ISO8601MilliTZ, loc2)
	assert.Equal(t, "2025-10-17T22:04:05.000+07:00", dateFormated)
}

func TestNowAndToday(t *testing.T) {
	if Now().IsZero() {
		t.Error("Now() returned zero value")
	}

	if Today().Hour() != 0 {
		t.Error("Today should start at 0")
	}
}

func TestAddDaysMonthsYears(t *testing.T) {
	base := Chronos{time.Date(2025, 10, 19, 0, 0, 0, 0, time.UTC)}
	if base.AddDays(1).Day() != 20 {
		t.Error("AddDays failed")
	}

	if base.AddMonths(1).Month() != 11 {
		t.Error("AddMonths failed")
	}

	if base.AddYears(1).Year() != 2026 {
		t.Error("AddYears failed")
	}
}

func TestAddHours(t *testing.T) {
	base, err := FromString(TestDateTime)
	require.NoError(t, err)
	// Int
	assert.Equal(t, "2025-10-17T18:04:05Z", base.AddHours(3).StandardFormat())
	assert.Equal(t, "2025-10-18T14:04:05Z", base.AddHours(23).StandardFormat())
	assert.Equal(t, "2025-10-17T13:04:05Z", base.AddHours(-2).StandardFormat())

	// Float
	assert.Equal(t, "2025-10-17T18:04:05Z", base.AddHoursFloat(3).StandardFormat())
	assert.Equal(t, "2025-10-18T14:04:05Z", base.AddHoursFloat(23).StandardFormat())
	assert.Equal(t, "2025-10-17T13:04:05Z", base.AddHoursFloat(-2).StandardFormat())
	assert.Equal(t, "2025-10-17T16:34:05Z", base.AddHoursFloat(1.5).StandardFormat())
}

func TestFormatLoc_ZeroValue(t *testing.T) {
	loc := time.FixedZone(fmt.Sprintf("%+03d:00", 1), int(1*3600))
	ch := FromTime(time.Time{})
	form := ch.FormatLoc(time.RFC3339, loc)
	assert.Empty(t, form)
}

func TestHoursDuration(t *testing.T) {
	start, err := FromString("2025-10-17 15:04:05")
	require.NoError(t, err)

	end, err := FromString("2025-10-17 19:49:05")
	require.NoError(t, err)

	assert.Equal(t, 4, HoursBetween(start, end))
	assert.Equal(t, 4, HoursBetween(end, start))

	// ModExchangeTime
	excStart := ModExchangeTime{start}
	excEnd := ModExchangeTime{end}

	hours := HoursBetween(excStart.Chronos, excEnd.Chronos)
	assert.Equal(t, 4, hours)
}

func TestToday(t *testing.T) {
	ch1 := Today()
	ch2 := Today()

	ch1 = ch1.StartOfDay()
	ch2 = ch2.EndOfDay()

	duration1 := ch2.Sub(ch1)
	assert.Equal(t, "23h59m59.999999999s", duration1.String())
}

func TestDaysBetween(t *testing.T) {
	ch1, err := FromString("2025-10-17 11:04:05")
	require.NoError(t, err)

	ch2, err := FromString("2025-10-20 20:04:05")
	require.NoError(t, err)

	assert.Equal(t, 3, ch1.DaysBetween(ch2))
}

func TestAge(t *testing.T) {
	ch1, err := FromString("2001-01-01")
	require.NoError(t, err)

	assert.GreaterOrEqual(t, ch1.Age(), 24)
}

func TestSatrEndWeek(t *testing.T) {
	ch1, err := FromString("2025-10-19 15:04:05")
	require.NoError(t, err)

	start := ch1.StartOfWeek()
	end := ch1.EndOfWeek()

	days := DaysBetween(start, end)
	assert.Equal(t, 6, days)
}

func TestIsLeapYear(t *testing.T) {
	ch1, err := FromString("2024-10-19 15:04:05")
	require.NoError(t, err)

	ch2, err := FromString("2025-10-19 15:04:05")
	require.NoError(t, err)

	ch3, err := FromString("2026-10-19 15:04:05")
	require.NoError(t, err)

	ch4, err := FromString("2027-10-19 15:04:05")
	require.NoError(t, err)

	assert.True(t, ch1.IsLeapYear())
	assert.False(t, ch2.IsLeapYear())
	assert.False(t, ch3.IsLeapYear())
	assert.False(t, ch4.IsLeapYear())
}

func TestMarshal(t *testing.T) {
	chron1, err := FromString(TestDateTime)
	require.NoError(t, err)

	date1 := TestDate{Date: chron1}
	dateJson, err := json.Marshal(date1)
	require.NoError(t, err)

	assert.JSONEq(t, `{"date":"2025-10-17T15:04:05Z"}`, string(dateJson))
}

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Standard", `{ "date": "2025-10-15T22:03:19+04:00"}`, false},
		{"DateOnly", `{ "date": "2025-10-15"}`, false},
		{"DateTime", `{ "date": "2025-10-15 22:01:25"}`, false},
		{"DateTimeNoSeconds", `{ "date": "2025-10-15 22:01"}`, false},
		{"UTCDateTime", `{ "date": "2025-10-15T22:01:25Z"}`, false},
		{"ISO8601MilliTZ", `{ "date": "2025-10-15T22:01:25.424+04:00"}`, false},
		{"ISO8601MilliUTC", `{ "date": "2025-10-15T22:01:25.424Z"}`, false},
		{"JavaScriptUTC", `{ "date": "Sun Oct 15 2025 22:01:25 GMT+0000 (Coordinated Universal Time)"}`, false},
		{"JavaScriptLocal", `{ "date": "Sun Oct 15 2025 22:01:25 GMT+0400 (+04)"}`, false},
		{"Invalid format", `{ "date": "Invalid"}`, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var date TestDate

			err := json.Unmarshal([]byte(tt.input), &date)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestComparisons(t *testing.T) {
	a := Chronos{time.Date(2025, 10, 19, 0, 0, 0, 0, time.UTC)}
	a2 := Chronos{time.Date(2025, 10, 19, 0, 0, 0, 0, time.UTC)}
	b := Chronos{time.Date(2025, 10, 20, 0, 0, 0, 0, time.UTC)}

	if !a.Before(b) {
		t.Error("Before() failed")
	}

	if !b.After(a) {
		t.Error("After() failed")
	}

	if !a.Equal(a2) {
		t.Error("Equal() failed")
	}
}

func TestIsTodayPastFuture(t *testing.T) {
	now := Now()
	if !now.IsToday() {
		t.Error("IsToday failed")
	}

	past := Chronos{time.Now().Add(-24 * time.Hour)}
	if !past.IsPast() {
		t.Error("IsPast failed")
	}

	future := Chronos{time.Now().Add(24 * time.Hour)}
	if !future.IsFuture() {
		t.Error("IsFuture failed")
	}
}

func TestDayAndMonthBoundaries(t *testing.T) {
	c := Chronos{time.Date(2025, 10, 19, 15, 0, 0, 0, time.UTC)}
	if c.StartOfDay().Hour() != 0 {
		t.Error("StartOfDay failed")
	}

	if c.EndOfDay().Hour() != 23 {
		t.Error("EndOfDay failed")
	}

	if c.FirstDayOfMonth().Day() != 1 {
		t.Error("FirstDayOfMonth failed")
	}

	if c.LastDayOfMonth().Day() != 31 {
		t.Error("LastDayOfMonth failed")
	}

	c2, err := FromString("2025-02-15")
	require.NoError(t, err)

	if c2.FirstDayOfMonth().Day() != 1 {
		t.Error("FirstDayOfMonth failed")
	}

	if c2.LastDayOfMonth().Day() != 28 {
		t.Error("LastDayOfMonth failed")
	}
}

func TestParseFromQueryFiberError(t *testing.T) {
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Request().URI().SetQueryString("date=2025-10-20&dateInvalid=Invalid&dateEmpty=")

	// invalid
	val := ctx.Query("dateInvalid")
	if val == "" {
		t.Fatalf("dateInvalid is empty")
	}

	ch, err := ParseTimeWithFormats("dateInvalid", val, []Format{Standard, DateOnly})
	require.Error(t, err)
	assert.True(t, ch.IsZero())

	err = fiber.NewError(fiber.StatusBadRequest, err.Error())
	assert.Equal(t, "parameter has invalid time format: parameter: dateInvalid, expected: [\"2006-01-02T15:04:05Z07:00\" \"2006-01-02\"], got: \"Invalid\"", err.Error())
}
func TestParseFromQuery(t *testing.T) {
	tests := []struct {
		name        string
		queryParam  string
		queryValue  string
		formats     []Format
		wantError   bool
		wantZero    bool
		description string
	}{
		{
			name:        "valid_date",
			queryParam:  "date",
			queryValue:  "2025-10-20",
			formats:     []Format{Standard, DateOnly},
			wantError:   false,
			wantZero:    false,
			description: "should parse valid date successfully",
		},
		{
			name:        "invalid_date",
			queryParam:  "dateInvalid",
			queryValue:  "Invalid",
			formats:     []Format{Standard, DateOnly},
			wantError:   true,
			wantZero:    true,
			description: "should return error for invalid date format",
		},
		{
			name:        "empty_date",
			queryParam:  "dateEmpty",
			queryValue:  "",
			formats:     []Format{Standard, DateOnly},
			wantError:   true,
			wantZero:    true,
			description: "should return error for empty date value",
		},
		{
			name:        "valid_date_with_different_format",
			queryParam:  "dateRFC3339",
			queryValue:  "2025-10-20T15:30:00Z",
			formats:     []Format{Standard, DateOnly},
			wantError:   false,
			wantZero:    false,
			description: "should parse RFC3339 format successfully",
		},
		{
			name:        "fail_with_different_format",
			queryParam:  "JavaScriptUTC",
			queryValue:  "Mon Jan 02 2006 15:04:05 GMT+0000 (Coordinated Universal Time)",
			formats:     []Format{Standard, DateOnly},
			wantError:   true,
			wantZero:    true,
			description: "should return error, different format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			ctx := app.AcquireCtx(&fasthttp.RequestCtx{})

			// Set up the query string
			ctx.Request().URI().SetQueryString(tt.queryParam + "=" + tt.queryValue)

			// Get the query value
			val := ctx.Query(tt.queryParam)

			// For empty values, verify they're actually empty
			if tt.queryValue == "" {
				assert.Empty(t, val, "expected empty query value")
			} else {
				assert.NotEmpty(t, val, "query value should not be empty")
			}

			// Parse the time
			result, err := ParseTimeWithFormats(tt.queryParam, val, tt.formats)

			// Check error expectations
			if tt.wantError {
				require.Error(t, err, "expected error for: %s", tt.description)
			} else {
				require.NoError(t, err, "expected no error for: %s", tt.description)
			}

			// Check zero time expectations
			if tt.wantZero {
				assert.True(t, result.IsZero(), "expected zero time for: %s", tt.description)
			} else {
				assert.False(t, result.IsZero(), "expected non-zero time for: %s", tt.description)
			}

			app.ReleaseCtx(ctx)
		})
	}
}

func TestParseTimeWithFormats(t *testing.T) {
	c, err := ParseTimeWithFormats("date", "2025-10-20", []Format{DateOnly})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}

	if c.IsZero() {
		t.Error("expected parsed time")
	}

	_, err = ParseTimeWithFormats("date", "bad", []Format{DateOnly})
	if err == nil {
		t.Error("expected error for invalid format")
	}
}

func TestSince(t *testing.T) {
	ch1 := Now()
	s1 := Since(ch1)
	assert.NotEqual(t, 0, s1.Nanoseconds())

	zero := FromTime(time.Time{})
	s2 := Since(zero)
	assert.Equal(t, int64(0), s2.Nanoseconds())
}

func TestWeekFunctions(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		wantISOYear     int
		wantISOWeek     int
		wantWeekOfMonth int
	}{
		{
			name:            "Start of year",
			input:           "2025-01-01", // Wednesday
			wantISOYear:     2025,
			wantISOWeek:     1,
			wantWeekOfMonth: 1,
		},
		{
			name:            "Start of February",
			input:           "2025-02-01", // Saturday
			wantISOYear:     2025,
			wantISOWeek:     5,
			wantWeekOfMonth: 1,
		},
		{
			name:            "Middle of year",
			input:           "2025-06-15", // Sunday
			wantISOYear:     2025,
			wantISOWeek:     24,
			wantWeekOfMonth: 3,
		},
		{
			name:            "End of October",
			input:           "2025-10-31", // Friday
			wantISOYear:     2025,
			wantISOWeek:     44,
			wantWeekOfMonth: 5,
		},
		{
			name:            "Leap year edge",
			input:           "2024-02-29", // Thursday (leap year)
			wantISOYear:     2024,
			wantISOWeek:     9,
			wantWeekOfMonth: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch, err := FromString(tt.input)
			require.NoError(t, err)

			isoYear, isoWeek := ch.ISOWeek()
			assert.Equal(t, tt.wantISOYear, isoYear)
			assert.Equal(t, tt.wantISOWeek, isoWeek)

			gotWeekOfMonth := ch.WeekOfMonth()
			assert.Equal(t, tt.wantWeekOfMonth, gotWeekOfMonth)
		})
	}
}

func TestChronos_WeekdayFunctions(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		wantDayNumber  int
		wantDayName    string
		wantDayISOWeek int
		wantDayWeek    int
	}{
		{
			name:           "Monday",
			input:          "2025-10-20", // Monday
			wantDayNumber:  20,
			wantDayName:    "Monday",
			wantDayISOWeek: 1,
			wantDayWeek:    0,
		},
		{
			name:           "Tuesday",
			input:          "2025-10-21", // Tuesday
			wantDayNumber:  21,
			wantDayName:    "Tuesday",
			wantDayISOWeek: 2,
			wantDayWeek:    1,
		},
		{
			name:           "Wednesday",
			input:          "2025-10-22", // Wednesday
			wantDayNumber:  22,
			wantDayName:    "Wednesday",
			wantDayISOWeek: 3,
			wantDayWeek:    2,
		},
		{
			name:           "Thursday",
			input:          "2025-10-23", // Thursday
			wantDayNumber:  23,
			wantDayName:    "Thursday",
			wantDayISOWeek: 4,
			wantDayWeek:    3,
		},
		{
			name:           "Friday",
			input:          "2025-10-24", // Friday
			wantDayNumber:  24,
			wantDayName:    "Friday",
			wantDayISOWeek: 5,
			wantDayWeek:    4,
		},
		{
			name:           "Saturday",
			input:          "2025-10-25", // Saturday
			wantDayNumber:  25,
			wantDayName:    "Saturday",
			wantDayISOWeek: 6,
			wantDayWeek:    5,
		},
		{
			name:           "Sunday",
			input:          "2025-10-26", // Sunday
			wantDayNumber:  26,
			wantDayName:    "Sunday",
			wantDayISOWeek: 7,
			wantDayWeek:    6,
		},
		{
			name:           "Friday2",
			input:          "2025-10-31", // Friday
			wantDayNumber:  31,
			wantDayName:    "Friday",
			wantDayISOWeek: 5,
			wantDayWeek:    4,
		},
		{
			name:           "Saturday2",
			input:          "2025-11-01", // Saturday
			wantDayNumber:  1,
			wantDayName:    "Saturday",
			wantDayISOWeek: 6,
			wantDayWeek:    5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch, err := FromString(tt.input)
			require.NoError(t, err)

			assert.Equal(t, tt.wantDayNumber, ch.DayNumber(), "DayNumber mismatch")
			assert.Equal(t, tt.wantDayName, ch.DayName(), "DayName mismatch")
			assert.Equal(t, tt.wantDayISOWeek, ch.WeekDayISONumber(), "WeekDayISONumber mismatch")
			assert.Equal(t, tt.wantDayWeek, ch.WeekDayNumber(), "WeekDayNumber mismatch")
		})
	}
}

func TestChronos_StartOfWeek(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected time.Time
	}{
		{
			name:     "Monday",
			input:    time.Date(2025, 10, 20, 15, 30, 45, 0, time.UTC), // Monday
			expected: time.Date(2025, 10, 20, 0, 0, 0, 0, time.UTC),    // Same day
		},
		{
			name:     "Tuesday",
			input:    time.Date(2025, 10, 21, 10, 0, 0, 0, time.UTC), // Tuesday
			expected: time.Date(2025, 10, 20, 0, 0, 0, 0, time.UTC),  // Monday
		},
		{
			name:     "Wednesday",
			input:    time.Date(2025, 10, 22, 20, 45, 30, 0, time.UTC), // Wednesday
			expected: time.Date(2025, 10, 20, 0, 0, 0, 0, time.UTC),    // Monday
		},
		{
			name:     "Thursday",
			input:    time.Date(2025, 10, 23, 8, 15, 0, 0, time.UTC), // Thursday
			expected: time.Date(2025, 10, 20, 0, 0, 0, 0, time.UTC),  // Monday
		},
		{
			name:     "Friday",
			input:    time.Date(2025, 10, 24, 18, 30, 0, 0, time.UTC), // Friday
			expected: time.Date(2025, 10, 20, 0, 0, 0, 0, time.UTC),   // Monday
		},
		{
			name:     "Saturday",
			input:    time.Date(2025, 10, 25, 12, 0, 0, 0, time.UTC), // Saturday
			expected: time.Date(2025, 10, 20, 0, 0, 0, 0, time.UTC),  // Monday
		},
		{
			name:     "Sunday",
			input:    time.Date(2025, 10, 26, 23, 59, 59, 0, time.UTC), // Sunday
			expected: time.Date(2025, 10, 20, 0, 0, 0, 0, time.UTC),    // Monday
		},
		{
			name:     "Week spanning month boundary",
			input:    time.Date(2025, 10, 2, 12, 0, 0, 0, time.UTC), // Thursday
			expected: time.Date(2025, 9, 29, 0, 0, 0, 0, time.UTC),  // Monday
		},
		{
			name:     "Week spanning year boundary",
			input:    time.Date(2025, 11, 3, 12, 0, 0, 0, time.UTC), // Monday
			expected: time.Date(2025, 11, 3, 0, 0, 0, 0, time.UTC),  // Same day at start
		},
		{
			name:     "Leap year",
			input:    time.Date(2024, 2, 29, 12, 0, 0, 0, time.UTC), // Thursday (leap day)
			expected: time.Date(2024, 2, 26, 0, 0, 0, 0, time.UTC),  // Monday
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := FromTime(tt.input)
			result := c.StartOfWeek()

			if !result.Time.Equal(tt.expected) {
				t.Errorf("StartOfWeek() for %v = %v, expected %v", tt.input, result.Time, tt.expected)
			}

			// Verify it's indeed Monday
			if result.Weekday() != time.Monday {
				t.Errorf("StartOfWeek() should return Monday, got %v", result.Weekday())
			}

			// Verify time is set to start of day
			if result.Hour() != 0 || result.Minute() != 0 || result.Second() != 0 || result.Nanosecond() != 0 {
				t.Errorf("StartOfWeek() should return start of day, got %v", result.Time)
			}
		})
	}
}

func TestChronos_EndOfWeek(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected time.Time
	}{
		{
			name:     "Monday",
			input:    time.Date(2025, 10, 20, 15, 30, 45, 0, time.UTC),         // Monday
			expected: time.Date(2025, 10, 26, 23, 59, 59, 999999999, time.UTC), // Sunday
		},
		{
			name:     "Tuesday",
			input:    time.Date(2025, 10, 21, 10, 0, 0, 0, time.UTC),           // Tuesday
			expected: time.Date(2025, 10, 26, 23, 59, 59, 999999999, time.UTC), // Sunday
		},
		{
			name:     "Wednesday",
			input:    time.Date(2025, 10, 22, 20, 45, 30, 0, time.UTC),         // Wednesday
			expected: time.Date(2025, 10, 26, 23, 59, 59, 999999999, time.UTC), // Sunday
		},
		{
			name:     "Thursday",
			input:    time.Date(2025, 10, 23, 8, 15, 0, 0, time.UTC),           // Thursday
			expected: time.Date(2025, 10, 26, 23, 59, 59, 999999999, time.UTC), // Sunday
		},
		{
			name:     "Friday",
			input:    time.Date(2025, 10, 24, 18, 30, 0, 0, time.UTC),          // Friday
			expected: time.Date(2025, 10, 26, 23, 59, 59, 999999999, time.UTC), // Sunday
		},
		{
			name:     "Saturday",
			input:    time.Date(2025, 10, 25, 12, 0, 0, 0, time.UTC),           // Saturday
			expected: time.Date(2025, 10, 26, 23, 59, 59, 999999999, time.UTC), // Sunday
		},
		{
			name:     "Sunday",
			input:    time.Date(2025, 10, 26, 23, 59, 59, 0, time.UTC),         // Sunday
			expected: time.Date(2025, 10, 26, 23, 59, 59, 999999999, time.UTC), // Same day at end
		},
		{
			name:     "Week spanning month boundary",
			input:    time.Date(2025, 9, 30, 12, 0, 0, 0, time.UTC),           // Tuesday
			expected: time.Date(2025, 10, 5, 23, 59, 59, 999999999, time.UTC), // Sunday
		},
		{
			name:     "Week spanning year boundary",
			input:    time.Date(2025, 12, 31, 12, 0, 0, 0, time.UTC),         // Last day of year
			expected: time.Date(2026, 1, 4, 23, 59, 59, 999999999, time.UTC), // New year Sunday
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Chronos{Time: tt.input}
			result := c.EndOfWeek()

			if !result.Time.Equal(tt.expected) {
				t.Errorf("EndOfWeek() for %v = %v, expected %v", tt.input, result.Time, tt.expected)
			}

			// Verify it's indeed Sunday
			if result.Weekday() != time.Sunday {
				t.Errorf("EndOfWeek() should return Sunday, got %v", result.Weekday())
			}

			// Verify time is set to end of day
			if result.Hour() != 23 || result.Minute() != 59 || result.Second() != 59 || result.Nanosecond() != 999999999 {
				t.Errorf("EndOfWeek() should return end of day, got %v", result.Time)
			}
		})
	}
}

func TestParseSameFormat(t *testing.T) {
	strt1 := ParseSameFormat(DateOnly, "2025-10-23")
	assert.Equal(t, "2025-10-23", strt1)

	strt2 := ParseSameFormat(DateOnly, "Invalid")
	assert.Empty(t, strt2)
}

func TestWeekdayRu(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Monday",
			input: "2025-10-20",
			want:  "пн",
		},
		{
			name:  "Tuesday",
			input: "2025-10-21",
			want:  "вт",
		},
		{
			name:  "Wednesday",
			input: "2025-10-22",
			want:  "ср",
		},
		{
			name:  "Thursday",
			input: "2025-10-23",
			want:  "чт",
		},
		{
			name:  "Friday",
			input: "2025-10-24",
			want:  "пт",
		},
		{
			name:  "Saturday",
			input: "2025-10-25",
			want:  "сб",
		},
		{
			name:  "Sunday",
			input: "2025-10-26",
			want:  "вс",
		},
		{
			name:  "Monday",
			input: "2025-10-27",
			want:  "пн",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch, err := FromString(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.want, ch.WeekdayRu(), "WeekdayRu mismatch")
		})
	}
}

func TestFormatDateRu(t *testing.T) {
	ch, err := FromString("2025-01-14")
	require.NoError(t, err)
	assert.Equal(t, "14 января", FormatDateRu(ch))
}

func TestUTC(t *testing.T) {
	ch, err := FromString("Mon Jan 02 2006 15:04:05 GMT-0700 (Coordinated Universal Time)")
	require.NoError(t, err)
	assert.Equal(t, "2006-01-02T22:04:05Z", ch.UTC().String())
}

func TestDayIndex(t *testing.T) {
	tests := []struct {
		day  string
		want int
	}{
		{"пн", 0},
		{"вт", 1},
		{"ср", 2},
		{"чт", 3},
		{"пт", 4},
		{"сб", 5},
		{"вс", 6},
		{"неизвестный", -1},
		{"", -1},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("dayIndex(%s)", tt.day), func(t *testing.T) {
			got := WeekdayRuIndex(tt.day)
			if got != tt.want {
				t.Errorf("dayIndex(%q) = %d, want %d", tt.day, got, tt.want)
			}
		})
	}
}

func TestInZone(t *testing.T) {
	// minus 1 hour
	ch, err := FromString("2026-04-09T22:13:05+04:00")
	require.NoError(t, err)
	assert.Equal(t, "2026-04-09 21:13:05", ch.InLocOffset(3).Format(DateTime))

	// utc zone
	ch, err = FromString("2026-04-09 15:04:05")
	require.NoError(t, err)
	assert.Equal(t, "2026-04-09 19:04:05", ch.InLocOffset(4).Format(DateTime))

	// same zone
	ch, err = FromString("2026-04-09T22:13:05+04:00")
	require.NoError(t, err)
	assert.Equal(t, "2026-04-09 22:13:05", ch.InLocOffset(4).Format(DateTime))
}

func TestSameChronosInUTC(t *testing.T) {
	ch, err := FromString("2026-04-09T22:13:05+04:00")
	require.NoError(t, err)
	assert.Equal(t, "2026-04-09T22:13:05Z", SameChronosInUTC(ch).String())

	ch, err = FromString("2026-04-09 15:04:05")
	require.NoError(t, err)
	assert.Equal(t, "2026-04-09T15:04:05Z", SameChronosInUTC(ch).String())
}

func TestSameChronosInLoc(t *testing.T) {
	ch, err := FromString("2026-06-24T18:00:00Z")
	require.NoError(t, err)

	assert.Equal(t, "2026-06-24T18:00:00+04:00", SameChronosInLoc(ch, 4).String())
}
