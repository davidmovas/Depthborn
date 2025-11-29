package components

import (
	"github.com/davidmovas/Depthborn/internal/ui/component"
	. "github.com/davidmovas/Depthborn/internal/ui/component/primitive"
	"github.com/davidmovas/Depthborn/internal/ui/style"
)

func Example(ctx *component.Context) component.Component {
	confirmState := component.UseState(ctx, false)

	return Button(InteractiveProps{
		FocusProps: FocusProps{
			OnClick: func() {
				confirmState.Set(!confirmState.Get())
			},
		},
	}, style.Sprintf("Confirm: %v", style.Val(confirmState.Get())))
}
