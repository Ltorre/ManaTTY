package utils

import (
	"fmt"
	"time"
)

// FormatDuration formats a duration as human-readable.
func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0fs", d.Seconds())
	}
	if d < time.Hour {
		minutes := int(d.Minutes())
		seconds := int(d.Seconds()) % 60
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	if d < 24*time.Hour {
		hours := int(d.Hours())
		minutes := int(d.Minutes()) % 60
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	return fmt.Sprintf("%dd %dh", days, hours)
}

// FormatDurationShort formats a duration in compact form.
func FormatDurationShort(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.0fm", d.Minutes())
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%.1fh", d.Hours())
	}
	return fmt.Sprintf("%.1fd", d.Hours()/24)
}

// FormatMilliseconds formats milliseconds as a duration.
func FormatMilliseconds(ms int64) string {
	return FormatDuration(time.Duration(ms) * time.Millisecond)
}

// FormatCooldown formats cooldown remaining in seconds.
func FormatCooldown(ms int64) string {
	if ms <= 0 {
		return "Ready!"
	}
	seconds := float64(ms) / 1000.0
	if seconds < 10 {
		return fmt.Sprintf("%.1fs", seconds)
	}
	return fmt.Sprintf("%.0fs", seconds)
}

// TimeSince returns the duration since a time as formatted string.
func TimeSince(t time.Time) string {
	return FormatDuration(time.Since(t))
}

// TimeUntil returns the duration until a time as formatted string.
func TimeUntil(t time.Time) string {
	d := time.Until(t)
	if d < 0 {
		return "now"
	}
	return FormatDuration(d)
}

// NowMillis returns current time in milliseconds.
func NowMillis() int64 {
	return time.Now().UnixMilli()
}

// MillisToTime converts milliseconds to time.Time.
func MillisToTime(ms int64) time.Time {
	return time.UnixMilli(ms)
}

// CalculateOfflineTime returns seconds since last save.
func CalculateOfflineTime(lastSaved time.Time) float64 {
	return time.Since(lastSaved).Seconds()
}
