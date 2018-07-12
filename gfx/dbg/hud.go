package dbg

import (
	"fmt"
)

func LogFPS(fps, drawCall int) {
	hud.fps = fps
	hud.drawCall = drawCall
}

func Hud(format string, args... interface{}) {
	hud.verbs = append(hud.verbs, fmt.Sprintf(format, args...))
}

func HudFunc(fn func() string) {
	hud.verbs = append(hud.verbs, fn())
}

// Game internal state.
type HudLog struct {
	verbs []string
	drawCall, fps int
}

func (hud *HudLog) draw() {
	var (
		x = gRender.view.x - gRender.view.w/2
		y = gRender.view.y - gRender.view.h/2
	)

	// draw fps & draw call
	if (d & FPS) != 0 {
		drawFps(x, y, hud.fps)
		drawDrawCall(x, y, hud.drawCall)
	}

	// draw string
	if (d & Stats) != 0 {
		d := float32(0)
		x += 10
		y += gRender.view.h - 20
		for _, str := range hud.verbs {
			gBuffer.String(x, y-d, str, .6)
			d += 10
		}
	}
}

func (hud *HudLog) reset() {
	hud.verbs = hud.verbs[:0]
}

func drawFps(x, y float32, fps int) {
	Color(0xFF000000)
	gBuffer.Rect(x+5, y+5, 50, 6)

	// format: RGBA
	Color(0xFF00FF00)

	w := float32(fps)/60 * 50
	gBuffer.Rect(x+5, y+5, w, 5)

	// format: RGBA
	Color(0xFF000000)

	gBuffer.String(x+5, y+10, fmt.Sprintf("%d fps", fps), .6)
}

func drawDrawCall(x,y float32, dc int) {
	gBuffer.String(x+5, y+25, fmt.Sprintf("DrawCall: %d", hud.drawCall), .6)
}