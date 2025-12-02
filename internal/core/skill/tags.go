package skill

// =============================================================================
// SKILL TAGS
// =============================================================================

// Common skill tags for filtering and modification

// Element tags
const (
	TagFire      = "fire"
	TagCold      = "cold"
	TagLightning = "lightning"
	TagPoison    = "poison"
	TagPhysical  = "physical"
	TagChaos     = "chaos"
	TagHoly      = "holy"
	TagShadow    = "shadow"
	TagArcane    = "arcane"
)

// Combat style tags
const (
	TagMelee   = "melee"
	TagRanged  = "ranged"
	TagSpell   = "spell"
	TagAttack  = "attack"
	TagDefense = "defense"
	TagSupport = "support"
)

// Mechanic tags
const (
	TagAoE        = "aoe"
	TagProjectile = "projectile"
	TagDuration   = "duration"
	TagInstant    = "instant"
	TagChannel    = "channel"
	TagMovement   = "movement"
	TagSummon     = "summon"
)

// Effect tags
const (
	TagDamage  = "damage"
	TagHeal    = "heal"
	TagBuff    = "buff"
	TagDebuff  = "debuff"
	TagControl = "control"
	TagDispel  = "dispel"
)

// Category tags
const (
	TagCrafting  = "crafting"
	TagTrading   = "trading"
	TagGathering = "gathering"
	TagSurvival  = "survival"
	TagSocial    = "social"
)

// Rarity/power tags
const (
	TagBasic     = "basic"
	TagAdvanced  = "advanced"
	TagExpert    = "expert"
	TagMaster    = "master"
	TagLegendary = "legendary"
)

// =============================================================================
// TAG UTILITIES
// =============================================================================

// HasTag checks if tag list contains specific tag
func HasTag(tags []string, tag string) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}

// HasAnyTag checks if tag list contains any of the specified tags
func HasAnyTag(tags []string, check []string) bool {
	for _, c := range check {
		if HasTag(tags, c) {
			return true
		}
	}
	return false
}

// HasAllTags checks if tag list contains all specified tags
func HasAllTags(tags []string, check []string) bool {
	for _, c := range check {
		if !HasTag(tags, c) {
			return false
		}
	}
	return true
}

// FilterByTag returns tags that match a specific tag
func FilterByTag(skills []Def, tag string) []Def {
	var result []Def
	for _, s := range skills {
		if s.HasTag(tag) {
			result = append(result, s)
		}
	}
	return result
}

// GetElementTags returns all element-related tags from list
func GetElementTags(tags []string) []string {
	elements := []string{
		TagFire, TagCold, TagLightning, TagPoison,
		TagPhysical, TagChaos, TagHoly, TagShadow, TagArcane,
	}

	var result []string
	for _, e := range elements {
		if HasTag(tags, e) {
			result = append(result, e)
		}
	}
	return result
}

// IsDamageSkill checks if skill has damage-related tags
func IsDamageSkill(tags []string) bool {
	return HasTag(tags, TagDamage) || HasTag(tags, TagAttack)
}

// IsDefenseSkill checks if skill has defense-related tags
func IsDefenseSkill(tags []string) bool {
	return HasTag(tags, TagDefense) || HasTag(tags, TagBuff)
}

// IsSupportSkill checks if skill has support-related tags
func IsSupportSkill(tags []string) bool {
	return HasTag(tags, TagSupport) || HasTag(tags, TagHeal) || HasTag(tags, TagBuff)
}
