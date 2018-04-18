package game

import (
	"korok.io/korok/gfx/dbg"
	"fmt"
)

// Game internal state.
type Stats struct {
	verbs []string
	drawCall int
}

func (st *Stats) V(str string) {
	st.verbs = append(st.verbs, str)
}

// screen height
func (st *Stats) printVerb() {
	for _, str := range st.verbs {
		dbg.DrawStrScaled(str, .6)
		dbg.Return()
	}
}

func (st *Stats) printDrawCall() {
	dbg.Move(5, 25)
	dbg.DrawStrScaled(fmt.Sprintf("drawCall: %d", st.drawCall), .6)
}

func (st *Stats) printFPS(fps int32) {
	dbg.Move(5, 5)

	dbg.Color(0xFF000000)
	dbg.DrawRect(0, 0, 50, 6)

	// format: RGBA
	dbg.Color(0xFF00FF00)

	w := float32(fps)/60 * 50
	dbg.DrawRect(0, 0, w, 5)

	// format: RGBA
	dbg.Color(0xFF000000)

	dbg.Move(5, 10)
	dbg.DrawStrScaled(fmt.Sprintf("%d fps", fps), 0.6)
}
