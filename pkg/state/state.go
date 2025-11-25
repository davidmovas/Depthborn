package state

import (
	"fmt"
	"log"
	"reflect"
	"strconv"

	"github.com/davidmovas/Depthborn/internal/infra"
)

type State struct {
	data map[string]any
	path string
}

type Pair struct {
	Key   string
	Value any
}

func New() *State {
	return &State{data: make(map[string]any), path: "root"}
}

func From(data map[string]any) *State {
	return &State{data: data, path: "root"}
}

func FromEntity(entity infra.Serializable) (*State, error) {
	data, err := entity.SerializeState()
	if err != nil {
		return nil, err
	}
	return From(data), nil
}

func (s *State) Data() map[string]any {
	result := make(map[string]any, len(s.data))
	for k, v := range s.data {
		result[k] = v
	}
	return result
}

func (s *State) RawData() map[string]any {
	return s.data
}

func (s *State) Has(key string) bool {
	_, exists := s.data[key]
	return exists
}

func (s *State) Len() int {
	return len(s.data)
}

func (s *State) Keys() []string {
	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	return keys
}

func (s *State) Get(key string) (any, bool) {
	value, exists := s.data[key]
	return value, exists
}

func (s *State) MustGet(key string) any {
	if value, exists := s.data[key]; exists {
		return value
	}
	panic(fmt.Sprintf("key %s.%s not found", s.path, key))
}

func (s *State) Extract(key string, target any) bool {
	value, exists := s.data[key]
	if !exists {
		return false
	}

	targetVal := reflect.ValueOf(target)
	if targetVal.Kind() != reflect.Ptr {
		return false
	}

	elem := targetVal.Elem()
	valueVal := reflect.ValueOf(value)

	if valueVal.Type().ConvertibleTo(elem.Type()) {
		elem.Set(valueVal.Convert(elem.Type()))
		return true
	}

	return s.trySpecialConversion(value, elem)
}

func (s *State) Try(key string, defaultValue any) any {
	if value, exists := s.data[key]; exists {
		return value
	}
	return defaultValue
}

func (s *State) Set(key string, value any) *State {
	s.data[key] = value
	return s
}

func (s *State) Batch(pair Pair, pairs ...Pair) *State {
	s.Set(pair.Key, pair.Value)
	for _, p := range pairs {
		s.Set(p.Key, p.Value)
	}
	return s
}

func BatchKV(kv ...any) (map[string]any, error) {
	state := New()

	if len(kv)%2 != 0 {
		return nil, fmt.Errorf("odd number of arguments: expected key/value pairs")
	}

	for i := 0; i < len(kv); i += 2 {
		key, ok := kv[i].(string)
		if !ok {
			return nil, fmt.Errorf("key at position %d is not string", i)
		}
		state.Set(key, kv[i+1])
	}

	return state.Data(), nil
}

func (s *State) SetAll(values map[string]any) *State {
	for k, v := range values {
		s.data[k] = v
	}
	return s
}

func (s *State) Delete(key string) *State {
	delete(s.data, key)
	return s
}

func (s *State) Merge(other *State) *State {
	for k, v := range other.data {
		s.data[k] = v
	}
	return s
}

func (s *State) Clone() *State {
	newData := make(map[string]any, len(s.data))
	for k, v := range s.data {
		newData[k] = v
	}
	return &State{data: newData, path: s.path}
}

func Set(key string, value any) *State {
	return New().Set(key, value)
}

func Batch(pair Pair, pairs ...Pair) *State {
	return New().Batch(pair, pairs...)
}

func Assign(s *State, targets map[string]any) error {
	for key, targetPtr := range targets {
		if !s.Extract(key, targetPtr) {
			log.Printf("Warning: failed to extract %s\n", key)
		}
	}
	return nil
}

func AssignStrict(s *State, targets map[string]any) error {
	for key, targetPtr := range targets {
		if !s.Extract(key, targetPtr) {
			return fmt.Errorf("failed to extract required field %s", key)
		}
	}
	return nil
}

func ExtractAs[T any](key string, data map[string]any) (T, bool) {
	var zero T
	value, exists := data[key]
	if !exists {
		return zero, false
	}

	if typed, ok := value.(T); ok {
		return typed, true
	}

	return convertValue[T](value)
}

func convertValue[T any](value any) (T, bool) {
	var zero T

	fromVal := reflect.ValueOf(value)
	toType := reflect.TypeOf(zero)

	if fromVal.Type().ConvertibleTo(toType) {
		return fromVal.Convert(toType).Interface().(T), true
	}

	switch any(zero).(type) {
	case int:
		switch v := value.(type) {
		case float64:
			return any(int(v)).(T), true
		case float32:
			return any(int(v)).(T), true
		case int64:
			return any(int(v)).(T), true
		case string:
			if i, err := strconv.Atoi(v); err == nil {
				return any(i).(T), true
			}
		case bool:
			if v {
				return any(1).(T), true
			}
			return any(0).(T), true
		}
	case float64:
		switch v := value.(type) {
		case int:
			return any(float64(v)).(T), true
		case int64:
			return any(float64(v)).(T), true
		case float32:
			return any(float64(v)).(T), true
		case string:
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				return any(f).(T), true
			}
		}
	case string:
		return any(fmt.Sprintf("%v", value)).(T), true
	case bool:
		switch v := value.(type) {
		case int:
			return any(v != 0).(T), true
		case float64:
			return any(v != 0).(T), true
		case string:
			if b, err := strconv.ParseBool(v); err == nil {
				return any(b).(T), true
			}
		}
	}

	return zero, false
}

func (s *State) trySpecialConversion(value any, target reflect.Value) bool {
	switch target.Kind() {
	case reflect.Int, reflect.Int64:
		switch v := value.(type) {
		case float64:
			target.SetInt(int64(v))
			return true
		case string:
			if i, err := strconv.ParseInt(v, 10, 64); err == nil {
				target.SetInt(i)
				return true
			}
		}
	case reflect.Float64:
		switch v := value.(type) {
		case int:
			target.SetFloat(float64(v))
			return true
		case string:
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				target.SetFloat(f)
				return true
			}
		}
	case reflect.String:
		target.SetString(fmt.Sprintf("%v", value))
		return true
	case reflect.Bool:
		switch v := value.(type) {
		case int:
			target.SetBool(v != 0)
			return true
		case float64:
			target.SetBool(v != 0)
			return true
		case string:
			if b, err := strconv.ParseBool(v); err == nil {
				target.SetBool(b)
				return true
			}
		}
	default:
		return false
	}
	return false
}
