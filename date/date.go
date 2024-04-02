package date

import "time"

// Parse 解析时间字符串，它接受一个时间字符串和一个可选的时区参数，返回解析后的时间对象。
func Parse(timeStr string, timezone ...string) time.Time {
	// 如果传入的时区参数数量大于0并且不为空
	if len(timezone) > 0 && timezone[0] != "" {
		// 使用指定时区解析时间字符串
		loc, _ := time.LoadLocation(timezone[0])
		t, err := time.ParseInLocation(time.DateTime, timeStr, loc)
		// 如果出现解析错误，则尝试使用其他时间格式进行解析
		if err != nil {
			t, err = time.Parse(time.RFC3339, timeStr)
			if err != nil {
				t, _ = time.Parse(time.DateOnly, timeStr)

				return t
			}

			return t
		}

		return t
	}

	// 如果未提供时区参数或者时区参数为空，则使用默认时区进行解析
	t, err := time.Parse(time.DateTime, timeStr)
	if err != nil {
		t, err = time.Parse(time.DateOnly, timeStr)
		if err != nil {
			t, _ = time.Parse(time.RFC3339, timeStr)

			return t
		}

		return t
	}

	return t
}

// ParseRFC3339 根据RFC3339格式解析时间字符串，可指定时区
func ParseRFC3339(timeStr string, timezone ...string) time.Time {
	if len(timezone) > 0 && timezone[0] != "" {
		loc, _ := time.LoadLocation(timezone[0])
		t, _ := time.ParseInLocation(time.RFC3339, timeStr, loc)

		return t
	}

	t, _ := time.Parse(time.RFC3339, timeStr)

	return t
}

// Now 函数返回当前时间，可以根据传入的时区参数进行时区转换
// 如果传入时区参数，且不为空，则根据指定的时区返回当前时间
// 否则返回本地时间的当前时间
func Now(timezone ...string) time.Time {
	if len(timezone) > 0 && timezone[0] != "" {
		loc, _ := time.LoadLocation(timezone[0])

		return time.Now().In(loc)
	}

	return time.Now()
}

// StartOfDay 函数接受一个时间 t，返回该时间对应的当天零点时间
func StartOfDay(t time.Time) time.Time {
	y, m, d := t.Date()

	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

// StartOfMonth 返回给定时间所在月份的第一天
func StartOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()

	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
}

// StartOfQuarter 函数接受一个时间参数 t，返回 t 所在季度的第一天的时间
func StartOfQuarter(t time.Time) time.Time {
	return time.Date(t.Year(), time.Month(((int(t.Month())-1)/3+1-1)*3+1), 1, 0, 0, 0, 0, t.Location())
}

// StartOfYear 函数接受一个时间参数 t，返回该年的第一天的开始时间
func StartOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
}

// EndOfDay 函数接受一个时间参数，返回当天的最后一秒时间
func EndOfDay(t time.Time) time.Time {
	y, m, d := t.Date()

	return time.Date(y, m, d, 23, 59, 59, 0, t.Location())
}

// EndOfMonth 函数接受一个时间参数，返回当月的最后一秒时间
func EndOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()

	return time.Date(year, month+1, 0, 23, 59, 59, 0, t.Location())
}

// EndOfYear 函数接受一个时间参数，返回当年的最后一秒时间
func EndOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 12, 31, 23, 59, 59, 0, t.Location())
}

// IsSameDay 检查两个时间 t1 和 t2 是否在同一天
func IsSameDay(t1, t2 time.Time) bool {
	return t1.YearDay() == t2.YearDay() && t1.Year() == t2.Year()
}

// IsSameMonth 检查两个时间是否是同一个月份
func IsSameMonth(t1, t2 time.Time) bool {
	t2 = t2.In(t1.Location())

	return t1.Year() == t2.Year() && t1.Month() == t2.Month()
}

// IsSameYear 检查两个时间是否是同一年
func IsSameYear(t1, t2 time.Time) bool {
	t2 = t2.In(t1.Location())

	return t1.Year() == t2.Year()
}

// IsToday 判断给定的时间 t 是否为今天
func IsToday(t time.Time) bool {
	return IsSameDay(t, time.Now().In(t.Location()))
}

// DiffInSecond 函数接受两个时间参数并返回它们之间的差异（以秒为单位）。
func DiffInSecond(t1, t2 time.Time) int {
	return int(t2.Unix() - t1.Unix())
}

// DiffAbsInSecond 函数接受两个时间参数并返回它们之间的绝对差异（以秒为单位）。
// 如果差异为负数，则返回其绝对值。
func DiffAbsInSecond(t1, t2 time.Time) int {
	diff := DiffInSecond(t1, t2)
	if diff < 0 {
		return -diff
	}

	return diff
}
