package skill

import (
	"errors"
	"sync"
)

type LoadoutImpl struct {
	slots map[int]Skill
	mu    sync.RWMutex
	size  int
}

func NewLoadout(slotCount int) *LoadoutImpl {
	return &LoadoutImpl{
		slots: make(map[int]Skill),
		size:  slotCount,
	}
}

func (l *LoadoutImpl) Equip(slot int, skill Skill) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if slot >= l.size || slot < 0 {
		return errors.New("invalid slot")
	}
	l.slots[slot] = skill
	return nil
}

func (l *LoadoutImpl) Unequip(slot int) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, ok := l.slots[slot]; !ok {
		return errors.New("slot empty")
	}
	delete(l.slots, slot)
	return nil
}

func (l *LoadoutImpl) GetSkill(slot int) (Skill, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	skill, ok := l.slots[slot]
	return skill, ok
}

func (l *LoadoutImpl) GetAllSkills() map[int]Skill {
	l.mu.RLock()
	defer l.mu.RUnlock()
	c := make(map[int]Skill, len(l.slots))
	for k, v := range l.slots {
		c[k] = v
	}
	return c
}

func (l *LoadoutImpl) SlotCount() int {
	return l.size
}

func (l *LoadoutImpl) Swap(slot1, slot2 int) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	s1, ok1 := l.slots[slot1]
	s2, ok2 := l.slots[slot2]
	if !ok1 && !ok2 {
		return errors.New("both slots empty")
	}
	l.slots[slot1] = s2
	l.slots[slot2] = s1
	return nil
}

func (l *LoadoutImpl) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.slots = make(map[int]Skill)
}
