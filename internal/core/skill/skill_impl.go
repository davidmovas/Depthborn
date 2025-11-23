package skill

import (
	"context"
	"errors"
)

var _ Skill = (*Impl)(nil)

type Impl struct {
	id           string
	name         string
	description  string
	skillType    Type
	level        int
	maxLevel     int
	cooldown     int64
	maxCooldown  int64
	manaCost     float64
	tags         []string
	metadata     map[string]any
	requirements Requirements
	useFunc      func(ctx context.Context, casterID string, params ActivationParams) (Result, error)
}

func NewSkill(id, name string, description string, skillType Type, maxLevel int, manaCost float64, cooldown int64, tags []string, requirements Requirements, useFunc func(ctx context.Context, casterID string, params ActivationParams) (Result, error)) *Impl {
	return &Impl{
		id:           id,
		name:         name,
		description:  description,
		skillType:    skillType,
		level:        0,
		maxLevel:     maxLevel,
		manaCost:     manaCost,
		cooldown:     0,
		maxCooldown:  cooldown,
		tags:         tags,
		metadata:     make(map[string]any),
		requirements: requirements,
		useFunc:      useFunc,
	}
}

func (s *Impl) ID() string                 { return s.id }
func (s *Impl) Type() string               { return string(s.skillType) }
func (s *Impl) Name() string               { return s.name }
func (s *Impl) SetName(name string)        { s.name = name }
func (s *Impl) Description() string        { return s.description }
func (s *Impl) SkillType() Type            { return s.skillType }
func (s *Impl) Level() int                 { return s.level }
func (s *Impl) MaxLevel() int              { return s.maxLevel }
func (s *Impl) Cooldown() int64            { return s.cooldown }
func (s *Impl) MaxCooldown() int64         { return s.maxCooldown }
func (s *Impl) SetCooldown(ms int64)       { s.cooldown = ms }
func (s *Impl) ManaCost() float64          { return s.manaCost }
func (s *Impl) Tags() []string             { return s.tags }
func (s *Impl) Metadata() map[string]any   { return s.metadata }
func (s *Impl) Requirements() Requirements { return s.requirements }

func (s *Impl) CanLevelUp() bool {
	return s.level < s.maxLevel
}

func (s *Impl) LevelUp() error {
	if !s.CanLevelUp() {
		return errors.New("skill already at max level")
	}
	s.level++
	return nil
}

func (s *Impl) CanUse(ctx context.Context, casterID string) bool {
	if s.requirements != nil && !s.requirements.Check(casterID) {
		return false
	}
	return s.level > 0
}

func (s *Impl) Use(ctx context.Context, casterID string, params ActivationParams) (Result, error) {
	if !s.CanUse(ctx, casterID) {
		return Result{Success: false, Message: "cannot use skill"}, nil
	}
	if s.useFunc != nil {
		return s.useFunc(ctx, casterID, params)
	}
	return Result{Success: true}, nil
}
