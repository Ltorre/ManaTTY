# Ritual Combo Effects (v1.2.0 Spec)

This document defines the **naming system** and **gameplay effects** for all 3-spell ritual combos.

With 12 spells, there are **C(12,3) = 220 possible combos**. Rather than hand-craft 220 entries, we use a **rule-based system** derived from spell elements and special spell identities.

---

## Spell Reference

| ID | Name | Element | Unlock Floor |
|----|------|---------|--------------|
| spell_fireball | Fireball | Fire | 1 |
| spell_inferno | Inferno | Fire | 25 |
| spell_meteor_strike | Meteor Strike | Fire | 75 (Prestige) |
| spell_frostbolt | Frostbolt | Ice | 3 |
| spell_blizzard | Blizzard | Ice | 35 |
| spell_frost_nova | Frost Nova | Ice | 60 |
| spell_lightning | Lightning | Thunder | 5 |
| spell_chain_lightning | Chain Lightning | Thunder | 45 |
| spell_thunderstorm | Thunderstorm | Thunder | 70 |
| spell_vortex | Arcane Vortex | Arcane | 10 |
| spell_echo | Spell Echo | Arcane | 55 |
| spell_arcane_blast | Arcane Blast | Arcane | 80 |

---

## Naming Rules

Ritual names are generated from two parts: **Element Adjective(s)** + **Power Noun**.

### Element Adjectives (by count)
- **Fire:** Blazing (1), Infernal (2), Volcanic (3)
- **Ice:** Frozen (1), Glacial (2), Permafrost (3)
- **Thunder:** Storm (1), Tempest (2), Cataclysm (3)
- **Arcane:** Runic (1), Ethereal (2), Astral (3)

### Power Noun (from highest-damage spell in the ritual)
| Spell | Noun |
|-------|------|
| Fireball | Ember |
| Inferno | Pyre |
| Meteor Strike | Cataclysm |
| Frostbolt | Shard |
| Blizzard | Gale |
| Frost Nova | Nova |
| Lightning | Bolt |
| Chain Lightning | Arc |
| Thunderstorm | Tempest |
| Arcane Vortex | Vortex |
| Spell Echo | Echo |
| Arcane Blast | Blast |

### Name Formula
```
"Ritual of [Adj1] [Adj2?] [Noun]"
```
- Single dominant element (2+ spells): use higher-tier adjective
- Mixed elements (1/1/1): list adjectives alphabetically

---

## Effect Rules

Effects are **passive bonuses** active while the ritual is equipped. Duration: always-on (no floor limit).

### Element-Based Effects

| Composition | Effect Type | Magnitude |
|-------------|-------------|-----------|
| **Pure (3 same element)** | Element's signature bonus | **+18%** |
| **Hybrid (2+1)** | Dominant element's bonus | **+12%** |
| **Triad (1/1/1)** | All three element bonuses | **+8% each** |

### Element Signature Bonuses
| Element | Bonus |
|---------|-------|
| Fire | +X% spell damage |
| Ice | -X% spell cooldown |
| Thunder | -X% mana cost |
| Arcane | +X% sigil charge rate |

### Spell Echo Kicker
If **Spell Echo** is one of the 3 spells, the ritual gains an additional **+5%** to all its effect magnitudes.

---

## Example Combos

### Pure Rituals (3 same element)

#### Ritual of Volcanic Cataclysm
- **Spells:** Fireball, Inferno, Meteor Strike
- **Effect:** Pure Fire: +18% spell damage

#### Ritual of Permafrost Nova
- **Spells:** Frostbolt, Blizzard, Frost Nova
- **Effect:** Pure Ice: -18% spell cooldown

#### Ritual of Cataclysm Tempest
- **Spells:** Lightning, Chain Lightning, Thunderstorm
- **Effect:** Pure Thunder: -18% mana cost

#### Ritual of Astral Blast
- **Spells:** Arcane Vortex, Spell Echo, Arcane Blast
- **Effect:** Pure Arcane: +18% sigil charge rate (+5% Echo kicker = +23%)

---

### Hybrid Rituals (2+1)

#### Ritual of Infernal Frozen Pyre
- **Spells:** Fireball, Inferno, Frostbolt
- **Effect:** Hybrid Fire-led: +12% spell damage

#### Ritual of Glacial Storm Nova
- **Spells:** Frostbolt, Frost Nova, Lightning
- **Effect:** Hybrid Ice-led: -12% spell cooldown

#### Ritual of Tempest Runic Arc
- **Spells:** Lightning, Chain Lightning, Arcane Vortex
- **Effect:** Hybrid Thunder-led: -12% mana cost

#### Ritual of Ethereal Blazing Echo
- **Spells:** Arcane Vortex, Spell Echo, Fireball
- **Effect:** Hybrid Arcane-led: +12% sigil charge rate (+5% Echo kicker = +17%)

---

### Triad Rituals (1/1/1)

#### Ritual of Blazing Frozen Storm Pyre
- **Spells:** Inferno, Blizzard, Chain Lightning
- **Effect:** Triad: +8% damage, -8% cooldown, -8% mana cost

#### Ritual of Blazing Runic Storm Echo
- **Spells:** Fireball, Spell Echo, Lightning
- **Effect:** Triad: +8% damage, +8% sigil charge, -8% mana cost (+5% Echo kicker = +13%/+13%/-13%)

---

## Notable "Signature" Combos

Some combos get special flavor names (displayed alongside the generated name):

| Spells | Signature Name | Special Note |
|--------|----------------|--------------|
| Spell Echo + any 2 same-element | "Resonant [Element]" | Echo amplifies element theme |
| Fireball + Lightning + Frostbolt | "Elemental Trinity" | Starter spells, first triad |
| Meteor Strike + Thunderstorm + Arcane Blast | "Apocalypse" | Three highest-damage spells |
| Inferno + Blizzard + Chain Lightning | "Convergence" | Mid-tier multi-element |

---

## Implementation Notes

1. **Ritual struct changes:** Add `EffectType` and `EffectMagnitude` fields
2. **Name generation:** Implement `generateRitualName(spellIDs)` using the rules above
3. **Effect application:** Hook into `CalculateManaPerSecond`, `CastSpell` (damage/cooldown/cost), and sigil charge logic
4. **UI:** Show ritual effect in Rituals view (e.g., "🔥 +12% damage")

---

## Balance Considerations

- Pure rituals reward mono-element builds (synergizes with Elemental Resonance)
- Triads are weaker per-stat but more versatile
- Echo kicker encourages using the iconic Spell Echo in ritual combos
- All effects are multiplicative with existing bonuses (prestige, synergy, etc.)
