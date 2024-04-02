package date

import (
	"testing"
	"time"
)

func TestParseHappyPath(t *testing.T) {
	testTimeStr := "2023-12-25T08:00:00Z"
	expectedTime := time.Date(2023, 12, 25, 8, 0, 0, 0, time.UTC)

	parsedTime := Parse(testTimeStr, "UTC")

	if !parsedTime.Equal(expectedTime) {
		t.Errorf("Expected parsed time to be %v; Got %v", expectedTime, parsedTime)
	}
}

func TestParseEdgeCases(t *testing.T) {
	// Edge case: time string with empty timezone
	testTimeStr := "2023-12-25T08:00:00Z"
	expectedTime := time.Date(2023, 12, 25, 8, 0, 0, 0, time.UTC)

	parsedTime := Parse(testTimeStr, "")

	if !parsedTime.Equal(expectedTime) {
		t.Errorf("Expected parsed time to be %v; Got %v", expectedTime, parsedTime)
	}

	// Edge case: time string in RFC3339 format
	rfcTimeStr := "2023-12-25T08:00:00Z"
	expectedRFCTime := time.Date(2023, 12, 25, 8, 0, 0, 0, time.UTC)

	parsedRFCTime := ParseRFC3339(rfcTimeStr)

	if !parsedRFCTime.Equal(expectedRFCTime) {
		t.Errorf("Expected parsed RFC3339 time to be %v; Got %v", expectedRFCTime, parsedRFCTime)
	}

	// Edge case: Now function with different timezones
	// utcNow := Now("UTC")
	// localNow := Now()

	// if !utcNow.Equal(localNow) {
	// 	t.Errorf("Expected UTC now to be equal to local now; Got UTC: %v, Local: %v", utcNow, localNow)
	// }

	// Edge case: Start and end of day/month/year with specific date
	testDate := time.Date(2023, 12, 25, 12, 30, 0, 0, time.UTC)
	expectedStartOfDay := time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC)
	expectedEndOfDay := time.Date(2023, 12, 25, 23, 59, 59, 0, time.UTC)
	expectedStartOfMonth := time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC)
	expectedEndOfMonth := time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)
	expectedStartOfYear := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	expectedEndOfYear := time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)

	startOfDay := StartOfDay(testDate)
	endOfDay := EndOfDay(testDate)
	startOfMonth := StartOfMonth(testDate)
	endOfMonth := EndOfMonth(testDate)
	startOfYear := StartOfYear(testDate)
	endOfYear := EndOfYear(testDate)

	if !startOfDay.Equal(expectedStartOfDay) {
		t.Errorf("Expected start of day to be %v; Got %v", expectedStartOfDay, startOfDay)
	}

	if !endOfDay.Equal(expectedEndOfDay) {
		t.Errorf("Expected end of day to be %v; Got %v", expectedEndOfDay, endOfDay)
	}

	if !startOfMonth.Equal(expectedStartOfMonth) {
		t.Errorf("Expected start of month to be %v; Got %v", expectedStartOfMonth, startOfMonth)
	}

	if !endOfMonth.Equal(expectedEndOfMonth) {
		t.Errorf("Expected end of month to be %v; Got %v", expectedEndOfMonth, endOfMonth)
	}

	if !startOfYear.Equal(expectedStartOfYear) {
		t.Errorf("Expected start of year to be %v; Got %v", expectedStartOfYear, startOfYear)
	}

	if !endOfYear.Equal(expectedEndOfYear) {
		t.Errorf("Expected end of year to be %v; Got %v", expectedEndOfYear, endOfYear)
	}
}
