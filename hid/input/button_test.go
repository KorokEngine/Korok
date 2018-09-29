package input

import "testing"

func TestButton(t *testing.T) {
	btn := button{}
	btn.Update(true)

	if !btn.JustPressed() {
		t.Error("button not pressed")
	}
	if !btn.Down() {
		t.Error("button not down")
	}
	btn.Reset()

	btn.Update(true)
	if btn.JustPressed() {
		t.Error("button should not be pressed")
	}
	if !btn.Down() {
		t.Error("button not down")
	}
	btn.Reset()

	btn.Update(false)
	if !btn.JustReleased() {
		t.Error("button not released")
	}
	if btn.JustPressed() {
		t.Error("button should not be pressed")
	}
	if btn.Down() {
		t.Error("button should not be down")
	}
	btn.Reset()
}

func TestButton_Update2(t *testing.T) {
	btn := button{}
	btn.Update(true)
	btn.Update(true)

	if !btn.JustPressed() {
		t.Error("button not pressed")
	}
	if !btn.Down() {
		t.Error("button not down")
	}
	btn.Reset()

	btn.Update(false)
	btn.Update(false)

	if !btn.JustReleased() {
		t.Error("button not released")
	}
	if btn.Down() {
		t.Errorf("button should not be down")
	}
	btn.Reset()

	btn.Update(true)
	btn.Update(false)

	if !btn.JustPressed() {
		t.Error("button not pressed")
	}
	if !btn.JustReleased() {
		t.Error("button not released")
	}
	btn.Reset()
}