package account

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/davidmovas/Depthborn/internal/item"
	"github.com/davidmovas/Depthborn/pkg/identifier"
	"github.com/davidmovas/Depthborn/pkg/persist"
)

// Stash represents account-wide shared storage with tabs
// Unlike inventory, stash has NO weight limits - only slot limits per tab
type Stash struct {
	mu sync.RWMutex

	tabs    []*StashTab
	maxTabs int
}

// StashConfig holds configuration for creating a stash
type StashConfig struct {
	InitialTabs int
	MaxTabs     int
	SlotsPerTab int // Number of slots per tab
}

// DefaultStashConfig returns default configuration
func DefaultStashConfig() StashConfig {
	return StashConfig{
		InitialTabs: 1,
		MaxTabs:     10,
		SlotsPerTab: 60, // 6x10 grid
	}
}

// NewStash creates a new stash
func NewStash(cfg StashConfig) *Stash {
	if cfg.MaxTabs <= 0 {
		cfg.MaxTabs = 10
	}
	if cfg.SlotsPerTab <= 0 {
		cfg.SlotsPerTab = 60
	}
	if cfg.InitialTabs <= 0 {
		cfg.InitialTabs = 1
	}
	if cfg.InitialTabs > cfg.MaxTabs {
		cfg.InitialTabs = cfg.MaxTabs
	}

	s := &Stash{
		tabs:    make([]*StashTab, 0, cfg.MaxTabs),
		maxTabs: cfg.MaxTabs,
	}

	// Create initial tabs
	for i := 0; i < cfg.InitialTabs; i++ {
		tab := NewStashTab(fmt.Sprintf("Stash %d", i+1), cfg.SlotsPerTab)
		s.tabs = append(s.tabs, tab)
	}

	return s
}

// --- Tab Management ---

// Tabs returns all stash tabs
func (s *Stash) Tabs() []*StashTab {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*StashTab, len(s.tabs))
	copy(result, s.tabs)
	return result
}

// GetTab returns tab by index
func (s *Stash) GetTab(index int) (*StashTab, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if index < 0 || index >= len(s.tabs) {
		return nil, false
	}
	return s.tabs[index], true
}

// AddTab creates a new stash tab
func (s *Stash) AddTab(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.tabs) >= s.maxTabs {
		return fmt.Errorf("maximum number of stash tabs reached (%d)", s.maxTabs)
	}

	// Use default slots per tab from first tab, or default
	slotsPerTab := 60
	if len(s.tabs) > 0 {
		slotsPerTab = s.tabs[0].SlotCount()
	}

	tab := NewStashTab(name, slotsPerTab)
	s.tabs = append(s.tabs, tab)
	return nil
}

// AddTabWithSlots creates a new stash tab with specific slot count
func (s *Stash) AddTabWithSlots(name string, slots int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.tabs) >= s.maxTabs {
		return fmt.Errorf("maximum number of stash tabs reached (%d)", s.maxTabs)
	}

	tab := NewStashTab(name, slots)
	s.tabs = append(s.tabs, tab)
	return nil
}

// RemoveTab removes stash tab by index
func (s *Stash) RemoveTab(index int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if index < 0 || index >= len(s.tabs) {
		return fmt.Errorf("tab index out of range: %d", index)
	}

	if len(s.tabs) == 1 {
		return fmt.Errorf("cannot remove the last stash tab")
	}

	// Check if tab is empty
	if s.tabs[index].ItemCount() > 0 {
		return fmt.Errorf("cannot remove non-empty stash tab")
	}

	s.tabs = append(s.tabs[:index], s.tabs[index+1:]...)
	return nil
}

// RenameTab updates tab name
func (s *Stash) RenameTab(index int, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if index < 0 || index >= len(s.tabs) {
		return fmt.Errorf("tab index out of range: %d", index)
	}

	s.tabs[index].SetName(name)
	return nil
}

// SwapTabs swaps positions of two tabs
func (s *Stash) SwapTabs(index1, index2 int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if index1 < 0 || index1 >= len(s.tabs) || index2 < 0 || index2 >= len(s.tabs) {
		return fmt.Errorf("tab index out of range")
	}

	s.tabs[index1], s.tabs[index2] = s.tabs[index2], s.tabs[index1]
	return nil
}

// TabCount returns number of tabs
func (s *Stash) TabCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.tabs)
}

// MaxTabs returns maximum number of tabs
func (s *Stash) MaxTabs() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.maxTabs
}

// --- Item Operations ---

// TransferToTab moves item to specified tab
func (s *Stash) TransferToTab(ctx context.Context, itm item.Item, tabIndex int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if tabIndex < 0 || tabIndex >= len(s.tabs) {
		return fmt.Errorf("tab index out of range: %d", tabIndex)
	}

	// Find and remove item from its current tab
	var sourceTab *StashTab
	for _, tab := range s.tabs {
		if tab.Contains(itm.ID()) {
			sourceTab = tab
			break
		}
	}

	if sourceTab != nil {
		if _, err := sourceTab.Remove(ctx, itm.ID()); err != nil {
			return fmt.Errorf("failed to remove from source tab: %w", err)
		}
	}

	// Add to destination tab
	if err := s.tabs[tabIndex].Add(ctx, itm); err != nil {
		// Rollback if source tab exists
		if sourceTab != nil {
			_ = sourceTab.Add(ctx, itm)
		}
		return fmt.Errorf("failed to add to destination tab: %w", err)
	}

	return nil
}

// TransferToSlot moves item to specific slot in specified tab
func (s *Stash) TransferToSlot(ctx context.Context, itm item.Item, tabIndex, slot int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if tabIndex < 0 || tabIndex >= len(s.tabs) {
		return fmt.Errorf("tab index out of range: %d", tabIndex)
	}

	// Find and remove item from its current tab
	var sourceTab *StashTab
	for _, tab := range s.tabs {
		if tab.Contains(itm.ID()) {
			sourceTab = tab
			break
		}
	}

	if sourceTab != nil {
		if _, err := sourceTab.Remove(ctx, itm.ID()); err != nil {
			return fmt.Errorf("failed to remove from source tab: %w", err)
		}
	}

	// Add to destination slot
	if err := s.tabs[tabIndex].AddToSlot(ctx, slot, itm); err != nil {
		// Rollback if source tab exists
		if sourceTab != nil {
			_ = sourceTab.Add(ctx, itm)
		}
		return fmt.Errorf("failed to add to destination slot: %w", err)
	}

	return nil
}

// FindItem searches all tabs for item
func (s *Stash) FindItem(itemID string) (item.Item, int, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i, tab := range s.tabs {
		if itm, ok := tab.Get(itemID); ok {
			return itm, i, true
		}
	}
	return nil, -1, false
}

// --- Search & Filter (across all tabs) ---

// Search finds items matching query string across all tabs
func (s *Stash) Search(query string) []item.Item {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []item.Item
	for _, tab := range s.tabs {
		result = append(result, tab.Search(query)...)
	}
	return result
}

// FindByType returns items of given type across all tabs
func (s *Stash) FindByType(itemType item.Type) []item.Item {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []item.Item
	for _, tab := range s.tabs {
		result = append(result, tab.FindByType(itemType)...)
	}
	return result
}

// FindByRarity returns items of given rarity across all tabs
func (s *Stash) FindByRarity(rarity item.Rarity) []item.Item {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []item.Item
	for _, tab := range s.tabs {
		result = append(result, tab.FindByRarity(rarity)...)
	}
	return result
}

// FindByTag returns items with given tag across all tabs
func (s *Stash) FindByTag(tag string) []item.Item {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []item.Item
	for _, tab := range s.tabs {
		result = append(result, tab.FindByTag(tag)...)
	}
	return result
}

// Filter returns items matching predicate across all tabs
func (s *Stash) Filter(predicate func(item.Item) bool) []item.Item {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []item.Item
	for _, tab := range s.tabs {
		result = append(result, tab.Filter(predicate)...)
	}
	return result
}

// --- Stats ---

// TotalSlots returns combined slot capacity of all tabs
func (s *Stash) TotalSlots() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var total int
	for _, tab := range s.tabs {
		total += tab.SlotCount()
	}
	return total
}

// TotalUsedSlots returns total used slots across all tabs
func (s *Stash) TotalUsedSlots() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var total int
	for _, tab := range s.tabs {
		total += tab.UsedSlots()
	}
	return total
}

// TotalFreeSlots returns total free slots across all tabs
func (s *Stash) TotalFreeSlots() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var total int
	for _, tab := range s.tabs {
		total += tab.FreeSlots()
	}
	return total
}

// TotalCount returns total items across all tabs (unique stacks)
func (s *Stash) TotalCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var total int
	for _, tab := range s.tabs {
		total += tab.ItemCount()
	}
	return total
}

// TotalItems returns total item count including stack sizes
func (s *Stash) TotalItems() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var total int
	for _, tab := range s.tabs {
		total += tab.TotalItems()
	}
	return total
}

// TotalValue returns combined value of all items
func (s *Stash) TotalValue() int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var total int64
	for _, tab := range s.tabs {
		total += tab.TotalValue()
	}
	return total
}

// --- Persistence ---

// StashState holds serializable stash state
type StashState struct {
	MaxTabs int             `msgpack:"max_tabs"`
	Tabs    []StashTabState `msgpack:"tabs"`
}

func (s *Stash) SerializeState() (map[string]any, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tabs := make([]StashTabState, len(s.tabs))
	for i, tab := range s.tabs {
		tabs[i] = tab.ToState()
	}

	state := StashState{
		MaxTabs: s.maxTabs,
		Tabs:    tabs,
	}

	data, err := persist.DefaultCodec().Encode(state)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := persist.DefaultCodec().Decode(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Stash) DeserializeState(stateData map[string]any) error {
	data, err := persist.DefaultCodec().Encode(stateData)
	if err != nil {
		return err
	}

	var state StashState
	if err := persist.DefaultCodec().Decode(data, &state); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.maxTabs = state.MaxTabs

	s.tabs = make([]*StashTab, len(state.Tabs))
	for i, tabState := range state.Tabs {
		s.tabs[i] = StashTabFromState(tabState)
	}

	return nil
}

// =============================================================================
// StashTab - Single tab in stash (slot-based, no weight limit)
// =============================================================================

// StashTab represents a single stash tab with slot-based storage
type StashTab struct {
	mu sync.RWMutex

	name      string
	icon      string
	color     string
	slots     []item.Item    // slot index -> item (nil = empty)
	itemIndex map[string]int // itemID -> slot index
}

// StashTabState holds serializable tab state
type StashTabState struct {
	Name    string   `msgpack:"name"`
	Icon    string   `msgpack:"icon"`
	Color   string   `msgpack:"color"`
	Slots   int      `msgpack:"slots"`
	ItemIDs []string `msgpack:"item_ids,omitempty"`
}

// NewStashTab creates a new stash tab
func NewStashTab(name string, slotCount int) *StashTab {
	if slotCount <= 0 {
		slotCount = 60
	}
	return &StashTab{
		name:      name,
		icon:      "default",
		color:     "#ffffff",
		slots:     make([]item.Item, slotCount),
		itemIndex: make(map[string]int),
	}
}

// StashTabFromState creates a tab from serialized state
func StashTabFromState(state StashTabState) *StashTab {
	tab := NewStashTab(state.Name, state.Slots)
	tab.icon = state.Icon
	tab.color = state.Color
	// Items need to be restored separately by repository
	return tab
}

// --- Metadata ---

func (t *StashTab) Name() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.name
}

func (t *StashTab) SetName(name string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.name = name
}

func (t *StashTab) Icon() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.icon
}

func (t *StashTab) SetIcon(icon string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.icon = icon
}

func (t *StashTab) Color() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.color
}

func (t *StashTab) SetColor(color string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.color = color
}

// --- Basic Operations ---

// Add adds item to first available slot (auto-stacks if possible)
func (t *StashTab) Add(ctx context.Context, itm item.Item) error {
	if itm == nil {
		return fmt.Errorf("cannot add nil item")
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	// Try to stack with existing item first
	if targetID, canStack := t.canStackWithLocked(itm); canStack {
		return t.mergeIntoExistingLocked(itm, targetID)
	}

	// Find free slot
	slot := t.findFreeSlotLocked()
	if slot == -1 {
		return fmt.Errorf("stash tab is full (no free slots)")
	}

	t.slots[slot] = itm
	t.itemIndex[itm.ID()] = slot
	return nil
}

// AddToSlot adds item to specific slot
func (t *StashTab) AddToSlot(ctx context.Context, slot int, itm item.Item) error {
	if itm == nil {
		return fmt.Errorf("cannot add nil item")
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	if slot < 0 || slot >= len(t.slots) {
		return fmt.Errorf("slot %d out of range (0-%d)", slot, len(t.slots)-1)
	}

	if t.slots[slot] != nil {
		return fmt.Errorf("slot %d is already occupied", slot)
	}

	t.slots[slot] = itm
	t.itemIndex[itm.ID()] = slot
	return nil
}

// Remove removes item by ID completely
func (t *StashTab) Remove(ctx context.Context, itemID string) (item.Item, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	slot, exists := t.itemIndex[itemID]
	if !exists {
		return nil, fmt.Errorf("item with ID %s not found", itemID)
	}

	itm := t.slots[slot]
	t.slots[slot] = nil
	delete(t.itemIndex, itemID)

	return itm, nil
}

// RemoveAmount removes specific amount from a stack
func (t *StashTab) RemoveAmount(ctx context.Context, itemID string, amount int) (item.Item, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	slot, exists := t.itemIndex[itemID]
	if !exists {
		return nil, fmt.Errorf("item with ID %s not found", itemID)
	}

	itm := t.slots[slot]
	currentStack := itm.StackSize()

	if amount >= currentStack {
		// Remove entire item
		t.slots[slot] = nil
		delete(t.itemIndex, itemID)
		return itm, nil
	}

	// Reduce stack size
	itm.RemoveStack(amount)

	// Create new item for removed portion (clone with new ID)
	removed := itm.Clone().(item.Item)
	if setter, ok := removed.(interface{ SetID(string) }); ok {
		setter.SetID(identifier.New())
	}
	// Reset clone's stack to the removed amount
	removed.RemoveStack(removed.StackSize() - 1) // Reset to 1
	removed.AddStack(amount - 1)                 // Set to amount

	return removed, nil
}

// Get returns item by ID
func (t *StashTab) Get(itemID string) (item.Item, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	slot, exists := t.itemIndex[itemID]
	if !exists {
		return nil, false
	}
	return t.slots[slot], true
}

// GetAtSlot returns item at specific slot
func (t *StashTab) GetAtSlot(slot int) (item.Item, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if slot < 0 || slot >= len(t.slots) {
		return nil, false
	}
	itm := t.slots[slot]
	return itm, itm != nil
}

// GetAll returns all items in tab
func (t *StashTab) GetAll() []item.Item {
	t.mu.RLock()
	defer t.mu.RUnlock()

	result := make([]item.Item, 0, len(t.itemIndex))
	for _, itm := range t.slots {
		if itm != nil {
			result = append(result, itm)
		}
	}
	return result
}

// Contains checks if item exists in tab
func (t *StashTab) Contains(itemID string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	_, exists := t.itemIndex[itemID]
	return exists
}

// Clear removes all items from tab
func (t *StashTab) Clear(ctx context.Context) []item.Item {
	t.mu.Lock()
	defer t.mu.Unlock()

	items := make([]item.Item, 0, len(t.itemIndex))
	for _, itm := range t.slots {
		if itm != nil {
			items = append(items, itm)
		}
	}

	t.slots = make([]item.Item, len(t.slots))
	t.itemIndex = make(map[string]int)

	return items
}

// --- Stack Operations ---

// SplitStack splits a stack into two, returns the new stack
func (t *StashTab) SplitStack(ctx context.Context, itemID string, amount int) (item.Item, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	slot, exists := t.itemIndex[itemID]
	if !exists {
		return nil, fmt.Errorf("item with ID %s not found", itemID)
	}

	itm := t.slots[slot]
	if itm.StackSize() <= amount {
		return nil, fmt.Errorf("cannot split: stack size %d is not greater than %d", itm.StackSize(), amount)
	}

	// Find free slot for new stack
	newSlot := t.findFreeSlotLocked()
	if newSlot == -1 {
		return nil, fmt.Errorf("no free slot for split stack")
	}

	// Remove from original stack
	itm.RemoveStack(amount)

	// Create new item (clone and set stack)
	newItem := itm.Clone().(item.Item)
	newItem.RemoveStack(newItem.StackSize() - 1) // Reset to 1
	newItem.AddStack(amount - 1)                 // Set to amount

	// Generate new ID for split item
	if setter, ok := newItem.(interface{ SetID(string) }); ok {
		setter.SetID(identifier.New())
	}

	t.slots[newSlot] = newItem
	t.itemIndex[newItem.ID()] = newSlot

	return newItem, nil
}

// MergeStacks merges source stack into target stack
func (t *StashTab) MergeStacks(ctx context.Context, sourceID, targetID string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	sourceSlot, sourceExists := t.itemIndex[sourceID]
	targetSlot, targetExists := t.itemIndex[targetID]

	if !sourceExists {
		return fmt.Errorf("source item %s not found", sourceID)
	}
	if !targetExists {
		return fmt.Errorf("target item %s not found", targetID)
	}

	source := t.slots[sourceSlot]
	target := t.slots[targetSlot]

	if !target.CanStackWith(source) {
		return fmt.Errorf("items cannot be stacked together")
	}

	availableSpace := target.MaxStackSize() - target.StackSize()
	if availableSpace <= 0 {
		return fmt.Errorf("target stack is full")
	}

	amountToMove := source.StackSize()
	if amountToMove > availableSpace {
		amountToMove = availableSpace
	}

	target.AddStack(amountToMove)
	source.RemoveStack(amountToMove)

	if source.StackSize() <= 0 {
		t.slots[sourceSlot] = nil
		delete(t.itemIndex, sourceID)
	}

	return nil
}

// CanStackWith checks if item can stack with existing items
func (t *StashTab) CanStackWith(itm item.Item) (string, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.canStackWithLocked(itm)
}

func (t *StashTab) canStackWithLocked(itm item.Item) (string, bool) {
	for _, existing := range t.slots {
		if existing != nil && existing.CanStackWith(itm) {
			if existing.StackSize() < existing.MaxStackSize() {
				return existing.ID(), true
			}
		}
	}
	return "", false
}

func (t *StashTab) mergeIntoExistingLocked(itm item.Item, targetID string) error {
	targetSlot := t.itemIndex[targetID]
	target := t.slots[targetSlot]

	availableSpace := target.MaxStackSize() - target.StackSize()
	amountToAdd := itm.StackSize()

	if amountToAdd <= availableSpace {
		// All fits in existing stack
		target.AddStack(amountToAdd)
		return nil
	}

	// Partial stack - add what fits, then add remainder as new item
	target.AddStack(availableSpace)

	// Remainder needs new slot
	slot := t.findFreeSlotLocked()
	if slot == -1 {
		return fmt.Errorf("stash tab is full")
	}

	itm.RemoveStack(availableSpace)
	t.slots[slot] = itm
	t.itemIndex[itm.ID()] = slot
	return nil
}

// --- Slot Management ---

// SlotCount returns number of slots
func (t *StashTab) SlotCount() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.slots)
}

// SetSlotCount changes number of slots
func (t *StashTab) SetSlotCount(count int) {
	if count <= 0 {
		return
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	currentCount := len(t.slots)
	if count == currentCount {
		return
	}

	if count > currentCount {
		// Expand
		newSlots := make([]item.Item, count)
		copy(newSlots, t.slots)
		t.slots = newSlots
	} else {
		// Shrink - only if trailing slots are empty
		canShrink := true
		for i := count; i < currentCount; i++ {
			if t.slots[i] != nil {
				canShrink = false
				break
			}
		}
		if canShrink {
			t.slots = t.slots[:count]
		}
	}
}

// UsedSlots returns number of occupied slots
func (t *StashTab) UsedSlots() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.itemIndex)
}

// FreeSlots returns number of available slots
func (t *StashTab) FreeSlots() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.slots) - len(t.itemIndex)
}

// IsFull returns true if no more items can be added
func (t *StashTab) IsFull() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.findFreeSlotLocked() == -1
}

// SwapSlots swaps items between two slots
func (t *StashTab) SwapSlots(ctx context.Context, slot1, slot2 int) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	maxSlot := len(t.slots)
	if slot1 < 0 || slot1 >= maxSlot || slot2 < 0 || slot2 >= maxSlot {
		return fmt.Errorf("slot out of range")
	}

	item1 := t.slots[slot1]
	item2 := t.slots[slot2]

	t.slots[slot1] = item2
	t.slots[slot2] = item1

	if item1 != nil {
		t.itemIndex[item1.ID()] = slot2
	}
	if item2 != nil {
		t.itemIndex[item2.ID()] = slot1
	}

	return nil
}

// MoveToSlot moves item to a different slot
func (t *StashTab) MoveToSlot(ctx context.Context, itemID string, targetSlot int) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	maxSlot := len(t.slots)
	if targetSlot < 0 || targetSlot >= maxSlot {
		return fmt.Errorf("target slot %d out of range", targetSlot)
	}

	currentSlot, exists := t.itemIndex[itemID]
	if !exists {
		return fmt.Errorf("item %s not found", itemID)
	}

	if currentSlot == targetSlot {
		return nil
	}

	if t.slots[targetSlot] != nil {
		return fmt.Errorf("target slot %d is occupied", targetSlot)
	}

	itm := t.slots[currentSlot]
	t.slots[currentSlot] = nil
	t.slots[targetSlot] = itm
	t.itemIndex[itemID] = targetSlot

	return nil
}

func (t *StashTab) findFreeSlotLocked() int {
	for i, itm := range t.slots {
		if itm == nil {
			return i
		}
	}
	return -1
}

// --- Search & Filter ---

// Search finds items matching query string (name contains)
func (t *StashTab) Search(query string) []item.Item {
	query = strings.ToLower(query)
	return t.Filter(func(itm item.Item) bool {
		return strings.Contains(strings.ToLower(itm.Name()), query)
	})
}

// FindByType returns items of given type
func (t *StashTab) FindByType(itemType item.Type) []item.Item {
	return t.Filter(func(itm item.Item) bool {
		return itm.ItemType() == itemType
	})
}

// FindByRarity returns items of given rarity
func (t *StashTab) FindByRarity(rarity item.Rarity) []item.Item {
	return t.Filter(func(itm item.Item) bool {
		return itm.Rarity() == rarity
	})
}

// FindByTag returns items with given tag
func (t *StashTab) FindByTag(tag string) []item.Item {
	return t.Filter(func(itm item.Item) bool {
		return itm.Tags().Has(tag)
	})
}

// FindByTags returns items having ALL specified tags
func (t *StashTab) FindByTags(tags ...string) []item.Item {
	return t.Filter(func(itm item.Item) bool {
		return itm.Tags().Contains(tags...)
	})
}

// FindByAnyTag returns items having ANY of specified tags
func (t *StashTab) FindByAnyTag(tags ...string) []item.Item {
	return t.Filter(func(itm item.Item) bool {
		return itm.Tags().ContainsAny(tags...)
	})
}

// FindStackable returns all stackable items
func (t *StashTab) FindStackable() []item.Item {
	return t.Filter(func(itm item.Item) bool {
		return itm.MaxStackSize() > 1
	})
}

// Filter returns items matching predicate
func (t *StashTab) Filter(predicate func(item.Item) bool) []item.Item {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var result []item.Item
	for _, itm := range t.slots {
		if itm != nil && predicate(itm) {
			result = append(result, itm)
		}
	}
	return result
}

// --- Stats ---

// ItemCount returns number of unique items (stacks count as 1)
func (t *StashTab) ItemCount() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.itemIndex)
}

// TotalItems returns total item count including stack sizes
func (t *StashTab) TotalItems() int {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var total int
	for _, itm := range t.slots {
		if itm != nil {
			total += itm.StackSize()
		}
	}
	return total
}

// TotalValue returns combined value of all items
func (t *StashTab) TotalValue() int64 {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var total int64
	for _, itm := range t.slots {
		if itm != nil {
			total += itm.Value() * int64(itm.StackSize())
		}
	}
	return total
}

// --- Persistence ---

func (t *StashTab) ToState() StashTabState {
	t.mu.RLock()
	defer t.mu.RUnlock()

	itemIDs := make([]string, len(t.slots))
	for i, itm := range t.slots {
		if itm != nil {
			itemIDs[i] = itm.ID()
		}
	}

	return StashTabState{
		Name:    t.name,
		Icon:    t.icon,
		Color:   t.color,
		Slots:   len(t.slots),
		ItemIDs: itemIDs,
	}
}

// GetItemIDs returns all item IDs in slot order
func (t *StashTab) GetItemIDs() []string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	ids := make([]string, len(t.slots))
	for i, itm := range t.slots {
		if itm != nil {
			ids[i] = itm.ID()
		}
	}
	return ids
}

// AddDirect adds an item without stacking (for deserialization)
func (t *StashTab) AddDirect(itm item.Item) error {
	if itm == nil {
		return fmt.Errorf("cannot add nil item")
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	slot := t.findFreeSlotLocked()
	if slot == -1 {
		return fmt.Errorf("no free slot")
	}

	t.slots[slot] = itm
	t.itemIndex[itm.ID()] = slot
	return nil
}

// AddDirectToSlot adds an item to specific slot without stacking
func (t *StashTab) AddDirectToSlot(slot int, itm item.Item) error {
	if itm == nil {
		return fmt.Errorf("cannot add nil item")
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	if slot < 0 || slot >= len(t.slots) {
		return fmt.Errorf("slot out of range")
	}

	t.slots[slot] = itm
	t.itemIndex[itm.ID()] = slot
	return nil
}

// CanAdd checks if item can be added (slot check or can stack)
func (t *StashTab) CanAdd(itm item.Item) bool {
	if itm == nil {
		return false
	}

	t.mu.RLock()
	defer t.mu.RUnlock()

	// Can stack?
	if _, canStack := t.canStackWithLocked(itm); canStack {
		return true
	}

	// Free slot?
	return t.findFreeSlotLocked() != -1
}
