# 🏰 ManaTTY - Mage Tower Ascension

A heroic fantasy terminal game with idle/incremental mechanics, built with Go and Bubble Tea TUI.

## 🎮 Game Concept

You are a wizard climbing a magical tower! Cast spells to deal damage and earn mana to ascend through floors. Each floor unlocks new spells that can be cast automatically. Combine spells into powerful rituals (3-spell combos) to boost your progression. When you reach the top, prestige to gain permanent bonuses and start again stronger!

**New:** Enter a nickname at launch to support multiple local save files!

## 📋 Project Status

**Current Phase:** ✅ Complete - Ready to Play!

### Milestones


#### v1.0.0 — Core Game
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

#### v1.1.0 — Events & Loadout Bonuses
- [x] **Milestone 13:** Elemental Resonance Loadout Bonus
  - 2+ spells of the same element in auto-cast grants a passive perk while equipped
  - Encourages themed loadouts and element-focused buildcrafting
- [x] **Milestone 14:** Floor Events (Lightweight)
  - Every 25 floors, a timed choice appears: +mana/sec (10 floors) vs +sigil charge rate (10 floors) vs -cooldowns (10 floors)
  - If you don't choose within 2 minutes, the event vanishes with no bonus (prevents idle lock-ups)

#### v1.2.0 — Named Rituals
- [x] **Milestone 15:** Named Ritual Combos & Effects
  - Each 3-spell ritual now has a unique generated name based on element composition
  - Rituals grant passive bonuses: Pure (3 same) = +18% (+20% for Arcane), Hybrid (2+1) = +12% dominant + +10% secondary, Triad (1/1/1) = +8% each
  - 🔥 Fire boosts damage, ❄️ Ice reduces cooldowns, ⚡ Thunder reduces mana cost, ✨ Arcane boosts mana generation
  - Including Spell Echo in a ritual adds a +5% kicker to all effects
  - Special "signature" combos have flavor names (Elemental Trinity, Apocalypse, Resonant Fire, etc.)
  - Live preview of ritual name and effects when selecting spells

#### v1.3.0 — Standalone & Cross-Platform
- [x] **Milestone 16:** Local JSON Storage & Windows Support
  - No database required! Game works out-of-the-box with local JSON saves
  - Saves stored in `~/.manatty/` (cross-platform)
  - Automatic fallback: tries MongoDB first, uses local storage if unavailable
  - Smart emoji detection: full emojis in Windows Terminal, ASCII fallback for legacy CMD
  - Pre-built binaries for Windows, macOS, Linux, FreeBSD (x64, ARM, 32-bit)
  - Path traversal protection with UUID validation
- [x] **v1.3.1 Balance & Polish Update**
  - **Sigil rebalance:** Damage requirements greatly increased (2.5x base from 200→500, exponent 1.6→1.8, floor factor 1.0→1.5) for meaningful damage gating
  - **Ritual redesign:** Arcane now grants mana generation (+20% pure bonus). Hybrid combos grant dual bonuses (+12% dominant + +10% secondary)
  - **Display fix:** Ritual effects now display in consistent order (no more flickering)

#### v1.4.0 — Ritual Synergies
- [x] **Milestone 18: Ritual Synergies & Chains** *(v1.4.0)*
  - Rituals interact when certain combinations are active together
  - Example synergies: Fire+Ice = "Thermal Shock" (+15% to both), Thunder+Arcane = "Mana Conduit" (+20% efficiency)
  - 6 synergy combinations: Thermal Shock, Mana Conduit, Volcanic Fury, Frozen Lightning, Arcane Inferno, Glacial Mystic
  - Transforms 3 independent ritual slots into interconnected puzzle
  - All bonuses are passive (idle-friendly)
  - Rewards strategic ritual composition planning
  - Performance optimized with synergy caching

#### v1.5.0 — Advanced Spell Rotation
- [x] **Milestone 19: Advanced Spell Rotation System** *(v1.5.0)*
  - Priority-based rotation planner (High/Medium/Low priority tiers)
  - Advanced conditions: "mana efficient", "during synergy", "sigil almost full"
  - Cooldown weaving optimization for maximum uptime
  - Mana reservation system (reserve % of mana pool)
  - Optimize for idle mode (spreads casts for sustained DPS vs burst)
  - Convert legacy auto-cast slots to rotation with one key press
  - Dramatically improves idle gameplay while rewarding optimization
  - Full UI with rotation view ([O] key) and real-time configuration

### Future Roadmap

#### v1.6.0 — Spell Combo System (Planned)

- [ ] **Milestone 20: Spell Combo System** *(v1.6.0)*
  - Casting specific spell sequences triggers combo bonuses
  - Examples: Fire→Fire→Ice = "Steam Burst" (+30% dmg), Ice→Thunder→Thunder = "Superconductor" (-50% cost)
  - 10-15 discoverable combos with visual feedback
  - Rewards thoughtful spell rotation planning
  - Works with auto-cast rotation system

#### Future Considerations

- [ ] **Spell Evolution & Mastery System**
  - Spells gain mastery points through usage (1 point per 1000 casts)
  - Permanent upgrades per spell: Fire +5% dmg, Ice -2% CD, Thunder -3% cost, Arcane +2% mana regen
  - Mastery persists through prestige, showing lifetime progress
  - Creates gradual permanent progression and rewards spell commitment

- [ ] **Prestige Artifact System**
  - Unlock rare artifacts at major milestones (floors 100, 250, 500, etc.)
  - Limited equipment slots (3-5 artifacts max)
  - Each artifact provides unique build-defining bonuses
  - Examples: "Mana Capacitor" (+50% max mana, -20% regen), "Temporal Loop" (every 10th spell triggers twice)
  - Forces meaningful equipment choices and build specialization

 
- [ ] **Mana Overflow Mechanics** - Convert excess mana generation into sigil charge, cooldown reduction, or burst damage
- [ ] **Ritual Maturation** - Rituals gain +1% effectiveness per hour active (caps at +50% after 50 hours)
- [ ] **Floor Event Memory** - Unlock permanent bonuses for consistently choosing same event type
- [ ] **Tower Milestone Auras** - Every 50 floors unlocks persistent aura (persists through prestige)
- [ ] **Ritual Presets** - Save and quickly recreate favorite ritual combinations
- [ ] **Enhanced Auto-Cast Conditions** - More sophisticated conditional triggers

## 🛠️ Tech Stack

- **Language:** Go 1.21+
- **TUI Framework:** [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **Styling:** [Lipgloss](https://github.com/charmbracelet/lipgloss)
- **Storage:** Local JSON files (default) or MongoDB (optional)

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

### Option 1: Download Pre-built Binary (Easiest)

Download the latest release for your platform from the `binaries/` folder:
- **Windows:** `manatty-windows-amd64.exe` (64-bit) or `manatty-windows-386.exe` (32-bit)
- **macOS:** `manatty-macos-arm64` (Apple Silicon) or `manatty-macos-amd64` (Intel)
- **Linux:** `manatty-linux-amd64`, `manatty-linux-arm64`, etc.

Just run it! No setup required - saves are stored locally in `~/.manatty/`.

### Option 2: Build from Source

#### Prerequisites
- Go 1.21 or later
- (Optional) MongoDB instance for cloud saves

#### Installation

```bash
# Clone the repository
git clone https://github.com/Ltorre/ManaTTY.git
cd ManaTTY

# Install dependencies
go mod download

# Run the game (no config needed!)
go run main.go

# Or build a binary
go build -o manatty .
./manatty
```

## ⚙️ Configuration (Optional)

The game works without any configuration! By default, it saves locally to `~/.manatty/`.

To use MongoDB instead, create a `.env` file in the project root:

```env
MONGODB_URI=mongodb://localhost:27017/mage_tower
LOG_LEVEL=info
GAME_TICK_RATE=10
AUTO_SAVE_INTERVAL=30
DEBUG=false
```

**Default values** (used when no `.env` exists):
- `LOG_LEVEL=info`
- `GAME_TICK_RATE=10`
- `AUTO_SAVE_INTERVAL=30`
- `DEBUG=false`
- Storage: Local JSON files

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
| `O` | Open Rotation view (v1.5.0) |
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

### Rotation View Controls (v1.5.0)

| Key | Action |
|-----|--------|
| `O` | Toggle rotation system on/off |
| `V` | Convert auto-cast slots to rotation |
| `W` | Toggle cooldown weaving |
| `Space` | Enable/disable selected spell |
| `↑/↓` | Navigate spell list |
| `Esc` | Return to tower |

## 📝 License

MIT License

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
