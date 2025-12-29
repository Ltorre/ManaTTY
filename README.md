# 🏰 ManaTTY - Mage Tower Ascension

A heroic fantasy terminal game with idle/incremental mechanics, built with Go and Bubble Tea TUI.

## 🎮 Game Concept

You are a wizard climbing a magical tower! Cast spells to earn mana/experience and ascend through floors. Each floor unlocks new spells that can be cast automatically. Combine spells into powerful rituals (3-spell combos) to boost your progression. When you reach the top, prestige to gain permanent bonuses and start again stronger!

## 📋 Project Status

**Current Phase:** Rituals & Prestige System Complete

### Milestones

- [x] **Milestone 1:** Project initialization & structure
- [x] **Milestone 2:** Core data models
- [x] **Milestone 3:** Game constants & formulas
- [x] **Milestone 4:** MongoDB storage layer
- [x] **Milestone 5:** Game engine core
- [x] **Milestone 6:** Bubble Tea UI foundation
- [x] **Milestone 7:** Main integration & Tower/Spell views
- [x] **Milestone 8:** Rituals & Prestige system
- [ ] **Milestone 9:** Offline progress & polish

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
- **Floor Climbing:** Spend mana to ascend to higher floors
- **Spells:** Unlock and cast spells with various cooldowns
- **Rituals:** Combine 3 spells for +15% mana generation per ritual
- **Prestige:** Reset at floor 100 for permanent multipliers

## 📝 License

MIT License

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
