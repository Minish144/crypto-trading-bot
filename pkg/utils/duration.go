package utils

import "time"

func DurationToString(duration time.Duration) string {
	minutes := duration.Minutes()
	switch {
	case minutes <= 1:
		return "1min"
	case minutes <= 2:
		return "2min"
	case minutes <= 3:
		return "3min"
	case minutes <= 5:
		return "5min"
	case minutes <= 10:
		return "10min"
	case minutes <= 15:
		return "15min"
	case minutes <= 30:
		return "30min"
	case minutes <= 60:
		return "hour"
	case minutes <= 1440:
		return "day"
	case minutes <= 10080:
		return "week"
	default:
		return "month"
	}
}
