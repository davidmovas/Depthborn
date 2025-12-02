package skill

import "github.com/davidmovas/Depthborn/internal/core/types"

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

// ElementTags is the list of all element tags
var ElementTags = []string{
	TagFire, TagCold, TagLightning, TagPoison,
	TagPhysical, TagChaos, TagHoly, TagShadow, TagArcane,
}

// FilterByTag returns skills that have specific tag
func FilterByTag(skills []Def, tag string) []Def {
	var result []Def
	for _, s := range skills {
		if s.Tags().Has(tag) {
			result = append(result, s)
		}
	}
	return result
}

// GetElementTags returns all element-related tags from TagSet
func GetElementTags(tags types.TagSet) []string {
	var result []string
	for _, e := range ElementTags {
		if tags.Has(e) {
			result = append(result, e)
		}
	}
	return result
}

// IsDamageSkill checks if skill has damage-related tags
func IsDamageSkill(tags types.TagSet) bool {
	return tags.Has(TagDamage) || tags.Has(TagAttack)
}

// IsDefenseSkill checks if skill has defense-related tags
func IsDefenseSkill(tags types.TagSet) bool {
	return tags.Has(TagDefense) || tags.Has(TagBuff)
}

// IsSupportSkill checks if skill has support-related tags
func IsSupportSkill(tags types.TagSet) bool {
	return tags.Has(TagSupport) || tags.Has(TagHeal) || tags.Has(TagBuff)
}
