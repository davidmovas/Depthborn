package affix

type BaseRequirements struct {
	minItemLevel int
	maxItemLevel int
	allowedTypes []string
	allowedSlots []string
}

func NewBaseRequirements(minItemLevel int) *BaseRequirements {
	return &BaseRequirements{
		minItemLevel: minItemLevel,
		maxItemLevel: 0, // No limit
		allowedTypes: make([]string, 0),
		allowedSlots: make([]string, 0),
	}
}

func (br *BaseRequirements) MinItemLevel() int {
	return br.minItemLevel
}

func (br *BaseRequirements) MaxItemLevel() int {
	return br.maxItemLevel
}

func (br *BaseRequirements) AllowedTypes() []string {
	return br.allowedTypes
}

func (br *BaseRequirements) AllowedSlots() []string {
	return br.allowedSlots
}

func (br *BaseRequirements) Check(itemType string, itemLevel int, slot string) bool {
	if itemLevel < br.minItemLevel {
		return false
	}
	if br.maxItemLevel > 0 && itemLevel > br.maxItemLevel {
		return false
	}

	if len(br.allowedTypes) > 0 {
		typeAllowed := false
		for _, allowedType := range br.allowedTypes {
			if allowedType == itemType {
				typeAllowed = true
				break
			}
		}
		if !typeAllowed {
			return false
		}
	}

	if len(br.allowedSlots) > 0 {
		slotAllowed := false
		for _, allowedSlot := range br.allowedSlots {
			if allowedSlot == slot {
				slotAllowed = true
				break
			}
		}
		if !slotAllowed {
			return false
		}
	}

	return true
}

func (br *BaseRequirements) SetMaxItemLevel(max int) {
	br.maxItemLevel = max
}

func (br *BaseRequirements) AddAllowedType(itemType string) {
	br.allowedTypes = append(br.allowedTypes, itemType)
}

func (br *BaseRequirements) AddAllowedSlot(slot string) {
	br.allowedSlots = append(br.allowedSlots, slot)
}
