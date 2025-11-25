package state

func (s *State) String(key string) (string, bool) {
	return ExtractAs[string](key, s.data)
}

func (s *State) StringOr(key string, def string) string {
	if val, ok := ExtractAs[string](key, s.data); ok {
		return val
	}
	return def
}

func (s *State) Int(key string) (int, bool) {
	return ExtractAs[int](key, s.data)
}

func (s *State) IntOr(key string, def int) int {
	if val, ok := ExtractAs[int](key, s.data); ok {
		return val
	}
	return def
}

func (s *State) Float(key string) (float64, bool) {
	return ExtractAs[float64](key, s.data)
}

func (s *State) FloatOr(key string, def float64) float64 {
	if val, ok := ExtractAs[float64](key, s.data); ok {
		return val
	}
	return def
}

func (s *State) Bool(key string) (bool, bool) {
	return ExtractAs[bool](key, s.data)
}

func (s *State) BoolOr(key string, def bool) bool {
	if val, ok := ExtractAs[bool](key, s.data); ok {
		return val
	}
	return def
}

func (s *State) Map(key string) (*State, bool) {
	if value, exists := s.data[key]; exists {
		if nestedMap, ok := value.(map[string]any); ok {
			return &State{
				data: nestedMap,
				path: s.path + "." + key,
			}, true
		}
	}
	return nil, false
}

func (s *State) MapOr(key string) *State {
	if nested, ok := s.Map(key); ok {
		return nested
	}
	return &State{
		data: make(map[string]any),
		path: s.path + "." + key,
	}
}

func (s *State) Slice(key string) ([]any, bool) {
	if value, exists := s.data[key]; exists {
		if slice, ok := value.([]any); ok {
			return slice, true
		}
	}
	return nil, false
}

func SliceTyped[T any](key string, s *State) ([]T, bool) {
	if slice, ok := s.Slice(key); ok {
		result := make([]T, 0, len(slice))
		for _, item := range slice {
			if typed, is := item.(T); is {
				result = append(result, typed)
			} else if converted, success := convertValue[T](item); success {
				result = append(result, converted)
			}
		}
		return result, true
	}
	return nil, false
}

func MapTyped[K comparable, V any](key string, s *State) (map[K]V, bool) {
	if nested, ok := s.Map(key); ok {
		result := make(map[K]V)
		for k, v := range nested.data {
			var keyTyped K
			var valTyped V

			if kConverted, is := convertValue[K](k); is {
				keyTyped = kConverted
			} else {
				continue
			}

			if vConverted, is := convertValue[V](v); is {
				valTyped = vConverted
			} else {
				continue
			}

			result[keyTyped] = valTyped
		}
		return result, len(result) > 0
	}
	return make(map[K]V), false
}
