package spine

import (
	"math"
)

type Timeline interface {
	Apply(skeleton *Skeleton, time, alpha float32)
}

type RotateTimeline struct {
	boneIndex int
	frames    []float32
	curve     *Curve
}

func NewRotateTimeline(l int) *RotateTimeline {
	timeline := new(RotateTimeline)
	timeline.frames = make([]float32, l*2)
	timeline.curve = NewCurve(l)
	return timeline
}

func (t *RotateTimeline) Apply(skeleton *Skeleton, time, alpha float32) {
	frames := t.frames
	if time < frames[0] {
		return
	}

	bone := skeleton.Bones[t.boneIndex]
	if time >= frames[len(frames)-2] {
		amount := bone.Data.rotation + frames[len(frames)-1] - bone.Rotation
		for amount > 180 {
			amount -= 360
		}

		for amount < -180 {
			amount += 360
		}
		bone.Rotation += amount * alpha
		return
	}

	frameIndex := binarySearch(frames, time, 2)
	lastFrameValue := frames[frameIndex-1]
	frameTime := frames[frameIndex]
	percent := 1 - (time-frameTime)/(frames[frameIndex-2]-frameTime)
	percent = t.curve.CurvePercent(frameIndex/2-1, percent)
	amount := frames[frameIndex+1] - lastFrameValue
	for amount > 180 {
		amount -= 360
	}
	for amount < -180 {
		amount += 360
	}
	amount = bone.Data.rotation + (lastFrameValue + amount*percent) - bone.Rotation
	for amount > 180 {
		amount -= 360
	}

	for amount < -180 {
		amount += 360
	}
	bone.Rotation += amount * alpha
}

func binarySearch(values []float32, target float32, step int) int {
	low := 0
	high := int(math.Floor(float64(len(values)/step))) - 2
	if high == 0 {
		return step
	}
	current := high >> 1
	for {
		if values[(current+1)*step] <= target {
			low = current + 1
		} else {
			high = current
		}
		if low == high {
			return (low + 1) * step
		}
		current = (low + high) >> 1
	}
}

func (t *RotateTimeline) setFrame(index int, time, angle float32) {
	frameIndex := index * 2
	t.frames[frameIndex] = time
	t.frames[frameIndex+1] = angle
}

func (t *RotateTimeline) frameCount() int {
	return len(t.frames) / 2
}

type TranslateTimeline struct {
	boneIndex int
	frames    []float32
	curve     *Curve
}

func NewTranslateTimeline(l int) *TranslateTimeline {
	timeline := new(TranslateTimeline)
	timeline.frames = make([]float32, l*3)
	timeline.curve = NewCurve(l)
	return timeline
}

func (t *TranslateTimeline) frameCount() int {
	return len(t.frames) / 3
}

func (t *TranslateTimeline) setFrame(index int, time, x, y float32) {
	frameIndex := index * 3
	t.frames[frameIndex] = time
	t.frames[frameIndex+1] = x
	t.frames[frameIndex+2] = y
}

func (t *TranslateTimeline) Apply(skeleton *Skeleton, time, alpha float32) {
	frames := t.frames
	if time < frames[0] {
		return
	}

	bone := skeleton.Bones[t.boneIndex]

	if time >= frames[len(frames)-3] {
		bone.X += (bone.Data.x + frames[len(frames)-2] - bone.X) * alpha
		bone.Y += (bone.Data.y + frames[len(frames)-1] - bone.Y) * alpha
		return
	}

	frameIndex := binarySearch(frames, time, 3)
	lastFrameX := frames[frameIndex-2]
	lastFrameY := frames[frameIndex-1]
	frameTime := frames[frameIndex]
	percent := 1 - (time-frameTime)/(frames[frameIndex-3]-frameTime)
	percent = t.curve.CurvePercent(frameIndex/3-1, percent)

	bone.X += (bone.Data.x + lastFrameX + (frames[frameIndex+1]-lastFrameX)*percent - bone.X) * alpha
	bone.Y += (bone.Data.y + lastFrameY + (frames[frameIndex+2]-lastFrameY)*percent - bone.Y) * alpha
}

type ScaleTimeline struct {
	boneIndex int
	frames    []float32
	curve     *Curve
}

func NewScaleTimeline(l int) *ScaleTimeline {
	timeline := new(ScaleTimeline)
	timeline.frames = make([]float32, l*3)
	timeline.curve = NewCurve(l)
	return timeline
}

func (t *ScaleTimeline) frameCount() int {
	return len(t.frames) / 3
}

func (t *ScaleTimeline) setFrame(index int, time, x, y float32) {
	frameIndex := index * 3
	t.frames[frameIndex] = time
	t.frames[frameIndex+1] = x
	t.frames[frameIndex+2] = y
}

func (t *ScaleTimeline) Apply(skeleton *Skeleton, time, alpha float32) {
	frames := t.frames
	if time < frames[0] {
		return
	}

	bone := skeleton.Bones[t.boneIndex]

	if time >= frames[len(frames)-3] {
		bone.ScaleX += (bone.Data.scaleX - 1 + frames[len(frames)-2] - bone.ScaleX) * alpha
		bone.ScaleY += (bone.Data.scaleY - 1 + frames[len(frames)-1] - bone.ScaleY) * alpha
		return
	}

	frameIndex := binarySearch(frames, time, 3)
	lastFrameX := frames[frameIndex-2]
	lastFrameY := frames[frameIndex-1]
	frameTime := frames[frameIndex]
	percent := 1 - (time-frameTime)/(frames[frameIndex-3]-frameTime)
	percent = t.curve.CurvePercent(frameIndex/3-1, percent)

	bone.ScaleX += (bone.Data.scaleX - 1 + lastFrameX + (frames[frameIndex+1]-lastFrameX)*percent - bone.ScaleX) * alpha
	bone.ScaleY += (bone.Data.scaleY - 1 + lastFrameY + (frames[frameIndex+2]-lastFrameY)*percent - bone.ScaleY) * alpha
}

type ColorTimeline struct {
	slotIndex int
	frames    []float32
	curve     *Curve
}

func NewColorTimeline(l int) *ColorTimeline {
	return &ColorTimeline{
		frames: make([]float32, l*5),
		curve:  NewCurve(l),
	}
}

func (t *ColorTimeline) frameCount() int {
	return t.curve.frameCount()
}

func (t *ColorTimeline) setFrame(index int, time, r, g, b, a float32) {
	index *= 5
	frames := t.frames
	frames[index] = time
	frames[index+1] = r
	frames[index+2] = g
	frames[index+3] = b
	frames[index+4] = a
}

func (t *ColorTimeline) Apply(skeleton *Skeleton, time, alpha float32) {
	frames := t.frames
	if time < frames[0] {
		return // Time is before first frame.
	}

	slot := skeleton.Slots[t.slotIndex]

	if time >= frames[len(t.frames)-5] { // Time is after last frame.
		i := len(frames) - 1
		slot.R = frames[i-3]
		slot.G = frames[i-2]
		slot.B = frames[i-1]
		slot.A = frames[i]
		return
	}

	// Interpolate between the last frame and the current frame.
	frameIndex := binarySearch(frames, time, 5)
	lastFrameR := frames[frameIndex-4]
	lastFrameG := frames[frameIndex-3]
	lastFrameB := frames[frameIndex-2]
	lastFrameA := frames[frameIndex-1]
	frameTime := frames[frameIndex]
	percent := 1 - (time-frameTime)/(frames[frameIndex-5]-frameTime)
	percent = t.curve.CurvePercent(frameIndex/5-1, percent)

	r := lastFrameR + (frames[frameIndex+1]-lastFrameR)*percent
	g := lastFrameG + (frames[frameIndex+2]-lastFrameG)*percent
	b := lastFrameB + (frames[frameIndex+3]-lastFrameB)*percent
	a := lastFrameA + (frames[frameIndex+4]-lastFrameA)*percent
	if alpha < 1 {
		slot.R += (r - slot.R) * alpha
		slot.G += (g - slot.G) * alpha
		slot.B += (b - slot.B) * alpha
		slot.A += (a - slot.A) * alpha
	} else {
		slot.R = r
		slot.G = g
		slot.B = b
		slot.A = a
	}
}

type AttachmentTimeline struct {
	slotIndex       int
	frames          []float32
	attachmentNames []string
}

func NewAttachmentTimeline(l int) *AttachmentTimeline {
	return &AttachmentTimeline{
		frames:          make([]float32, l),
		attachmentNames: make([]string, l),
	}
}

func (t *AttachmentTimeline) frameCount() int {
	return len(t.frames)
}

func (t *AttachmentTimeline) setFrame(index int, time float32, attachmentName string) {
	t.frames[index] = time
	t.attachmentNames[index] = attachmentName
}

func (t *AttachmentTimeline) Apply(skeleton *Skeleton, time, alpha float32) {
	frames := t.frames
	if time < frames[0] {
		return // Time is before first frame.
	}

	var frameIndex int
	if time >= frames[len(frames)-1] { // Time is after last frame.
		frameIndex = len(frames) - 1
	} else {
		frameIndex = binarySearch(frames, time, 1) - 1
	}

	attachmentName := t.attachmentNames[frameIndex]
	var attachment Attachment
	if attachmentName != "" {
		attachment = skeleton.AttachmentBySlotIndex(t.slotIndex, attachmentName)
	}
	skeleton.Slots[t.slotIndex].Attachment = attachment
}

type Animation struct {
	name      string
	timelines []Timeline
	duration  float32
}

func NewAnimation(name string, timelines []Timeline, duration float32) *Animation {
	anim := new(Animation)
	anim.name = name
	anim.timelines = timelines
	anim.duration = duration
	return anim
}

func (a *Animation) Apply(skeleton *Skeleton, time float32, loop bool) {
	if loop && a.duration != 0 {
		time = float32(math.Mod(float64(time), float64(a.duration)))
	}

	for _, timeline := range a.timelines {
		timeline.Apply(skeleton, time, 1)
	}
}

func (a *Animation) Mix(skeleton *Skeleton, time float32, loop bool, alpha float32) {
	if loop && a.duration != 0 {
		time = float32(math.Mod(float64(time), float64(a.duration)))
	}
	for _, timeline := range a.timelines {
		timeline.Apply(skeleton, time, alpha)
	}
}

func (a *Animation) Duration() float32 {
	return a.duration
}
