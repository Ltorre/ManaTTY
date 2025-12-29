# 🏰 ManaTTY - Mage Tower Ascension

A heroic fantasy terminal game with idle/incremental mechanics, built with Go and Bubble Tea TUI.

## 🎮 Game Concept

You are a wizard climbing a magical tower! Cast spells to deal damage and earn mana to ascend through floors. Each floor unlocks new spells that can be cast automatically. Combine spells into powerful rituals (3-spell combos) to boost your progression. When you reach the top, prestige to gain permanent bonuses and start again stronger!

**New:** Enter a nickname at launch to support multiple local save files!

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
- [x] **Milestone 10:** Buildcrafting & QoL improvements
  - Spell leveling (max 10) with damage/cooldown/cost scaling
  - Auto-cast slot priority reordering
  - Element synergy system (3 same-element casts = 20% buff)
  - Aggregated mana hints in loadout panel
  - Mana display in spells view
  - Nickname prompt for multiple save files
- [x] **Milestone 11:** Ascension Sigil system
  - Damage-based gate for floor climbing (mana + sigil required)
  - Spell damage stats visible in UI
  - Balanced scaling for mid/late game challenge
- [x] **Milestone 12:** Spell Specialization & Auto-Cast Rules
  - Tier 1 specializations at level 5: Crit Chance (+15% for 2x dmg) or Mana Efficiency (-20% cost)
  - Tier 2 specializations at level 10: Burst Damage (+30% dmg) or Rapid Cast (-25% CD)
  - Auto-cast conditional rules: Always, Mana>50%, Mana>75%, Sigil not full, Synergy active
  - Press `X` to specialize, `C` to cycle conditions
- [x] **Milestone 13 (v1.1.0):** Elemental Resonance Loadout Bonus
  - 2+ spells of the same element in auto-cast grants a passive perk while equipped
  - Encourages themed loadouts and element-focused buildcrafting
- [x] **Milestone 14 (v1.1.0):** Floor Events (Lightweight)
  - Every 25 floors, a timed choice appears: +mana/sec (10 floors) vs +sigil charge rate (10 floors) vs -cooldowns (10 floors)
  - If you don’t choose within 2 minutes, the event vanishes with no bonus (prevents idle lock-ups)
- [x] **Milestone 15 (v1.2.0):** Named Ritual Combos & Effects
  - Each 3-spell ritual now has a unique generated name based on element composition
  - Rituals grant passive bonuses: Pure (3 same element) = +18%, Hybrid (2+1) = +12%, Triad (1/1/1) = +8% each
  - 🔥 Fire rituals boost spell damage, ❄️ Ice reduces cooldowns, ⚡ Thunder reduces mana cost, ✨ Arcane boosts sigil charge
  - Including Spell Echo in a ritual adds a +5% kicker to all effects
  - Special "signature" combos have flavor names (Elemental Trinity, Apocalypse, Resonant Fire, etc.)
  - Live preview of ritual name and effects when selecting spells

### Future Roadmap

- [ ] **(Maybe) - Ritual Mastery & Evolution**
  - Casting the same ritual combo repeatedly "masters" it
  - Mastered rituals grant minor permanent effects or reduced cooldowns
  - Adds long-term goals without new UI complexity

- [ ] **(Maybe) - Loadout Perks (Tiny Set Bonuses)**
  - Simple 2-piece / 3-piece bonuses across elements (e.g., Fire+Thunder = +sigil rate)
  - Always-on while equipped; no extra clicks

- [ ] **(Maybe) - Ritual Presets / Favorites**
  - Mark a ritual as a favorite and quickly re-create it after prestige
  - Keeps the loop fast without adding new gameplay burden

- [ ] **(Maybe) - Auto-Cast “Focus Mode”**
  - Optional toggle: prefer casting spells matching active synergy (when possible)
  - Reduces micro-management while encouraging cohesive builds

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
- **Ascension Sigil:** Deal damage by casting spells to charge the sigil—both mana AND sigil must be full to climb!
- **Floor Climbing:** Spend mana to ascend to higher floors (cost scales exponentially)
- **Spells:** Unlock and cast 12 unique spells across 4 elements (Fire, Ice, Thunder, Arcane)
- **Spell Leveling:** Spend mana to upgrade spells (max level 10) for reduced cooldown, lower mana cost, and increased damage
- **Spell Specialization:** At levels 5 and 10, choose a specialization path (Crit, Mana Efficiency, Burst Damage, or Rapid Cast)
- **Auto-Cast Slots:** Assign spells to limited auto-cast slots (base 2, up to 4 with prestige)
- **Auto-Cast Conditions:** Set rules per slot: Always, Mana>50%, Mana>75%, Sigil not full, or Synergy active
- **Slot Priority:** Reorder auto-cast slots to control which spells cast first
- **Element Synergies:** Cast 3 spells of the same element in a row for a 10-second buff (20% reduced cost, cooldown & +20% damage)
- **Mana Economy:** All spells (auto and manual) consume mana—choose your auto-cast loadout wisely!
- **Manual Casting:** Cast any spell manually (+10% mana cost) for tactical control
- **Rituals:** Combine 3 spells for +15% mana generation per ritual
- **Prestige:** Reset at floor 100 for permanent multipliers, more ritual slots, and more auto-cast slots
- **Offline Progress:** Earn mana even while away (50% efficiency)

## ⌨️ Controls

| Key | Action |
|-----|--------|
| `S` | Open Spells view |
| `R` | Open Rituals view |
| `T` | Open Stats view |
| `P` | Open Prestige view (at floor 100+) |
| `M` | Open Menu |
| `A` | Toggle Auto-cast on/off |
| `Space` | Toggle spell in auto-cast slot (Spells view) |
| `U` | Upgrade selected spell (Spells view) |
| `X` | Open specialization menu (Spells view, level 5/10) |
| `C` | Cycle auto-cast condition (Spells view) |
| `<` / `>` | Reorder auto-cast slot priority (Spells view) |
| `↑/↓` | Navigate lists |
| `Enter` | Select/Cast spell manually |
| `Ctrl+S` | Manual Save |
| `Q` | Quit (auto-saves) |

## 📝 License

MIT License

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
