package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/davidmovas/Depthborn/internal/ui/component"
	. "github.com/davidmovas/Depthborn/internal/ui/component/primitive"
	"github.com/davidmovas/Depthborn/internal/ui/style"
)

// Demo creates a comprehensive demo component showing various primitives
func Demo(ctx *component.Context) component.Component {
	// State
	counter := component.UseState(ctx, 0, "counter")
	selectedTab := component.UseState(ctx, 0, "tab")
	inputValue := component.UseState(ctx, "", "input")
	checkboxValue := component.UseState(ctx, false, "checkbox")
	modalOpen := component.UseState(ctx, false, "modal")

	// Screen info
	width, height := ctx.ScreenSize()

	return component.Func(func(ctx *component.Context) string {
		// Header with debug info
		header := renderHeader(ctx, width, height, counter.Get())

		// Tab buttons (row 0)
		tabs := renderTabs(ctx, selectedTab)

		// Content based on selected tab (row 1+ after FocusRowBreak)
		var content component.Component
		switch selectedTab.Get() {
		case 0:
			content = renderButtonsDemo(ctx, counter)
		case 1:
			content = renderInputDemo(ctx, inputValue, checkboxValue)
		case 2:
			content = renderLayoutDemo(ctx)
		case 3:
			content = renderModalDemo(ctx, modalOpen)
		default:
			content = Text(TextProps{Content: "Unknown tab"})
		}

		// Footer with instructions
		footer := renderFooter()

		// Combine all - FocusRowBreak() moves to next row for navigation
		mainContent := VStack(ContainerProps{
			ChildrenProps: Children(
				component.Raw(header),
				VSpacer(1),
				tabs,
				component.FocusRowBreak(), // <- Row break between tabs and content
				VSpacer(1),
				content,
				VSpacer(1),
				footer,
			),
			LayoutProps: LayoutProps{
				Width: width,
			},
		}, 0)

		// Modal overlay
		modal := Modal(ModalProps{
			Open:  modalOpen.Get(),
			Title: "Example Modal",
			OnClose: func() {
				modalOpen.Set(false)
			},
			Size:           ModalSizeSmall,
			CloseOnEscape:  true,
			CloseOnOverlay: true,
			ShowCloseBtn:   true,
			Overlay:        true,
			ContainerProps: ContainerProps{
				ChildrenProps: Children(
					Text(TextProps{Content: "This is a modal window!"}),
					VSpacer(1),
					Text(TextProps{Content: "Press ESC or click Close to dismiss."}),
				),
			},
			Footer: []component.Component{
				Button(InteractiveProps{
					FocusProps: FocusProps{
						AutoFocus: true,
						OnClick: func() {
							modalOpen.Set(false)
						},
					},
				}, "[Close]"),
			},
		})

		return VStack(ContainerProps{
			ChildrenProps: Children(mainContent, modal),
		}, 0).Render(ctx)
	})
}

func renderHeader(ctx *component.Context, width, height int, counter int) string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(style.Primary).
		Padding(0, 1)

	infoStyle := lipgloss.NewStyle().
		Foreground(style.Grey500)

	title := titleStyle.Render("Depthborn UI Demo")

	// Debug: show current focus position
	pos := ctx.Focus().CurrentPosition()
	posStr := "nil"
	if pos != nil {
		posStr = fmt.Sprintf("r%d,c%d", pos.Row, pos.Col)
	}
	info := infoStyle.Render(fmt.Sprintf("Screen: %dx%d | Counter: %d | Focus: %s", width, height, counter, posStr))

	return lipgloss.JoinHorizontal(lipgloss.Top, title, "  ", info)
}

func renderTabs(ctx *component.Context, selectedTab *component.State[int]) component.Component {
	tabs := []string{"Buttons", "Inputs", "Layout", "Modal"}

	var tabComponents []component.Component
	for i, tabName := range tabs {
		idx := i // capture for closure
		isSelected := selectedTab.Get() == idx

		// Minimalist tab styling: cyan text, bold when selected
		var tabStyle lipgloss.Style
		if isSelected {
			tabStyle = lipgloss.NewStyle().
				Foreground(style.InteractiveFocus).
				Bold(true).
				PaddingLeft(2).
				PaddingRight(2)
		} else {
			tabStyle = lipgloss.NewStyle().
				Foreground(style.Interactive).
				PaddingLeft(2).
				PaddingRight(2)
		}

		tabComponents = append(tabComponents, Button(InteractiveProps{
			StyleProps: WithStyle(tabStyle),
			FocusProps: FocusProps{
				AutoFocus: idx == 0,
				OnClick: func() {
					selectedTab.Set(idx)
				},
			},
		}, tabName))

		// Add spacer between tabs
		if i < len(tabs)-1 {
			tabComponents = append(tabComponents, HSpacer(1))
		}
	}

	return HStack(ContainerProps{
		ChildrenProps: ChildrenProps{Children: tabComponents},
	}, 0)
}

func renderButtonsDemo(ctx *component.Context, counter *component.State[int]) component.Component {
	return VStack(ContainerProps{
		ChildrenProps: Children(
			Heading(TextProps{Content: "Button Examples"}),
			VSpacer(1),

			// Row of buttons - minimalist styling (uses default Button style now)
			HStack(ContainerProps{
				ChildrenProps: Children(
					Button(InteractiveProps{
						FocusProps: FocusProps{
							OnClick: func() {
								counter.Set(counter.Get() + 1)
							},
						},
					}, "[+] Increment"),

					Button(InteractiveProps{
						FocusProps: FocusProps{
							OnClick: func() {
								counter.Set(counter.Get() - 1)
							},
						},
					}, "[-] Decrement"),

					Button(InteractiveProps{
						FocusProps: FocusProps{
							OnClick: func() {
								counter.Set(0)
							},
						},
					}, "[x] Reset"),
				),
			}, 2),

			VSpacer(1),
			Text(TextProps{
				Content: fmt.Sprintf("Counter value: %d", counter.Get()),
				StyleProps: StyleProps{
					Style: Ptr(style.Merge(style.Bold, style.Fg(style.Interactive))),
				},
			}),

			VSpacer(1),
			Divider(DividerProps{
				BaseProps: BaseProps{LayoutProps: LayoutProps{Width: 50}},
				Label:     "Badge Examples",
			}),
			VSpacer(1),

			HStack(ContainerProps{
				ChildrenProps: Children(
					Badge(BadgeProps{Variant: BadgeDefault, Content: "Default"}),
					Badge(BadgeProps{Variant: BadgeSuccess, Content: "Success"}),
					Badge(BadgeProps{Variant: BadgeWarning, Content: "Warning"}),
					Badge(BadgeProps{Variant: BadgeError, Content: "Error"}),
					Badge(BadgeProps{Variant: BadgeInfo, Content: "Info"}),
				),
			}, 1),
		),
	}, 1)
}

func renderInputDemo(ctx *component.Context, inputValue *component.State[string], checkboxValue *component.State[bool]) component.Component {
	return VStack(ContainerProps{
		ChildrenProps: Children(
			Heading(TextProps{Content: "Input Examples"}),
			VSpacer(1),

			// Text input
			VStack(ContainerProps{
				ChildrenProps: Children(
					Label(TextProps{Content: "Text Input:"}),
					Input(InputProps{
						BaseProps: BaseProps{
							LayoutProps: LayoutProps{Width: 30},
						},
						Value:       inputValue.Get(),
						Placeholder: "Type something...",
						OnChange: func(val string) {
							inputValue.Set(val)
						},
					}),
				),
			}, 0),

			VSpacer(1),
			Text(TextProps{
				Content: fmt.Sprintf("Input value: %q", inputValue.Get()),
				StyleProps: StyleProps{
					Style: Ptr(style.Fg(style.Grey600)),
				},
			}),

			VSpacer(1),
			Divider(DividerProps{
				BaseProps: BaseProps{LayoutProps: LayoutProps{Width: 50}},
				Label:     "Progress Examples",
			}),
			VSpacer(1),

			// Progress bars
			VStack(ContainerProps{
				ChildrenProps: Children(
					Text(TextProps{Content: "Progress (25%):"}),
					ProgressBar(ProgressBarProps{
						BaseProps: BaseProps{LayoutProps: LayoutProps{Width: 40}},
						Value:     0.25,
						Variant:   "bar",
					}),

					VSpacer(1),
					Text(TextProps{Content: "Progress (75%):"}),
					ProgressBar(ProgressBarProps{
						BaseProps: BaseProps{LayoutProps: LayoutProps{Width: 40}},
						Value:     0.75,
						Variant:   "blocks",
					}),
				),
			}, 0),
		),
	}, 1)
}

func renderLayoutDemo(ctx *component.Context) component.Component {
	return VStack(ContainerProps{
		ChildrenProps: Children(
			Heading(TextProps{Content: "Layout Examples"}),
			VSpacer(1),

			// Boxes
			HStack(ContainerProps{
				ChildrenProps: Children(
					Card(ContainerProps{
						LayoutProps: LayoutProps{Width: 20, Padding: 1},
						ChildrenProps: Children(
							Text(TextProps{Content: "Card 1"}),
						),
					}),

					Card(ContainerProps{
						LayoutProps: LayoutProps{Width: 20, Padding: 1},
						ChildrenProps: Children(
							Text(TextProps{Content: "Card 2"}),
						),
					}),

					Card(ContainerProps{
						LayoutProps: LayoutProps{Width: 20, Padding: 1},
						ChildrenProps: Children(
							Text(TextProps{Content: "Card 3"}),
						),
					}),
				),
			}, 2),

			VSpacer(1),
			Divider(DividerProps{
				BaseProps: BaseProps{LayoutProps: LayoutProps{Width: 50}},
				Label:     "Alert Examples",
			}),
			VSpacer(1),

			// Alerts
			VStack(ContainerProps{
				ChildrenProps: Children(
					Alert(AlertProps{
						Type:    BadgeInfo,
						Message: "This is an info message",
						BaseProps: BaseProps{
							LayoutProps: LayoutProps{Width: 50},
						},
					}),
					Alert(AlertProps{
						Type:    BadgeSuccess,
						Message: "Operation completed successfully!",
						BaseProps: BaseProps{
							LayoutProps: LayoutProps{Width: 50},
						},
					}),
					Alert(AlertProps{
						Type:    BadgeWarning,
						Message: "Warning: Please review your input",
						BaseProps: BaseProps{
							LayoutProps: LayoutProps{Width: 50},
						},
					}),
				),
			}, 1),
		),
	}, 1)
}

func renderModalDemo(ctx *component.Context, modalOpen *component.State[bool]) component.Component {
	return VStack(ContainerProps{
		ChildrenProps: Children(
			Heading(TextProps{Content: "Modal Example"}),
			VSpacer(1),

			Text(TextProps{
				Content: "Click the button below to open a modal dialog.",
			}),

			VSpacer(1),

			Button(InteractiveProps{
				FocusProps: FocusProps{
					OnClick: func() {
						modalOpen.Set(true)
					},
				},
			}, "[Open Modal]"),

			VSpacer(2),

			Text(TextProps{
				Content: "Modal Status: " + func() string {
					if modalOpen.Get() {
						return "OPEN"
					} else {
						return "CLOSED"
					}
				}(),
				StyleProps: StyleProps{Style: Ptr(style.Fg(style.Grey500))},
			}),
		),
	}, 1)
}

func renderFooter() component.Component {
	footerStyle := style.Merge(
		style.Fg(style.Grey500),
		style.Italic,
	)

	return Text(TextProps{
		Content: "Navigation: ↑↓←→ arrows | Activate: Enter/Space | Quit: Ctrl+C",
		StyleProps: StyleProps{
			Style: Ptr(footerStyle),
		},
	})
}
