package ui

import (
	"os"
	"runtime"
)

// Symbols holds all the icons/emojis used in the UI.
// Modern terminals get emojis, legacy CMD gets ASCII fallbacks.
type Symbols struct {
	// Elements
	Fire    string
	Ice     string
	Thunder string
	Arcane  string
	Default string

	// UI icons
	Tower    string
	Stats    string
	Ritual   string
	AutoCast string
	Event    string
	Mana     string
	Floor    string
	Damage   string
	Cooldown string
	Wave     string
	Prestige string
	Synergy  string

	// Status
	Ready     string
	Cooldown2 string
	Active    string
	Locked    string

	// Misc
	Bullet string
	Arrow  string
	Star   string
	Check  string
	Cross  string
}

// symbols is the global symbol set, initialized based on terminal capabilities.
var symbols Symbols

// supportsEmoji detects if the current terminal supports emoji display.
// Returns true for:
// - All non-Windows systems (macOS, Linux, etc.)
// - Windows Terminal (WT_SESSION env var)
// - VS Code integrated terminal (TERM_PROGRAM=vscode)
// - PowerShell in modern terminals
// Returns false for legacy CMD.exe
func supportsEmoji() bool {
	// Non-Windows always supports emoji
	if runtime.GOOS != "windows" {
		return true
	}

	// Windows Terminal sets WT_SESSION
	if os.Getenv("WT_SESSION") != "" {
		return true
	}

	// VS Code terminal
	if os.Getenv("TERM_PROGRAM") == "vscode" {
		return true
	}

	// ConEmu/Cmder set ConEmuANSI
	if os.Getenv("ConEmuANSI") == "ON" {
		return true
	}

	// Alacritty, Hyper, etc. set TERM
	if term := os.Getenv("TERM"); term != "" && term != "dumb" {
		return true
	}

	// ANSICON for some modern CMD setups
	if os.Getenv("ANSICON") != "" {
		return true
	}

	// Default: assume legacy CMD on Windows = no emoji
	return false
}

func init() {
	if supportsEmoji() {
		// Full emoji support
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
	} else {
		// ASCII fallbacks for legacy Windows CMD
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
	}
}

// GetSymbols returns the current symbol set.
func GetSymbols() Symbols {
	return symbols
}

// SupportsEmoji returns true if the terminal supports emoji display.
func SupportsEmoji() bool {
	return supportsEmoji()
}
