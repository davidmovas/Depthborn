package components

import (
	"fmt"

	"github.com/davidmovas/Depthborn/internal/ui/component"
	. "github.com/davidmovas/Depthborn/internal/ui/component/primitive"
	"github.com/davidmovas/Depthborn/internal/ui/component/utils"
	"github.com/davidmovas/Depthborn/internal/ui/style"
)

func ExampleConfirmDialog(ctx *component.Context) component.Component {
	dialogOpen := component.UseState(ctx, false, "dialog_open")

	return VStack(ContainerProps{
		ChildrenProps: Children(
			Heading(TextProps{
				Content: "Confirm Dialog Example",
			}),

			Button(
				InteractiveProps{
					StyleProps: WithStyle(
						style.Merge(style.Fg(style.White), style.P(1), style.Br()),
					),
					FocusProps: FocusProps{
						OnClick: func() {
							dialogOpen.Set(true)
						},
					},
				},
				"Delete Item",
			),

			// Dialog
			Dialog(DialogProps{
				Open:    dialogOpen.Value(),
				Title:   "Confirm Delete",
				Message: "Are you sure you want to delete this item? This action cannot be undone.",
				OnConfirm: func() {
					dialogOpen.Set(false)
				},
				OnCancel: func() {
					dialogOpen.Set(false)
				},
				Variant: Ptr(BadgeError),
			}),
		),
	},
		2,
	)
}

func ExampleFormModal(ctx *component.Context) component.Component {
	modalOpen := component.UseState(ctx, false, "modal_open")
	username := component.UseState(ctx, "", "username")
	email := component.UseState(ctx, "", "email")

	return VStack(
		ContainerProps{
			ChildrenProps: Children(
				Heading(TextProps{
					Content: "Form Modal Example",
				}),

				Button(
					InteractiveProps{
						StyleProps: WithStyle(
							style.Merge(style.Bg(style.Primary), style.Fg(style.White), style.P(1), style.Br()),
						),
						FocusProps: FocusProps{
							OnClick: func() {
								modalOpen.Set(true)
							},
						},
					},
					"Open Form",
				),

				// Custom Modal
				Modal(ModalProps{
					Open: modalOpen.Value(),
					OnClose: func() {
						modalOpen.Set(false)
					},
					Title:       Ptr("Create New User"),
					Description: Ptr("Fill in the details to create a new user account."),
					Size:        ModalSizeMedium,
					ContainerProps: ContainerProps{
						ChildrenProps: ChildrenProps{
							Children: []component.Component{
								// Form fields
								VStack(ContainerProps{
									ChildrenProps: Children(
										VStack(ContainerProps{
											ChildrenProps: Children(
												Label(TextProps{Content: "Username"}),
												Box(ContainerProps{
													StyleProps: WithStyle(
														style.Merge(style.Bg(style.Grey100), style.P(1), style.Br()),
													),
													ChildrenProps: ChildrenProps{
														Children: []component.Component{
															Text(TextProps{
																Content: func() string {
																	if username.Value() == "" {
																		return "Enter username..."
																	}
																	return username.Value()
																}(),
															}),
														},
													},
												}),
											),
										}, 1),
									),
								}, 2),

								// Email field
								VStack(ContainerProps{
									ChildrenProps: Children(
										Label(TextProps{Content: "Email"}),
										Box(ContainerProps{
											StyleProps: WithStyle(
												style.Merge(style.Bg(style.Grey100), style.P(1), style.Br()),
											),
											ChildrenProps: ChildrenProps{
												Children: []component.Component{
													Text(TextProps{
														Content: func() string {
															if email.Value() == "" {
																return "Enter email..."
															}
															return email.Value()
														}(),
													}),
												},
											},
										}),
									),
								}, 1),
							},
						},
					},
					Footer: []component.Component{
						HStack(ContainerProps{
							ChildrenProps: Children(
								Button(
									InteractiveProps{
										StyleProps: WithStyle(
											style.Merge(style.Bg(style.Grey200), style.Fg(style.Grey700), style.P(1), style.Br()),
										),
										FocusProps: FocusProps{
											OnClick: func() {
												modalOpen.Set(false)
											},
										},
									},
									"Cancel",
								),

								Button(
									InteractiveProps{
										StyleProps: WithStyle(
											style.Merge(style.Bg(style.Primary), style.Fg(style.White), style.P(1), style.Br()),
										),
										FocusProps: FocusProps{
											OnClick: func() {
												// Submit form
												modalOpen.Set(false)
											},
											AutoFocus: Ptr(true),
										},
									},
									"Create",
								),
							),
						}, 2),
					},
					CloseOnEscape:  Ptr(true),
					CloseOnOverlay: Ptr(true),
					ShowCloseBtn:   Ptr(true),
				}),
			),
		},
		2,
	)
}

func ExampleModernFormModal(ctx *component.Context) component.Component {
	// State management
	modalOpen := component.UseState(ctx, false, "modal_open")
	username := component.UseState(ctx, "", "username")
	email := component.UseState(ctx, "", "email")
	bio := component.UseState(ctx, "", "bio")
	country := component.UseState(ctx, "", "country")

	// Validation errors
	usernameError := component.UseState(ctx, "", "username_error")
	emailError := component.UseState(ctx, "", "email_error")

	// Submit handler
	handleSubmit := func() {
		// Validate
		valid := true

		if username.Value() == "" {
			usernameError.Set("Username is required")
			valid = false
		} else {
			usernameError.Set("")
		}

		if email.Value() == "" {
			emailError.Set("Email is required")
			valid = false
		} else if utils.ValidateEmail(email.Value()) == nil {
			emailError.Set("")
		} else {
			emailError.Set("Invalid email format")
			valid = false
		}

		if valid {
			// Form is valid, submit
			fmt.Printf("Creating user: %s (%s)\n", username.Value(), email.Value())
			modalOpen.Set(false)

			// Reset form
			username.Set("")
			email.Set("")
			bio.Set("")
			country.Set("")
		}
	}

	return VStack(
		ContainerProps{
			ChildrenProps: Children(
				Heading(TextProps{
					Content: "Modern Form Modal Example",
				}),

				Text(TextProps{
					Content: "Click the button to open a modal with a fully functional form.",
					StyleProps: StyleProps{
						Style: Ptr(style.Fg(style.Grey600)),
					},
				}),

				// Open Modal Button
				Button(
					InteractiveProps{
						StyleProps: StyleProps{
							Style: Ptr(
								style.Merge(
									style.Bg(style.Primary),
									style.Fg(style.White),
									style.P(1),
									style.Br(),
								),
							),
						},
						FocusProps: FocusProps{
							OnClick: func() {
								modalOpen.Set(true)
							},
						},
					},
					"🚀 Create New User",
				),

				// The Modal
				Modal(ModalProps{
					Open: modalOpen.Value(),
					OnClose: func() {
						modalOpen.Set(false)
						// Reset errors on close
						usernameError.Set("")
						emailError.Set("")
					},
					Title:       Ptr("Create New User"),
					Description: Ptr("Fill in the details to create a new user account."),
					Size:        ModalSizeMedium,

					ContainerProps: ContainerProps{
						ChildrenProps: ChildrenProps{
							Children: []component.Component{
								// Form content
								VStack(
									ContainerProps{
										ChildrenProps: Children(
											// Username Input
											VStack(
												ContainerProps{
													ChildrenProps: Children(
														Label(TextProps{
															Content: "Username *",
														}),
														Input(InputProps{
															BaseProps: BaseProps{
																LayoutProps: LayoutProps{
																	Width: Ptr(40),
																},
															},
															Value:       username.Value(),
															Placeholder: "Enter username...",
															Type:        "text",
															Prefix:      "👤 ",
															ErrorText:   usernameError.Value(),
															OnChange: func(val string) {
																username.Set(val)
																if val != "" {
																	usernameError.Set("")
																}
															},
														}),
													),
												},
												0,
											),

											// Email Input
											VStack(
												ContainerProps{
													ChildrenProps: Children(
														Label(TextProps{
															Content: "Email Address *",
														}),
														Input(InputProps{
															BaseProps: BaseProps{
																LayoutProps: LayoutProps{
																	Width: Ptr(40),
																},
															},
															Value:       email.Value(),
															Placeholder: "your@email.com",
															Type:        "text",
															Prefix:      "📧 ",
															ErrorText:   emailError.Value(),
															OnChange: func(val string) {
																email.Set(val)
																if val != "" {
																	emailError.Set("")
																}
															},
														}),
													),
												},
												0,
											),

											// Country Select
											VStack(
												ContainerProps{
													ChildrenProps: Children(
														Label(TextProps{
															Content: "Country",
														}),
														Select(SelectProps{
															BaseProps: BaseProps{
																LayoutProps: LayoutProps{
																	Width: Ptr(40),
																},
															},
															Options: []SelectOption{
																{Label: "🇺🇸 United States", Value: "us"},
																{Label: "🇬🇧 United Kingdom", Value: "uk"},
																{Label: "🇩🇪 Germany", Value: "de"},
																{Label: "🇫🇷 France", Value: "fr"},
																{Label: "🇯🇵 Japan", Value: "jp"},
																{Label: "🇨🇦 Canada", Value: "ca"},
															},
															Value:       country.Value(),
															Placeholder: "Select your country...",
															OnChange: func(val string) {
																country.Set(val)
															},
														}),
													),
												},
												0,
											),

											// Bio TextArea
											VStack(
												ContainerProps{
													ChildrenProps: Children(
														Label(TextProps{
															Content: "Bio (optional)",
														}),
														TextArea(TextAreaProps{
															BaseProps: BaseProps{
																LayoutProps: LayoutProps{
																	Width:  Ptr(40),
																	Height: Ptr(4),
																},
															},
															Value:       bio.Value(),
															Placeholder: "Tell us about yourself...",
															OnChange: func(val string) {
																bio.Set(val)
															},
														}),
													),
												},
												0,
											),
										),
									},
									2, // Gap between fields
								),
							},
						},
					},

					// Modal Footer with buttons
					Footer: []component.Component{
						HStack(
							ContainerProps{
								ChildrenProps: Children(
									// Cancel Button
									Button(
										InteractiveProps{
											StyleProps: StyleProps{
												Style: Ptr(
													style.Merge(
														style.Bg(style.Grey200),
														style.Fg(style.Grey700),
														style.P(1),
														style.Br(),
													),
												),
											},
											FocusProps: FocusProps{
												OnClick: func() {
													modalOpen.Set(false)
													usernameError.Set("")
													emailError.Set("")
												},
											},
										},
										"Cancel",
									),

									// Submit Button
									Button(
										InteractiveProps{
											StyleProps: StyleProps{
												Style: Ptr(
													style.Merge(
														style.Bg(style.Primary),
														style.Fg(style.White),
														style.P(1),
														style.Br(),
													),
												),
											},
											FocusProps: FocusProps{
												OnClick:   handleSubmit,
												AutoFocus: Ptr(true),
											},
										},
										"✓ Create User",
									),
								),
							},
							2,
						),
					},

					CloseOnEscape:  Ptr(true),
					CloseOnOverlay: Ptr(true),
					ShowCloseBtn:   Ptr(true),
				}),
			),
		},
		2,
	)
}

func ExampleProgressModal(ctx *component.Context) component.Component {
	modalOpen := component.UseState(ctx, false, "progress_modal_open")
	progress := component.UseState(ctx, 0.0, "progress")
	frame := component.UseState(ctx, 0, "anim_frame")

	component.UseEffect(ctx, func() {
		if modalOpen.Value() && progress.Value() < 1.0 {
			newProgress := progress.Value() + 0.001
			if newProgress > 1.0 {
				newProgress = 1.0
			}
			progress.Set(newProgress)
			frame.Set(frame.Value() + 1)
		}
	}, []any{modalOpen.Value(), progress.Value()}, "progress_effect")

	return VStack(
		ContainerProps{
			ChildrenProps: Children(
				Heading(TextProps{
					Content: "Progress Modal Example",
				}),

				Button(
					InteractiveProps{
						StyleProps: StyleProps{
							Style: Ptr(
								style.Merge(
									style.Bg(style.Info),
									style.Fg(style.White),
									style.P(1),
									style.Br(),
								),
							),
						},
						FocusProps: FocusProps{
							OnClick: func() {
								modalOpen.Set(true)
								progress.Set(0.0)
							},
						},
					},
					"⬇️ Start Download",
				),

				Modal(ModalProps{
					Open:        modalOpen.Value(),
					OnClose:     func() { modalOpen.Set(false) },
					Title:       Ptr("Downloading..."),
					Description: Ptr("Please wait while we download your files."),
					Size:        ModalSizeSmall,

					ContainerProps: ContainerProps{
						ChildrenProps: ChildrenProps{
							Children: []component.Component{
								VStack(
									ContainerProps{
										ChildrenProps: Children(
											// Spinner
											Spinner(SpinnerProps{
												Variant: "dots",
												Frame:   frame.Value(),
												Label:   fmt.Sprintf("Downloading... %.0f%%", progress.Value()*100),
											}),

											// Progress Bar
											ProgressBar(ProgressBarProps{
												BaseProps: BaseProps{
													LayoutProps: LayoutProps{
														Width: Ptr(50),
													},
												},
												Value:          progress.Value(),
												Variant:        "bar",
												ShowPercentage: true,
											}),

											// File count
											Text(TextProps{
												Content: fmt.Sprintf("%d / 100 files downloaded", int(progress.Value()*100)),
												StyleProps: StyleProps{
													Style: Ptr(style.Fg(style.Grey600)),
												},
											}),
										),
									},
									2,
								),
							},
						},
					},

					Footer: []component.Component{
						Button(
							InteractiveProps{
								StyleProps: StyleProps{
									Style: Ptr(
										style.Merge(
											style.Bg(style.Error),
											style.Fg(style.White),
											style.P(1),
											style.Br(),
										),
									),
								},
								FocusProps: FocusProps{
									OnClick: func() {
										modalOpen.Set(false)
										progress.Set(0.0)
									},
								},
							},
							"Cancel Download",
						),
					},

					CloseOnEscape:  Ptr(false),
					CloseOnOverlay: Ptr(false),
					ShowCloseBtn:   Ptr(false),
				}),
			),
		},
		2,
	)
}
