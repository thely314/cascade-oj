package util

import (
	"errors"
	"time"
)

// TimeZone 时区类型
type TimeZone int

// 时区映射，从 UTC 偏移量（小时）到时区字符串
var offsetToTimeZone = map[int]string{
	-12: "Pacific/Fiji",         // -12
	-11: "Pacific/Midway",       // -11
	-10: "Pacific/Honolulu",     // -10
	-9:  "America/Anchorage",    // -9
	-8:  "America/Los_Angeles",  // -8
	-7:  "America/Denver",       // -7
	-6:  "America/Chicago",      // -6
	-5:  "America/New_York",     // -5
	-4:  "America/Sao_Paulo",    // -4
	-3:  "America/Buenos_Aires", // -3
	-2:  "America/Noronha",      // -2
	-1:  "Atlantic/Azores",      // -1
	0:   "UTC",                  // 0
	1:   "Europe/Paris",         // +1
	2:   "Europe/Istanbul",      // +2
	3:   "Europe/Moscow",        // +3
	4:   "Asia/Dubai",           // +4
	5:   "Asia/Karachi",         // +5
	6:   "Asia/Dhaka",           // +6
	7:   "Asia/Bangkok",         // +7
	8:   "Asia/Shanghai",        // +8
	9:   "Asia/Tokyo",           // +9
	10:  "Australia/Sydney",     // +10
	11:  "Pacific/Guadalcanal",  // +11
	12:  "Pacific/Auckland",     // +12
}

// GetTimeZone 获取时区对象
func GetTimeZone(zone TimeZone) (time.Location, error) {
	// 检查参数是否合法
	offset := int(zone)
	if _, ok := offsetToTimeZone[offset]; !ok {
		return time.Location{}, errors.New("invalid time zone offset")
	}

	location, err := time.LoadLocation(offsetToTimeZone[offset])
	if err != nil {
		return time.Location{}, err
	}

	return *location, nil
}

// ConvertFrontendTimeToDatabaseTime 将前端的本地时间转换为数据库的 UTC 时间
func ConvertFrontendTimeToDatabaseTime(t time.Time, frontendZone TimeZone) (time.Time, error) {
	// 获取前端时区
	frontendLocation, err := GetTimeZone(frontendZone)
	if err != nil {
		return time.Time{}, err
	}

	// 将时间转换为前端时区，再转换为 UTC
	localTime := t.In(&frontendLocation)
	return localTime.UTC(), nil
}

// ConvertDatabaseTimeToFrontendTime 将数据库的 UTC 时间转换为前端的本地时间
func ConvertDatabaseTimeToFrontendTime(t time.Time, frontendZone TimeZone) (time.Time, error) {
	// 获取前端时区
	frontendLocation, err := GetTimeZone(frontendZone)
	if err != nil {
		return time.Time{}, err
	}

	// 将 UTC 时间转换为前端时区
	return t.In(&frontendLocation), nil
}
