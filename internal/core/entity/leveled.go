package entity

type BaseLeveled struct {
	level int
}

func NewBaseLeveled(level int) *BaseLeveled {
	return &BaseLeveled{
		level: level,
	}
}

func (bl *BaseLeveled) Level() int {
	return bl.level
}

func (bl *BaseLeveled) SetLevel(level int) {
	if level < 1 {
		level = 1
	}
	bl.level = level
}
