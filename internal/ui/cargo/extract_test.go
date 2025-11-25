package cargo

import (
	"reflect"
	"testing"
)

type LivingConfigTest struct {
	Health    float64
	MaxHealth float64
	Tags      []int
	Meta      map[string]string
}

// Helper function to create a sample instance for testing
func sampleLivingConfig() LivingConfigTest {
	return LivingConfigTest{
		Health:    100.5,
		MaxHealth: 200.0,
		Tags:      []int{1, 2, 3},
		Meta:      map[string]string{"role": "warrior"},
	}
}

func TestExtract(t *testing.T) {
	SetGlobal(sampleLivingConfig())

	t.Run("ExtractKeys", func(t *testing.T) {
		t.Run("Basic fields", func(t *testing.T) {
			var hp, maxHp float64
			var tags []int
			Extract[LivingConfigTest](
				"health", &hp,
				"MAXHEALTH", &maxHp,
				"tags", &tags,
			)

			if hp != 100.5 {
				t.Errorf("hp = %v, want 100.5", hp)
			}
			if maxHp != 200.0 {
				t.Errorf("maxHp = %v, want 200.0", maxHp)
			}
			if !reflect.DeepEqual(tags, []int{1, 2, 3}) {
				t.Errorf("tags = %v, want [1 2 3]", tags)
			}
		})

		t.Run("Case-insensitive keys", func(t *testing.T) {
			var hp, maxHp float64
			Extract[LivingConfigTest](
				"HEALTH", &hp,
				"maxhealth", &maxHp,
			)
			if hp != 100.5 || maxHp != 200.0 {
				t.Errorf("case-insensitive keys failed")
			}
		})

		t.Run("Missing field sets zero", func(t *testing.T) {
			var unknown int
			Extract[LivingConfigTest](
				"notexist", &unknown,
			)
			if unknown != 0 {
				t.Errorf("unknown = %v, want 0 for missing field", unknown)
			}
		})
	})

	t.Run("ExtractInto", func(t *testing.T) {
		t.Run("Basic assignment", func(t *testing.T) {
			var out struct {
				Health float64
				Tags   []int
			}
			ExtractInto[LivingConfigTest](&out)
			if out.Health != 100.5 {
				t.Errorf("Health = %v, want 100.5", out.Health)
			}
			if !reflect.DeepEqual(out.Tags, []int{1, 2, 3}) {
				t.Errorf("Tags = %v, want [1 2 3]", out.Tags)
			}
		})

		t.Run("Missing fields zeroed", func(t *testing.T) {
			var out struct {
				Unknown int
			}
			ExtractInto[LivingConfigTest](&out)
			if out.Unknown != 0 {
				t.Errorf("Unknown = %v, want 0", out.Unknown)
			}
		})
	})
}
