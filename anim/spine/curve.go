package spine

type Curve struct {
	curves []float32
}

func NewCurve(frameCount int) *Curve {
	curve := new(Curve)
	curve.curves = make([]float32, (frameCount-1)*6)
	return curve
}

func (c *Curve) frameCount() int {
	return len(c.curves)/6 + 1
}

func (c *Curve) SetLinear(index int) {
	c.curves[index*6] = 0
}

func (c *Curve) SetStepped(index int) {
	c.curves[index*6] = -1
}

func (c *Curve) SetCurve(index int, cx1, cy1, cx2, cy2 float32) {
	subdiv_step := float32(1) / 10
	subdiv_step2 := subdiv_step * subdiv_step
	subdiv_step3 := subdiv_step2 * subdiv_step
	pre1 := 3 * subdiv_step
	pre2 := 3 * subdiv_step2
	pre4 := 6 * subdiv_step2
	pre5 := 6 * subdiv_step3
	tmp1x := -cx1*2 + cx2
	tmp1y := -cy1*2 + cy2
	tmp2x := (cx1-cx2)*3 + 1
	tmp2y := (cy1-cy2)*3 + 1
	i := index * 6
	curves := c.curves
	curves[i] = cx1*pre1 + tmp1x*pre2 + tmp2x*subdiv_step3
	curves[i+1] = cy1*pre1 + tmp1y*pre2 + tmp2y*subdiv_step3
	curves[i+2] = tmp1x*pre4 + tmp2x*pre5
	curves[i+3] = tmp1y*pre4 + tmp2y*pre5
	curves[i+4] = tmp2x * pre5
	curves[i+5] = tmp2y * pre5
}

func (c *Curve) CurvePercent(index int, percent float32) float32 {
	if percent < 0 {
		percent = 0
	} else if percent > 1 {
		percent = 1
	}

	curveIndex := index * 6
	curves := c.curves
	dfx := curves[curveIndex]
	if dfx == 0 {
		return percent
	}
	if dfx == -1 {
		return 0
	}

	dfy := curves[curveIndex+1]
	ddfx := curves[curveIndex+2]
	ddfy := curves[curveIndex+3]
	dddfx := curves[curveIndex+4]
	dddfy := curves[curveIndex+5]
	x := dfx
	y := dfy
	i := 8

	for {
		if x >= percent {
			lastX := x - dfx
			lastY := y - dfy
			return lastY + (y-lastY)*(percent-lastX)/(x-lastX)
		}
		if i == 0 {
			break
		}
		i -= 1
		dfx += ddfx
		dfy += ddfy
		ddfx += dddfx
		ddfy += dddfy
		x += dfx
		y += dfy
	}
	return y + (1-y)*(percent-x)/(1-x)
}
