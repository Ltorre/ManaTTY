# 🏰 ManaTTY - Mage Tower Ascension

A heroic fantasy terminal game with idle/incremental mechanics, built with Go and Bubble Tea TUI.

## 🎮 Game Concept

You are a wizard climbing a magical tower! Cast spells to earn mana/experience and ascend through floors. Each floor unlocks new spells that can be cast automatically. Combine spells into powerful rituals (3-spell combos) to boost your progression. When you reach the top, prestige to gain permanent bonuses and start again stronger!

## 📋 Project Status

**Current Phase:** ✅ Complete - Ready to Play!

### Milestones

- [x] **Milestone 1:** Project initialization & structure
- [x] **Milestone 2:** Core data models
- [x] **Milestone 3:** Game constants & formulas
- [x] **Milestone 4:** MongoDB storage layer
- [x] **Milestone 5:** Game engine core
- [x] **Milestone 6:** Bubble Tea UI foundation
- [x] **Milestone 7:** Main integration & Tower/Spell views
- [x] **Milestone 8:** Rituals & Prestige system
- [x] **Milestone 9:** Offline progress & polish

## 🛠️ Tech Stack

- **Language:** Go 1.21+
- **TUI Framework:** [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **Styling:** [Lipgloss](https://github.com/charmbracelet/lipgloss)
- **Database:** MongoDB (Atlas Cloud or local)

## 📁 Project Structure

```
mage-tower-ascension/
├── main.go                 # Entry point
├── config/                 # Configuration management
├── models/                 # Data models (Game, Player, Spell, etc.)
├── storage/                # MongoDB connection & repositories
├── engine/                 # Game logic & calculations
├── ui/                     # Bubble Tea TUI components
│   ├── screens/            # Individual view screens
│   └── components/         # Reusable UI components
├── game/                   # Game constants & formulas
└── utils/                  # Helper utilities
```

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or later
- MongoDB instance (local or Atlas)

### Installation

```bash
# Clone the repository
git clone https://github.com/Ltorre/ManaTTY.git
cd ManaTTY

# Install dependencies
go mod download

# Set up environment variables
cp .env.example .env
# Edit .env with your MongoDB URI

# Run the game
go run main.go
```

## ⚙️ Configuration

Create a `.env` file in the project root:

```env
MONGODB_URI=mongodb://localhost:27017/mage_tower
LOG_LEVEL=info
GAME_TICK_RATE=10
AUTO_SAVE_INTERVAL=30
DEBUG=false
```

## 🎯 Core Mechanics

- **Mana Generation:** Earn mana passively based on your current floor
- **Floor Climbing:** Spend mana to ascend to higher floors (cost scales exponentially)
- **Spells:** Unlock and cast 12 unique spells across 4 elements (Fire, Ice, Thunder, Arcane)
- **Auto-Cast:** Toggle automatic spell casting for hands-free progression
- **Rituals:** Combine 3 spells for +15% mana generation per ritual
- **Prestige:** Reset at floor 100 for permanent multipliers and bonuses
- **Offline Progress:** Earn mana even while away (50% efficiency)

## ⌨️ Controls

| Key | Action |
|-----|--------|
| `S` | Open Spells view |
| `R` | Open Rituals view |
| `T` | Open Stats view |
| `P` | Open Prestige view (at floor 100+) |
| `M` | Open Menu |
| `A` | Toggle Auto-cast |
| `↑/↓` | Navigate lists |
| `Enter` | Select/Cast |
| `Ctrl+S` | Manual Save |
| `Q` | Quit (auto-saves) |

## 📝 License

MIT License

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
