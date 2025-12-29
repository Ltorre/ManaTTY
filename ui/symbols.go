package ui

import (
	"runtime"
)

// Symbols holds all the icons/emojis used in the UI.
// On Windows, ASCII fallbacks are used for better compatibility.
type Symbols struct {
	// Elements
	Fire    string
	Ice     string
	Thunder string
	Arcane  string
	Default string

	// UI icons
	Tower     string
	Stats     string
	Ritual    string
	AutoCast  string
	Event     string
	Mana      string
	Floor     string
	Damage    string
	Cooldown  string
	Wave      string
	Prestige  string
	Synergy   string

	// Status
	Ready     string
	Cooldown2 string
	Active    string
	Locked    string

	// Misc
	Bullet    string
	Arrow     string
	Star      string
	Check     string
	Cross     string
}

// symbols is the global symbol set, initialized based on OS.
var symbols Symbols

func init() {
	if runtime.GOOS == "windows" {
		// ASCII fallbacks for Windows Command Prompt compatibility
		symbols = Symbols{
			// Elements - use colored text markers instead of emojis
			Fire:    "[F]",
			Ice:     "[I]",
			Thunder: "[T]",
			Arcane:  "[A]",
			Default: "*",

			// UI icons
			Tower:    "#",
			Stats:    "*",
			Ritual:   "~",
			AutoCast: ">",
			Event:    "!",
			Mana:     "o",
			Floor:    "^",
			Damage:   "!",
			Cooldown: "-",
			Wave:     "~",
			Prestige: "+",
			Synergy:  "*",

			// Status
			Ready:     "+",
			Cooldown2: "-",
			Active:    "*",
			Locked:    "x",

			// Misc
			Bullet: "*",
			Arrow:  ">",
			Star:   "*",
			Check:  "+",
			Cross:  "x",
		}
	} else {
		// Full emoji support for Unix-like systems
		symbols = Symbols{
			// Elements
			Fire:    "🔥",
			Ice:     "❄️",
			Thunder: "⚡",
			Arcane:  "✨",
			Default: "•",

			// UI icons
			Tower:    "🏰",
			Stats:    "📊",
			Ritual:   "🔥",
			AutoCast: "⚡",
			Event:    "✨",
			Mana:     "💎",
			Floor:    "🏔️",
			Damage:   "💥",
			Cooldown: "⏱️",
			Wave:     "🌊",
			Prestige: "⭐",
			Synergy:  "🔥",

			// Status
			Ready:     "✓",
			Cooldown2: "⏳",
			Active:    "●",
			Locked:    "🔒",

			// Misc
			Bullet: "•",
			Arrow:  "→",
			Star:   "★",
			Check:  "✓",
			Cross:  "✗",
		}
	}
}

// GetSymbols returns the current symbol set.
func GetSymbols() Symbols {
	return symbols
}

// IsWindows returns true if running on Windows.
func IsWindows() bool {
	return runtime.GOOS == "windows"
}
