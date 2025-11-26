package style

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// GradientBuilder строит градиент из цветов
type GradientBuilder struct {
	colors    []Color
	steps     int
	direction GradientDirection
}

// GradientDirection определяет направление градиента
type GradientDirection int

const (
	// Для текста (Foreground/Background)
	DirectionLeftToRight  GradientDirection = iota // По умолчанию →
	DirectionRightToLeft                           // ←
	DirectionTopToBottom                           // ↓
	DirectionBottomToTop                           // ↑
	DirectionDiagonalTLBR                          // ↘ (Top-Left to Bottom-Right)
	DirectionDiagonalTRBL                          // ↙ (Top-Right to Bottom-Left)
	DirectionDiagonalBLTR                          // ↗ (Bottom-Left to Top-Right)
	DirectionDiagonalBRTL                          // ↖ (Bottom-Right to Top-Left)

	// Для border
	DirectionBorderClockwise        // По часовой: Top → Right → Bottom → Left
	DirectionBorderCounterClockwise // Против часовой: Top → Left → Bottom → Right
	DirectionBorderVertical         // Вертикально: Top и Bottom одинаковые, Left и Right одинаковые
	DirectionBorderHorizontal       // Горизонтально: Left и Right одинаковые, Top и Bottom одинаковые
)

// NewGradient создает новый градиент
func NewGradient(colors ...Color) *GradientBuilder {
	return &GradientBuilder{
		colors:    colors,
		direction: DirectionLeftToRight, // По умолчанию
	}
}

// Direction устанавливает направление градиента
func (g *GradientBuilder) Direction(dir GradientDirection) *GradientBuilder {
	g.direction = dir
	return g
}

// Steps устанавливает количество шагов градиента
func (g *GradientBuilder) Steps(steps int) *GradientBuilder {
	g.steps = steps
	return g
}

// Foreground применяет градиент к тексту (цвет текста)
// Умная версия - пропускает уже стилизованные части
func (g *GradientBuilder) Foreground(text string) string {
	// Для текста учитываем направление
	result := g.applySmartGradient(text, "foreground")
	return g.applyDirection(result, text)
}

// applyDirection применяет направление к уже отрендеренному тексту
func (g *GradientBuilder) applyDirection(renderedText, originalText string) string {
	// Для простых направлений (left-to-right, right-to-left) можно просто реверсировать
	switch g.direction {
	case DirectionRightToLeft:
		// Реверсируем градиент
		return g.reverseGradient(originalText, "foreground")
	default:
		return renderedText
	}
}

// reverseGradient создает градиент в обратном направлении
func (g *GradientBuilder) reverseGradient(text string, applyTo string) string {
	runes := []rune(text)
	if len(runes) == 0 {
		return ""
	}

	// Генерируем градиент и реверсируем его
	gradient := g.generate(len(runes))
	for i, j := 0, len(gradient)-1; i < j; i, j = i+1, j-1 {
		gradient[i], gradient[j] = gradient[j], gradient[i]
	}

	var result string
	for i, char := range runes {
		var style lipgloss.Style
		switch applyTo {
		case "foreground":
			style = lipgloss.NewStyle().Foreground(gradient[i])
		case "background":
			style = lipgloss.NewStyle().Background(gradient[i])
		}
		result += style.Render(string(char))
	}

	return result
}

// ForegroundRaw применяет градиент ко всему тексту без учета стилей
func (g *GradientBuilder) ForegroundRaw(text string) string {
	runes := []rune(text)
	if len(runes) == 0 {
		return ""
	}

	gradient := g.generate(len(runes))
	var result string

	for i, char := range runes {
		style := lipgloss.NewStyle().Foreground(gradient[i])
		result += style.Render(string(char))
	}

	return result
}

// Background применяет градиент как фон
func (g *GradientBuilder) Background(text string) string {
	return g.applySmartGradient(text, "background")
}

// BackgroundRaw применяет градиент как фон ко всему тексту
func (g *GradientBuilder) BackgroundRaw(text string) string {
	runes := []rune(text)
	if len(runes) == 0 {
		return ""
	}

	gradient := g.generate(len(runes))
	var result string

	for i, char := range runes {
		style := lipgloss.NewStyle().Background(gradient[i])
		result += style.Render(string(char))
	}

	return result
}

// applySmartGradient применяет градиент, пропуская ANSI escape последовательности
func (g *GradientBuilder) applySmartGradient(text string, applyTo string) string {
	if len(text) == 0 {
		return ""
	}

	// Паттерн для ANSI escape последовательностей
	ansiPattern := regexp.MustCompile(`\x1b\[[0-9;]*m`)

	// Разбиваем текст на части: обычный текст и ANSI коды
	parts := []struct {
		text   string
		isAnsi bool
	}{}

	lastIndex := 0
	matches := ansiPattern.FindAllStringIndex(text, -1)

	for _, match := range matches {
		// Добавляем обычный текст перед ANSI кодом
		if match[0] > lastIndex {
			parts = append(parts, struct {
				text   string
				isAnsi bool
			}{text[lastIndex:match[0]], false})
		}
		// Добавляем ANSI код
		parts = append(parts, struct {
			text   string
			isAnsi bool
		}{text[match[0]:match[1]], true})
		lastIndex = match[1]
	}

	// Добавляем оставшийся текст
	if lastIndex < len(text) {
		parts = append(parts, struct {
			text   string
			isAnsi bool
		}{text[lastIndex:], false})
	}

	// Если нет обычного текста (весь текст - ANSI коды), возвращаем как есть
	if len(parts) == 0 {
		return text
	}

	// Подсчитываем количество видимых символов
	visibleChars := 0
	for _, part := range parts {
		if !part.isAnsi {
			visibleChars += len([]rune(part.text))
		}
	}

	if visibleChars == 0 {
		return text
	}

	// Генерируем градиент для видимых символов
	gradient := g.generate(visibleChars)

	// Применяем градиент
	var result strings.Builder
	gradientIndex := 0

	for _, part := range parts {
		if part.isAnsi {
			// ANSI коды оставляем как есть
			result.WriteString(part.text)
		} else {
			// Применяем градиент к обычному тексту
			for _, char := range part.text {
				if gradientIndex >= len(gradient) {
					break
				}

				var style lipgloss.Style
				switch applyTo {
				case "foreground":
					style = lipgloss.NewStyle().Foreground(gradient[gradientIndex])
				case "background":
					style = lipgloss.NewStyle().Background(gradient[gradientIndex])
				}

				result.WriteString(style.Render(string(char)))
				gradientIndex++
			}
		}
	}

	return result.String()
}

// ForegroundStyle возвращает массив стилей с градиентом для foreground
func (g *GradientBuilder) ForegroundStyle(steps int) []Style {
	gradient := g.generate(steps)
	styles := make([]Style, steps)

	for i, color := range gradient {
		styles[i] = lipgloss.NewStyle().Foreground(color)
	}

	return styles
}

// BackgroundStyle возвращает массив стилей с градиентом для background
func (g *GradientBuilder) BackgroundStyle(steps int) []Style {
	gradient := g.generate(steps)
	styles := make([]Style, steps)

	for i, color := range gradient {
		styles[i] = lipgloss.NewStyle().Background(color)
	}

	return styles
}

// Style применяет градиент к базовому стилю
// baseStyle - ваш заранее настроенный стиль (с border, padding и т.д.)
// applyTo - куда применить: "foreground"/"fg", "background"/"bg", "border"
// Для border применяется первый цвет градиента (если нужен полный градиент - используйте BorderGradient)
func (g *GradientBuilder) Style(baseStyle Style, applyTo string) Style {
	gradient := g.generate(1)
	if len(gradient) == 0 {
		return baseStyle
	}

	result := baseStyle.Copy()

	switch applyTo {
	case "foreground", "fg", "color":
		result = result.Foreground(gradient[0])
	case "background", "bg":
		result = result.Background(gradient[0])
	case "border":
		result = result.BorderForeground(gradient[0])
	}

	return result
}

// BorderGradient применяет градиент к каждой стороне border
// ВАЖНО: lipgloss позволяет установить только один цвет на каждую сторону border
// Для настоящего градиента используйте BorderGradientBox()
func (g *GradientBuilder) BorderGradient(baseStyle Style) Style {
	var topColor, rightColor, bottomColor, leftColor Color

	switch g.direction {
	case DirectionBorderClockwise:
		// Top → Right → Bottom → Left (по часовой)
		gradient := g.generate(4)
		if len(gradient) >= 4 {
			topColor = gradient[0]
			rightColor = gradient[1]
			bottomColor = gradient[2]
			leftColor = gradient[3]
		}

	case DirectionBorderCounterClockwise:
		// Top → Left → Bottom → Right (против часовой)
		gradient := g.generate(4)
		if len(gradient) >= 4 {
			topColor = gradient[0]
			leftColor = gradient[1]
			bottomColor = gradient[2]
			rightColor = gradient[3]
		}

	case DirectionBorderVertical:
		// Вертикальный градиент: Top → Bottom
		// Left и Right получают интерполированные значения
		gradient := g.generate(4)
		if len(gradient) >= 4 {
			topColor = gradient[0]
			bottomColor = gradient[3]
			// Left и Right - средние значения
			leftColor = gradient[1]
			rightColor = gradient[2]
		}

	case DirectionBorderHorizontal:
		// Горизонтальный градиент: Left → Right
		// Top и Bottom получают интерполированные значения
		gradient := g.generate(4)
		if len(gradient) >= 4 {
			leftColor = gradient[0]
			rightColor = gradient[3]
			// Top и Bottom - средние значения
			topColor = gradient[1]
			bottomColor = gradient[2]
		}

	case DirectionDiagonalTLBR:
		// Top-Left to Bottom-Right: градиент от левого верхнего к правому нижнему
		gradient := g.generate(4)
		if len(gradient) >= 4 {
			topColor = gradient[0]
			leftColor = gradient[1]
			rightColor = gradient[2]
			bottomColor = gradient[3]
		}

	case DirectionDiagonalBLTR:
		// Bottom-Left to Top-Right: градиент от левого нижнего к правому верхнему
		gradient := g.generate(4)
		if len(gradient) >= 4 {
			bottomColor = gradient[0]
			leftColor = gradient[1]
			topColor = gradient[2]
			rightColor = gradient[3]
		}

	case DirectionDiagonalTRBL:
		// Top-Right to Bottom-Left: градиент от правого верхнего к левому нижнему
		gradient := g.generate(4)
		if len(gradient) >= 4 {
			topColor = gradient[0]
			rightColor = gradient[1]
			leftColor = gradient[2]
			bottomColor = gradient[3]
		}

	case DirectionDiagonalBRTL:
		// Bottom-Right to Top-Left: градиент от правого нижнего к левому верхнему
		gradient := g.generate(4)
		if len(gradient) >= 4 {
			bottomColor = gradient[0]
			rightColor = gradient[1]
			topColor = gradient[2]
			leftColor = gradient[3]
		}

	default:
		// По умолчанию - по часовой
		gradient := g.generate(4)
		if len(gradient) >= 4 {
			topColor = gradient[0]
			rightColor = gradient[1]
			bottomColor = gradient[2]
			leftColor = gradient[3]
		}
	}

	result := baseStyle.Copy()

	// Применяем цвета
	result = result.
		BorderTopForeground(topColor).
		BorderRightForeground(rightColor).
		BorderBottomForeground(bottomColor).
		BorderLeftForeground(leftColor)

	return result
}

// BorderGradientBox создает рамку с настоящим градиентом, используя символы
// Это позволяет создать плавный градиент вдоль каждой стороны
func (g *GradientBuilder) BorderGradientBox(content string, width int) string {
	lines := strings.Split(content, "\n")
	height := len(lines)

	// Периметр = верх + право + низ + лево (без учета углов, они между сегментами)
	// +4 для углов в градиенте
	perimeterLength := width + height + width + height + 4
	gradient := g.generate(perimeterLength)

	if len(gradient) == 0 {
		gradient = []Color{lipgloss.Color("#FFFFFF")}
	}

	var result strings.Builder
	gradIndex := 0

	// Функция для безопасного получения цвета
	getColor := func() Color {
		if gradIndex >= len(gradient) {
			return gradient[len(gradient)-1]
		}
		color := gradient[gradIndex]
		gradIndex++
		return color
	}

	// Верхний левый угол
	cornerStyle := lipgloss.NewStyle().Foreground(getColor())
	result.WriteString(cornerStyle.Render("┌"))

	// Верхняя граница
	for i := 0; i < width; i++ {
		s := lipgloss.NewStyle().Foreground(getColor())
		result.WriteString(s.Render("─"))
	}

	// Верхний правый угол
	cornerStyle = lipgloss.NewStyle().Foreground(getColor())
	result.WriteString(cornerStyle.Render("┐"))
	result.WriteString("\n")

	// Средние строки
	for i := 0; i < height; i++ {
		// Левая граница
		leftStyle := lipgloss.NewStyle().Foreground(getColor())
		result.WriteString(leftStyle.Render("│"))

		// Контент
		if i < len(lines) {
			// Обрезаем или дополняем строку до нужной ширины
			line := lines[i]
			runes := []rune(line)
			if len(runes) > width {
				result.WriteString(string(runes[:width]))
			} else {
				result.WriteString(line)
				result.WriteString(strings.Repeat(" ", width-len(runes)))
			}
		} else {
			result.WriteString(strings.Repeat(" ", width))
		}

		// Правая граница
		rightStyle := lipgloss.NewStyle().Foreground(getColor())
		result.WriteString(rightStyle.Render("│"))
		result.WriteString("\n")
	}

	// Нижний левый угол
	cornerStyle = lipgloss.NewStyle().Foreground(getColor())
	result.WriteString(cornerStyle.Render("└"))

	// Нижняя граница
	for i := 0; i < width; i++ {
		s := lipgloss.NewStyle().Foreground(getColor())
		result.WriteString(s.Render("─"))
	}

	// Нижний правый угол
	cornerStyle = lipgloss.NewStyle().Foreground(getColor())
	result.WriteString(cornerStyle.Render("┘"))

	return result.String()
}

// Styles возвращает массив стилей с примененным градиентом
// baseStyle - ваш заранее настроенный стиль
// steps - количество вариаций
// applyTo - куда применить градиент
func (g *GradientBuilder) Styles(baseStyle Style, steps int, applyTo string) []Style {
	gradient := g.generate(steps)
	styles := make([]Style, steps)

	for i, color := range gradient {
		result := baseStyle.Copy()

		switch applyTo {
		case "foreground", "fg", "color":
			result = result.Foreground(color)
		case "background", "bg":
			result = result.Background(color)
		case "border":
			result = result.BorderForeground(color)
		}

		styles[i] = result
	}

	return styles
}

// Colors возвращает массив интерполированных цветов
func (g *GradientBuilder) Colors(steps int) []Color {
	return g.generate(steps)
}

// generate создает градиент с указанным количеством шагов
func (g *GradientBuilder) generate(steps int) []Color {
	if len(g.colors) == 0 {
		return []Color{}
	}
	if len(g.colors) == 1 {
		result := make([]Color, steps)
		for i := range result {
			result[i] = g.colors[0]
		}
		return result
	}
	if steps <= 0 {
		return []Color{}
	}
	if steps == 1 {
		return []Color{g.colors[0]}
	}

	gradient := make([]Color, steps)

	for i := 0; i < steps; i++ {
		// Нормализованная позиция от 0.0 до 1.0
		position := float64(i) / float64(steps-1)

		// Размер одного сегмента
		segmentSize := 1.0 / float64(len(g.colors)-1)

		// Определяем индекс сегмента
		segmentIndex := int(position / segmentSize)

		// Защита от выхода за границы
		if segmentIndex < 0 {
			segmentIndex = 0
		}
		if segmentIndex >= len(g.colors)-1 {
			segmentIndex = len(g.colors) - 2
		}

		// Локальная позиция внутри сегмента (0.0 - 1.0)
		localPosition := (position - float64(segmentIndex)*segmentSize) / segmentSize

		// Защита от NaN и Inf
		if localPosition < 0 {
			localPosition = 0
		}
		if localPosition > 1 {
			localPosition = 1
		}

		color1 := g.colors[segmentIndex]
		color2 := g.colors[segmentIndex+1]

		gradient[i] = interpolateColor(color1, color2, localPosition)
	}

	return gradient
}

// Утилиты для быстрого использования

// GradientFg применяет градиент к тексту (foreground), пропуская стилизованные части
func GradientFg(text string, colors ...Color) string {
	return NewGradient(colors...).Foreground(text)
}

// GradientBg применяет градиент как фон, пропуская стилизованные части
func GradientBg(text string, colors ...Color) string {
	return NewGradient(colors...).Background(text)
}

// GradientFgRaw применяет градиент ко всему тексту
func GradientFgRaw(text string, colors ...Color) string {
	return NewGradient(colors...).ForegroundRaw(text)
}

// GradientBgRaw применяет градиент как фон ко всему тексту
func GradientBgRaw(text string, colors ...Color) string {
	return NewGradient(colors...).BackgroundRaw(text)
}

// GradientBorder создает стиль с градиентной рамкой из базового стиля
// ВАЖНО: каждая сторона border будет одного цвета (ограничение lipgloss)
// Для настоящего плавного градиента используйте GradientBorderBox
func GradientBorder(baseStyle Style, colors ...Color) Style {
	return NewGradient(colors...).Style(baseStyle, "border")
}

// GradientBorderFull создает стиль с полным градиентом на рамке
// ВАЖНО: каждая сторона border будет одного цвета (ограничение lipgloss)
// Для настоящего плавного градиента используйте GradientBorderBox
func GradientBorderFull(baseStyle Style, colors ...Color) Style {
	return NewGradient(colors...).BorderGradient(baseStyle)
}

// GradientBorderBox создает рамку с настоящим плавным градиентом
// Это рисует рамку символами, что позволяет создать плавный градиент
func GradientBorderBox(content string, width int, colors ...Color) string {
	return NewGradient(colors...).BorderGradientBox(content, width)
}

// GradientStyles возвращает массив стилей с градиентом
func GradientStyles(baseStyle Style, steps int, applyTo string, colors ...Color) []Style {
	return NewGradient(colors...).Styles(baseStyle, steps, applyTo)
}

func interpolateColor(color1, color2 Color, t float64) Color {
	r1, g1, b1, _ := color1.RGBA()
	r2, g2, b2, _ := color2.RGBA()

	r18 := uint8(r1 >> 8)
	g18 := uint8(g1 >> 8)
	b18 := uint8(b1 >> 8)

	r28 := uint8(r2 >> 8)
	g28 := uint8(g2 >> 8)
	b28 := uint8(b2 >> 8)

	r := uint8(float64(r18) + (float64(r28)-float64(r18))*t)
	g := uint8(float64(g18) + (float64(g28)-float64(g18))*t)
	b := uint8(float64(b18) + (float64(b28)-float64(b18))*t)

	return lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", r, g, b))
}
