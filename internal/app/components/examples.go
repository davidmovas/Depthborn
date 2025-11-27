package components

import (
	"github.com/davidmovas/Depthborn/internal/ui/component"
	. "github.com/davidmovas/Depthborn/internal/ui/component/primitive"
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
