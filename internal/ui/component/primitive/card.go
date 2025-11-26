package primitive

import (
	"strings"

	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/style"
)

type LabelPosition int

const (
	LabelTopLeft LabelPosition = iota
	LabelTopCenter
	LabelTopRight
)

type CardProps struct {
	CommonProps
	Children      []component.Component
	Label         *string        // Стилизованная строка лейбла
	LabelPosition *LabelPosition // Только верхние позиции
	Padding       *int
	Margin        *int
	Width         *int
	Height        *int
	BorderStyle   *style.BoxBorderStyle
}

func Card(props CardProps) component.Component {
	return component.Func(func(ctx *component.Context) string {
		// Рендерим детей
		childrenOutput := ""
		for _, child := range props.Children {
			childrenOutput += child.Render(ctx)
		}

		if props.Label != nil {
			borderStyle := style.BoxBorderRounded
			if props.BorderStyle != nil {
				borderStyle = *props.BorderStyle
			}
			return renderCardWithTopLabel(*props.Label, childrenOutput, props.LabelPosition, props, borderStyle)
		}

		// Карточка без лейбла
		return renderSimpleCard(childrenOutput, props)
	})
}

func renderCardWithTopLabel(label, content string, position *LabelPosition, props CardProps, borderStyle style.BoxBorderStyle) string {
	pos := LabelTopLeft
	if position != nil {
		pos = *position
	}

	// Получаем чистый контент для расчета ширины
	cleanContent := style.StripAnsi(content)
	contentLines := strings.Split(cleanContent, "\n")

	// Вычисляем максимальную ширину контента
	maxContentWidth := 0
	for _, line := range contentLines {
		width := style.CalculateWidth(line)
		if width > maxContentWidth {
			maxContentWidth = width
		}
	}

	// Применяем padding
	padding := 0
	if props.Padding != nil {
		padding = *props.Padding
	}

	// Общая ширина = контент + padding слева и справа
	totalWidth := maxContentWidth + padding*2

	// Чистый лейбл для расчета ширины
	cleanLabel := style.StripAnsi(label)
	labelWidth := style.CalculateWidth(cleanLabel)

	// Строим карточку
	var result strings.Builder

	// 1. Верхняя граница с лейблом
	result.WriteString(buildTopBorderWithLabel(borderStyle, label, labelWidth, totalWidth, pos))
	result.WriteString("\n")

	// 2. Верхний отступ (пустые строки)
	for i := 0; i < padding; i++ {
		result.WriteString(borderStyle.Left)
		result.WriteString(strings.Repeat(" ", totalWidth))
		result.WriteString(borderStyle.Right)
		result.WriteString("\n")
	}

	// 3. Контент
	originalContentLines := strings.Split(content, "\n")
	for _, line := range originalContentLines {
		result.WriteString(borderStyle.Left)

		// Чистая ширина этой строки
		cleanLine := style.StripAnsi(line)
		cleanWidth := style.CalculateWidth(cleanLine)

		// Вычисляем отступы для центрирования
		leftPadding := padding
		rightPadding := totalWidth - cleanWidth - leftPadding

		// Левый отступ
		result.WriteString(strings.Repeat(" ", leftPadding))
		// Контент (оригинал со стилями)
		result.WriteString(line)
		// Правый отступ
		if rightPadding > 0 {
			result.WriteString(strings.Repeat(" ", rightPadding))
		}

		result.WriteString(borderStyle.Right)
		result.WriteString("\n")
	}

	// 4. Нижний отступ (пустые строки)
	for i := 0; i < padding; i++ {
		result.WriteString(borderStyle.Left)
		result.WriteString(strings.Repeat(" ", totalWidth))
		result.WriteString(borderStyle.Right)
		result.WriteString("\n")
	}

	// 5. Нижняя граница
	result.WriteString(buildBottomBorder(borderStyle, totalWidth))

	// Применяем margin (с дефолтным значением если не указан)
	return applyMarginWithDefault(result.String(), props)
}

func renderSimpleCard(content string, props CardProps) string {
	cardStyle := buildCardStyle(props)
	return applyMarginWithDefault(cardStyle.Render(content), props)
}

func buildTopBorder(border style.BoxBorderStyle, width int) string {
	return border.TopLeft + strings.Repeat(border.Top, width) + border.TopRight
}

func buildTopBorderWithLabel(border style.BoxBorderStyle, styledLabel string, labelWidth, totalWidth int, position LabelPosition) string {
	var result strings.Builder

	switch position {
	case LabelTopLeft:
		// ┌─[Label]─────┐
		result.WriteString(border.TopLeft)
		result.WriteString(strings.Repeat(border.Top, 2)) // Отступ 2 символа
		result.WriteString(styledLabel)                   // Стилизованный лейбл
		remainingWidth := totalWidth - labelWidth - 2
		if remainingWidth > 0 {
			result.WriteString(strings.Repeat(border.Top, remainingWidth))
		}
		result.WriteString(border.TopRight)

	case LabelTopCenter:
		// ┌───[Label]───┐
		sideWidth := (totalWidth - labelWidth) / 2
		result.WriteString(border.TopLeft)
		result.WriteString(strings.Repeat(border.Top, sideWidth))
		result.WriteString(styledLabel) // Стилизованный лейбл
		remainingWidth := totalWidth - sideWidth - labelWidth
		if remainingWidth > 0 {
			result.WriteString(strings.Repeat(border.Top, remainingWidth))
		}
		result.WriteString(border.TopRight)

	case LabelTopRight:
		// ┌─────[Label]─┐
		result.WriteString(border.TopLeft)
		remainingWidth := totalWidth - labelWidth - 2
		if remainingWidth > 0 {
			result.WriteString(strings.Repeat(border.Top, remainingWidth))
		}
		result.WriteString(styledLabel)                   // Стилизованный лейбл
		result.WriteString(strings.Repeat(border.Top, 2)) // Отступ 2 символа
		result.WriteString(border.TopRight)
	}

	return result.String()
}

func buildBottomBorder(border style.BoxBorderStyle, width int) string {
	return border.BottomLeft + strings.Repeat(border.Bottom, width) + border.BottomRight
}

func buildCardStyle(props CardProps) style.Style {
	cardStyle := style.Merge(
		style.BorderRounded,
		style.Padding2,
	)

	if props.Padding != nil {
		cardStyle = cardStyle.Padding(*props.Padding)
	}
	if props.Margin != nil {
		cardStyle = cardStyle.Margin(*props.Margin)
	}
	if props.Width != nil {
		cardStyle = cardStyle.Width(*props.Width)
	}
	if props.Height != nil {
		cardStyle = cardStyle.Height(*props.Height)
	}

	if props.Style != nil {
		cardStyle = cardStyle.Inherit(*props.Style)
	}

	return cardStyle
}

func applyMargin(content string, props CardProps) string {
	if props.Margin != nil && *props.Margin > 0 {
		marginStyle := style.New().Margin(*props.Margin)
		return marginStyle.Render(content)
	}
	return content
}

// applyMarginWithDefault применяет margin с дефолтным значением
func applyMarginWithDefault(content string, props CardProps) string {
	margin := 1 // Дефолтный вертикальный margin
	if props.Margin != nil {
		margin = *props.Margin
	}

	if margin > 0 {
		marginStyle := style.New().MarginTop(margin).MarginBottom(margin)
		return marginStyle.Render(content)
	}
	return content
}
