package components

import (
	"time"

	"github.com/charmbracelet/lipgloss"
)

// NotificationType indicates the severity of a notification.
type NotificationType string

const (
	NotifyInfo    NotificationType = "info"
	NotifySuccess NotificationType = "success"
	NotifyWarning NotificationType = "warning"
	NotifyError   NotificationType = "error"
)

// Notification represents a toast notification.
type Notification struct {
	Text      string
	Type      NotificationType
	CreatedAt time.Time
	Duration  time.Duration
}

// Notification styles
var (
	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3B82F6")).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#22C55E")).
			Bold(true)

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EAB308")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EF4444")).
			Bold(true)
)

// NewNotification creates a new notification.
func NewNotification(text string, notifyType NotificationType) *Notification {
	return &Notification{
		Text:      text,
		Type:      notifyType,
		CreatedAt: time.Now(),
		Duration:  3 * time.Second,
	}
}

// IsExpired returns true if the notification has expired.
func (n *Notification) IsExpired() bool {
	return time.Since(n.CreatedAt) > n.Duration
}

// Render renders the notification.
func (n *Notification) Render() string {
	var style lipgloss.Style
	var icon string

	switch n.Type {
	case NotifySuccess:
		style = successStyle
		icon = "✓"
	case NotifyWarning:
		style = warningStyle
		icon = "⚠"
	case NotifyError:
		style = errorStyle
		icon = "✗"
	default:
		style = infoStyle
		icon = "ℹ"
	}

	return style.Render(icon + " " + n.Text)
}

// RenderNotification renders a notification string with type.
func RenderNotification(text string, notifyType NotificationType) string {
	n := NewNotification(text, notifyType)
	return n.Render()
}
