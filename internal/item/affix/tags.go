package affix

// Tag represents affix tag for filtering and categorization
type Tag string

// Element tags - damage and resistance types
const (
	TagFire      Tag = "fire"
	TagCold      Tag = "cold"
	TagLightning Tag = "lightning"
	TagPoison    Tag = "poison"
	TagChaos     Tag = "chaos"
	TagPhysical  Tag = "physical"
)

// Magic tags - arcane damage types
const (
	TagArcane Tag = "arcane"
	TagHoly   Tag = "holy"
	TagShadow Tag = "shadow"
	TagVoid   Tag = "void"
)

// Combat category tags
const (
	TagAttack  Tag = "attack"  // Affects attacks
	TagSpell   Tag = "spell"   // Affects spells
	TagMelee   Tag = "melee"   // Melee-specific
	TagRanged  Tag = "ranged"  // Ranged-specific
	TagDefense Tag = "defense" // Defensive bonuses
)

// Resource tags
const (
	TagLife     Tag = "life"     // Health-related
	TagMana     Tag = "mana"     // Mana-related
	TagStamina  Tag = "stamina"  // Stamina-related
	TagResource Tag = "resource" // Any resource
)

// Attribute tags
const (
	TagStrength     Tag = "strength"
	TagDexterity    Tag = "dexterity"
	TagIntelligence Tag = "intelligence"
	TagVitality     Tag = "vitality"
	TagAttribute    Tag = "attribute" // Any primary attribute
)

// Effect tags
const (
	TagCritical   Tag = "critical"   // Crit-related
	TagSpeed      Tag = "speed"      // Speed-related (attack, move, cast)
	TagLeech      Tag = "leech"      // Life/mana steal
	TagRegen      Tag = "regen"      // Regeneration
	TagResistance Tag = "resistance" // Any resistance
	TagPenetration Tag = "penetration" // Resistance penetration
)

// Special effect tags
const (
	TagCurse    Tag = "curse"    // Curse effects
	TagBlessing Tag = "blessing" // Blessing effects
	TagAura     Tag = "aura"     // Aura effects
	TagOnHit    Tag = "on_hit"   // Proc on hit
	TagOnKill   Tag = "on_kill"  // Proc on kill
	TagOnCrit   Tag = "on_crit"  // Proc on critical
)

// Item-specific tags
const (
	TagWeapon    Tag = "weapon"    // Weapon-only affixes
	TagArmor     Tag = "armor"     // Armor-only affixes
	TagAccessory Tag = "accessory" // Accessory-only affixes
	TagJewelry   Tag = "jewelry"   // Rings/amulets
	TagShield    Tag = "shield"    // Shield-specific
)

// Rarity/quality tags
const (
	TagCommon    Tag = "common"    // Can appear on common items
	TagRare      Tag = "rare"      // Rare+ items only
	TagLegendary Tag = "legendary" // Legendary+ items only
	TagCrafted   Tag = "crafted"   // Crafting-only affixes
	TagCorrupted Tag = "corrupted" // Corruption-only affixes
)

// Utility tags
const (
	TagLoot       Tag = "loot"       // Loot quantity/quality
	TagExperience Tag = "experience" // Experience gain
	TagGold       Tag = "gold"       // Gold find
	TagMovement   Tag = "movement"   // Movement-related
)

// String returns tag as string
func (t Tag) String() string {
	return string(t)
}

// AllElementTags returns all elemental tags
func AllElementTags() []Tag {
	return []Tag{TagFire, TagCold, TagLightning, TagPoison, TagChaos, TagPhysical}
}

// AllMagicTags returns all magic type tags
func AllMagicTags() []Tag {
	return []Tag{TagArcane, TagHoly, TagShadow, TagVoid}
}

// AllDamageTags returns all damage-related tags
func AllDamageTags() []Tag {
	return []Tag{
		TagFire, TagCold, TagLightning, TagPoison, TagChaos, TagPhysical,
		TagArcane, TagHoly, TagShadow, TagVoid,
	}
}

// AllCombatTags returns all combat category tags
func AllCombatTags() []Tag {
	return []Tag{TagAttack, TagSpell, TagMelee, TagRanged, TagDefense}
}

// AllResourceTags returns all resource tags
func AllResourceTags() []Tag {
	return []Tag{TagLife, TagMana, TagStamina, TagResource}
}

// AllAttributeTags returns all attribute tags
func AllAttributeTags() []Tag {
	return []Tag{TagStrength, TagDexterity, TagIntelligence, TagVitality, TagAttribute}
}

// HasTag checks if tag slice contains specific tag
func HasTag(tags []Tag, target Tag) bool {
	for _, t := range tags {
		if t == target {
			return true
		}
	}
	return false
}

// HasAnyTag checks if tag slice contains any of target tags
func HasAnyTag(tags []Tag, targets []Tag) bool {
	for _, target := range targets {
		if HasTag(tags, target) {
			return true
		}
	}
	return false
}

// HasAllTags checks if tag slice contains all target tags
func HasAllTags(tags []Tag, targets []Tag) bool {
	for _, target := range targets {
		if !HasTag(tags, target) {
			return false
		}
	}
	return true
}

// StringsToTags converts string slice to tag slice
func StringsToTags(strings []string) []Tag {
	tags := make([]Tag, len(strings))
	for i, s := range strings {
		tags[i] = Tag(s)
	}
	return tags
}

// TagsToStrings converts tag slice to string slice
func TagsToStrings(tags []Tag) []string {
	strings := make([]string, len(tags))
	for i, t := range tags {
		strings[i] = string(t)
	}
	return strings
}
