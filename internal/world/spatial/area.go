package spatial

import "math"

// Area represents spatial region
type Area interface {
	// Shape returns area shape
	Shape() AreaShape

	// Center returns area center position
	Center() Position

	// Contains checks if position is in area
	Contains(pos Position) bool

	// GetPositions returns all positions in area
	GetPositions() []Position

	// Bounds returns bounding box
	Bounds() Rectangle

	// Size returns number of positions in area
	Size() int
}

// AreaShape defines region geometry
type AreaShape string

const (
	ShapeCircle    AreaShape = "circle"
	ShapeSquare    AreaShape = "square"
	ShapeRectangle AreaShape = "rectangle"
	ShapeLine      AreaShape = "line"
	ShapeCone      AreaShape = "cone"
	ShapeRing      AreaShape = "ring"
	ShapeCross     AreaShape = "cross"
	ShapeStar      AreaShape = "star"
)

// CircleArea represents circular region
type CircleArea struct {
	CenterPos Position
	Radius    float64
}

// NewCircleArea creates circular area
func NewCircleArea(center Position, radius float64) *CircleArea {
	return &CircleArea{CenterPos: center, Radius: radius}
}

// Shape returns area shape
func (c *CircleArea) Shape() AreaShape {
	return ShapeCircle
}

// Center returns center position
func (c *CircleArea) Center() Position {
	return c.CenterPos
}

// Contains checks if position is in circle
func (c *CircleArea) Contains(pos Position) bool {
	return c.CenterPos.DistanceTo(pos) <= c.Radius
}

// GetPositions returns all positions in circle
func (c *CircleArea) GetPositions() []Position {
	positions := make([]Position, 0)
	r := int(c.Radius)

	for dx := -r; dx <= r; dx++ {
		for dy := -r; dy <= r; dy++ {
			pos := c.CenterPos.Add(dx, dy, 0)
			if c.Contains(pos) {
				positions = append(positions, pos)
			}
		}
	}

	return positions
}

// Bounds returns bounding box
func (c *CircleArea) Bounds() Rectangle {
	r := int(c.Radius)
	return Rectangle{
		MinX: c.CenterPos.X - r,
		MinY: c.CenterPos.Y - r,
		MaxX: c.CenterPos.X + r,
		MaxY: c.CenterPos.Y + r,
	}
}

// Size returns approximate number of positions
func (c *CircleArea) Size() int {
	return int(math.Pi * c.Radius * c.Radius)
}

// RectangleArea represents rectangular region
type RectangleArea struct {
	Min Position
	Max Position
}

// NewRectangleArea creates rectangular area
func NewRectangleArea(min, max Position) *RectangleArea {
	return &RectangleArea{Min: min, Max: max}
}

// Shape returns area shape
func (r *RectangleArea) Shape() AreaShape {
	return ShapeRectangle
}

// Center returns center position
func (r *RectangleArea) Center() Position {
	return Position{
		X: (r.Min.X + r.Max.X) / 2,
		Y: (r.Min.Y + r.Max.Y) / 2,
		Z: (r.Min.Z + r.Max.Z) / 2,
	}
}

// Contains checks if position is in rectangle
func (r *RectangleArea) Contains(pos Position) bool {
	return pos.X >= r.Min.X && pos.X <= r.Max.X &&
		pos.Y >= r.Min.Y && pos.Y <= r.Max.Y &&
		pos.Z >= r.Min.Z && pos.Z <= r.Max.Z
}

// GetPositions returns all positions in rectangle
func (r *RectangleArea) GetPositions() []Position {
	positions := make([]Position, 0)

	for x := r.Min.X; x <= r.Max.X; x++ {
		for y := r.Min.Y; y <= r.Max.Y; y++ {
			for z := r.Min.Z; z <= r.Max.Z; z++ {
				positions = append(positions, Position{X: x, Y: y, Z: z})
			}
		}
	}

	return positions
}

// Bounds returns bounding box
func (r *RectangleArea) Bounds() Rectangle {
	return Rectangle{
		MinX: r.Min.X,
		MinY: r.Min.Y,
		MaxX: r.Max.X,
		MaxY: r.Max.Y,
	}
}

// Size returns number of positions
func (r *RectangleArea) Size() int {
	width := r.Max.X - r.Min.X + 1
	height := r.Max.Y - r.Min.Y + 1
	depth := r.Max.Z - r.Min.Z + 1
	return width * height * depth
}

// Rectangle represents 2D bounds
type Rectangle struct {
	MinX int
	MinY int
	MaxX int
	MaxY int
}

// Contains checks if position is in rectangle
func (r Rectangle) Contains(x, y int) bool {
	return x >= r.MinX && x <= r.MaxX && y >= r.MinY && y <= r.MaxY
}

// Width returns rectangle width
func (r Rectangle) Width() int {
	return r.MaxX - r.MinX + 1
}

// Height returns rectangle height
func (r Rectangle) Height() int {
	return r.MaxY - r.MinY + 1
}

// ConeArea represents cone-shaped region
type ConeArea struct {
	Origin    Position
	Direction Direction
	Length    float64
	Width     float64 // Cone width at end
}

// NewConeArea creates cone area
func NewConeArea(origin Position, direction Direction, length, width float64) *ConeArea {
	return &ConeArea{
		Origin:    origin,
		Direction: direction.Normalize(),
		Length:    length,
		Width:     width,
	}
}

// Shape returns area shape
func (c *ConeArea) Shape() AreaShape {
	return ShapeCone
}

// Center returns origin position
func (c *ConeArea) Center() Position {
	return c.Origin
}

// Contains checks if position is in cone
func (c *ConeArea) Contains(pos Position) bool {
	if pos.Z != c.Origin.Z {
		return false
	}

	// Calculate distance along cone axis
	dx := float64(pos.X - c.Origin.X)
	dy := float64(pos.Y - c.Origin.Y)

	dirX := float64(c.Direction.DX)
	dirY := float64(c.Direction.DY)

	// Project position onto cone direction
	distance := dx*dirX + dy*dirY

	if distance < 0 || distance > c.Length {
		return false
	}

	// Calculate perpendicular distance
	perpX := dx - distance*dirX
	perpY := dy - distance*dirY
	perpDist := math.Sqrt(perpX*perpX + perpY*perpY)

	// Calculate cone width at this distance
	widthAtDistance := (distance / c.Length) * c.Width

	return perpDist <= widthAtDistance
}

// GetPositions returns all positions in cone
func (c *ConeArea) GetPositions() []Position {
	positions := make([]Position, 0)
	bounds := c.Bounds()

	for x := bounds.MinX; x <= bounds.MaxX; x++ {
		for y := bounds.MinY; y <= bounds.MaxY; y++ {
			pos := Position{X: x, Y: y, Z: c.Origin.Z}
			if c.Contains(pos) {
				positions = append(positions, pos)
			}
		}
	}

	return positions
}

// Bounds returns bounding box
func (c *ConeArea) Bounds() Rectangle {
	maxExtent := int(math.Max(c.Length, c.Width))
	return Rectangle{
		MinX: c.Origin.X - maxExtent,
		MinY: c.Origin.Y - maxExtent,
		MaxX: c.Origin.X + maxExtent,
		MaxY: c.Origin.Y + maxExtent,
	}
}

// Size returns approximate number of positions
func (c *ConeArea) Size() int {
	return int(c.Length * c.Width / 2)
}

// LineArea represents line-shaped region
type LineArea struct {
	Start Position
	End   Position
	Width float64
}

// NewLineArea creates line area
func NewLineArea(start, end Position, width float64) *LineArea {
	return &LineArea{Start: start, End: end, Width: width}
}

// Shape returns area shape
func (l *LineArea) Shape() AreaShape {
	return ShapeLine
}

// Center returns midpoint
func (l *LineArea) Center() Position {
	return Position{
		X: (l.Start.X + l.End.X) / 2,
		Y: (l.Start.Y + l.End.Y) / 2,
		Z: (l.Start.Z + l.End.Z) / 2,
	}
}

// Contains checks if position is in line
func (l *LineArea) Contains(pos Position) bool {
	if pos.Z != l.Start.Z {
		return false
	}

	// Calculate distance from point to line segment
	dx := float64(l.End.X - l.Start.X)
	dy := float64(l.End.Y - l.Start.Y)
	length := math.Sqrt(dx*dx + dy*dy)

	if length == 0 {
		return l.Start.DistanceTo(pos) <= l.Width
	}

	// Normalize direction
	dx /= length
	dy /= length

	// Vector from start to point
	px := float64(pos.X - l.Start.X)
	py := float64(pos.Y - l.Start.Y)

	// Project onto line
	projection := px*dx + py*dy

	if projection < 0 || projection > length {
		return false
	}

	// Calculate perpendicular distance
	perpX := px - projection*dx
	perpY := py - projection*dy
	perpDist := math.Sqrt(perpX*perpX + perpY*perpY)

	return perpDist <= l.Width
}

// GetPositions returns all positions in line
func (l *LineArea) GetPositions() []Position {
	positions := make([]Position, 0)
	bounds := l.Bounds()

	for x := bounds.MinX; x <= bounds.MaxX; x++ {
		for y := bounds.MinY; y <= bounds.MaxY; y++ {
			pos := Position{X: x, Y: y, Z: l.Start.Z}
			if l.Contains(pos) {
				positions = append(positions, pos)
			}
		}
	}

	return positions
}

// Bounds returns bounding box
func (l *LineArea) Bounds() Rectangle {
	w := int(l.Width)
	return Rectangle{
		MinX: minInt(l.Start.X, l.End.X) - w,
		MinY: minInt(l.Start.Y, l.End.Y) - w,
		MaxX: maxInt(l.Start.X, l.End.X) + w,
		MaxY: maxInt(l.Start.Y, l.End.Y) + w,
	}
}

// Size returns approximate number of positions
func (l *LineArea) Size() int {
	distance := l.Start.DistanceTo(l.End)
	return int(distance * l.Width * 2)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
