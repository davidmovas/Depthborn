package entity

import "github.com/davidmovas/Depthborn/internal/core/types"

var _ types.Tagged = (*BaseTagged)(nil)

type BaseTagged struct {
	tagSet types.TagSet
}

func NewBaseTagged() *BaseTagged {
	return &BaseTagged{
		tagSet: NewBaseTagSet(),
	}
}

func (bt *BaseTagged) Tags() types.TagSet {
	return bt.tagSet
}

type BaseTagSet struct {
	tags map[string]bool
}

func NewBaseTagSet() *BaseTagSet {
	return &BaseTagSet{
		tags: make(map[string]bool),
	}
}

func (bts *BaseTagSet) Add(tag string) {
	bts.tags[tag] = true
}

func (bts *BaseTagSet) Remove(tag string) {
	delete(bts.tags, tag)
}

func (bts *BaseTagSet) Has(tag string) bool {
	return bts.tags[tag]
}

func (bts *BaseTagSet) Contains(tags ...string) bool {
	for _, tag := range tags {
		if !bts.tags[tag] {
			return false
		}
	}
	return true
}

func (bts *BaseTagSet) ContainsAny(tags ...string) bool {
	for _, tag := range tags {
		if bts.tags[tag] {
			return true
		}
	}
	return false
}

func (bts *BaseTagSet) All() []string {
	result := make([]string, 0, len(bts.tags))
	for tag := range bts.tags {
		result = append(result, tag)
	}
	return result
}

func (bts *BaseTagSet) Clear() {
	bts.tags = make(map[string]bool)
}
