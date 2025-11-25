package state

import (
	"fmt"

	"github.com/davidmovas/Depthborn/internal/infra"
)

func (s *State) SetEntity(key string, entity infra.Serializable) error {
	if entity == nil {
		s.data[key] = nil
		return nil
	}

	stateData, err := entity.SerializeState()
	if err != nil {
		return fmt.Errorf("failed to serialize entity %s: %w", key, err)
	}

	s.data[key] = stateData
	return nil
}

func (s *State) GetEntity(key string, entity infra.Serializable) error {
	stateData, exists := s.data[key]
	if !exists {
		return fmt.Errorf("entity %s not found", key)
	}

	if stateMap, ok := stateData.(map[string]any); ok {
		return entity.DeserializeState(stateMap)
	}

	if stateWrapper, ok := stateData.(*State); ok {
		return entity.DeserializeState(stateWrapper.data)
	}

	return fmt.Errorf("invalid state data for entity %s", key)
}

func (s *State) ExtractEntity(key string, entity infra.Serializable) bool {
	if err := s.GetEntity(key, entity); err != nil {
		return false
	}
	return true
}

func (s *State) EntityState(key string) (*State, bool) {
	if value, exists := s.data[key]; exists {
		switch v := value.(type) {
		case map[string]any:
			return &State{
				data: v,
				path: s.path + "." + key,
			}, true
		case *State:
			return v, true
		}
	}
	return nil, false
}

func (s *State) EntityStateOr(key string) *State {
	if nested, ok := s.EntityState(key); ok {
		return nested
	}
	return &State{
		data: make(map[string]any),
		path: s.path + "." + key,
	}
}
