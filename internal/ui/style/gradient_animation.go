package style

import (
	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/component/layout"
)

// AnimatedGradient creates animated gradient colors by rotating them
// ctx - component context (pass your component.Context)
// speed - frames per color shift (higher = slower)
//   - speed 60  = 1 shift/sec (very slow, smooth)
//   - speed 30  = 2 shifts/sec (slow)
//   - speed 15  = 4 shifts/sec (medium)
//   - speed 5   = 12 shifts/sec (fast)
//   - speed 1   = 60 shifts/sec (very fast)
//
// colors - gradient colors to rotate
func AnimatedGradient(ctx *component.Context, speed int, colors ...Color) []Color {
	frame := getAnimationFrame(ctx)
	return animatedGradientRotate(frame, speed, colors...)
}

// AnimatedGradientReverse creates animated gradient by rotating in reverse
func AnimatedGradientReverse(ctx *component.Context, speed int, colors ...Color) []Color {
	frame := getAnimationFrame(ctx)
	return animatedGradientRotateReverse(frame, speed, colors...)
}

// AnimatedGradientPingPong creates animated gradient with ping-pong effect
func AnimatedGradientPingPong(ctx *component.Context, speed int, colors ...Color) []Color {
	frame := getAnimationFrame(ctx)
	return animatedGradientPingPong(frame, speed, colors...)
}

// AnimatedGradientWave creates wave-like animation
func AnimatedGradientWave(ctx *component.Context, speed int, colors ...Color) []Color {
	frame := getAnimationFrame(ctx)
	return animatedGradientWave(frame, speed, colors...)
}

func animatedGradientRotate(frame int, speed int, colors ...Color) []Color {
	if len(colors) == 0 {
		return colors
	}
	if speed <= 0 {
		speed = 1
	}

	offset := (frame / speed) % len(colors)
	rotated := make([]Color, len(colors))
	for i := range colors {
		rotated[i] = colors[(i+offset)%len(colors)]
	}
	return rotated
}

func animatedGradientRotateReverse(frame int, speed int, colors ...Color) []Color {
	if len(colors) == 0 {
		return colors
	}
	if speed <= 0 {
		speed = 1
	}

	offset := (frame / speed) % len(colors)
	rotated := make([]Color, len(colors))
	for i := range colors {
		rotated[i] = colors[(len(colors)-offset+i)%len(colors)]
	}
	return rotated
}

func animatedGradientPingPong(frame int, speed int, colors ...Color) []Color {
	if len(colors) == 0 {
		return colors
	}
	if speed <= 0 {
		speed = 1
	}

	cycle := len(colors) * 2
	pos := (frame / speed) % cycle

	if pos >= len(colors) {
		pos = cycle - pos - 1
	}

	rotated := make([]Color, len(colors))
	for i := range colors {
		rotated[i] = colors[(i+pos)%len(colors)]
	}
	return rotated
}

func animatedGradientWave(frame int, speed int, colors ...Color) []Color {
	if len(colors) < 2 {
		return colors
	}
	if speed <= 0 {
		speed = 1
	}

	extended := make([]Color, len(colors)*2)
	for i := range extended {
		extended[i] = colors[i%len(colors)]
	}

	offset := (frame / speed) % len(colors)
	result := make([]Color, len(colors))
	for i := range result {
		result[i] = extended[(i+offset)%len(extended)]
	}
	return result
}

func getAnimationFrame(ctx *component.Context) int {
	/*frameState := component.UseState(ctx, 0)

	controller := layout.NewAnimationController(func() {
		frameState.Set(frameState.Value() + 1)
	})

	controllerState := component.UseState(ctx, controller)

	component.UseEffect(ctx, func() {
		controllerState.Value().Stop()
		fmt.Println("Animation controller stopped")
	}, []any{})

	return controllerState.Value().GetFrame()*/

	frameState := component.UseState(ctx, 0)

	animController := component.UseState(ctx, (*layout.AnimationController)(nil))

	// Initialize controller if needed
	if animController.Value() == nil {
		// Создаём контроллер с коллбеком который обновляет state
		controller := layout.NewAnimationController(func() {
			// Обновляем state чтобы вызвать ререндер
			frameState.Set(frameState.Value() + 1)
		})
		animController.Set(controller)
	}

	effect := func() {
		if animController.Value() != nil {
			animController.Value().Stop()
		}
	}

	// Clean up on unmount
	component.UseEffect(ctx, effect, []any{})

	// Get current frame from controller
	frame := 0
	if animController.Value() != nil {
		frame = animController.Value().GetFrame()
	}

	return frame

	/*
	 frameState := component.UseState(ctx, 0, "frame")

	    // Animation controller
	    animController := component.UseState(ctx, (*layout.AnimationController)(nil), "anim")

	    // Initialize controller if needed
	    if animController.Value() == nil {
	       // Создаём контроллер с коллбеком который обновляет state
	       controller := layout.NewAnimationController(func() {
	          // Обновляем state чтобы вызвать ререндер
	          frameState.Set(frameState.Value() + 1)
	       })
	       animController.Set(controller)
	    }

	    effect := func() {
	       if animController.Value() != nil {
	          animController.Value().Stop()
	       }
	    }

	    // Clean up on unmount
	    component.UseEffect(ctx, effect, []any{})

	    // Get current frame from controller
	    frame := 0
	    if animController.Value() != nil {
	       frame = animController.Value().GetFrame()
	    }
	*/
}
