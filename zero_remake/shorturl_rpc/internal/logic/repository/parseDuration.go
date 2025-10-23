package repository

import (
	"strconv"
	"strings"
	"time"
)

type DurationType struct{}

// parseCustomDuration 自定义函数，支持解析包含 'd' 单位的时间字符串
func (DurationType *DurationType) ParseCustomDuration(s string) (time.Duration, error) {
	if strings.HasSuffix(s, "d") {
		daysStr := strings.TrimSuffix(s, "d")
		days, err := strconv.Atoi(daysStr)
		if err != nil {
			return 0, err
		}
		return time.Duration(days) * 24 * time.Hour, nil
	}
	return time.ParseDuration(s)
}

var Duration = DurationType{}
