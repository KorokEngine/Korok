package gui

import (
	"unicode/utf8"
	"korok.io/korok/math"
	"korok.io/korok/math/f32"
	"korok.io/korok/gfx/font"
)

// 工具结构，负责把字符串转化为顶点..
// 拥有所有需要的条件属性
type FontRender struct {
	*DrawList
	fontSize float32
	font font.Font
	color uint32
	lineSpace float32
}

// 当前的实现中，不考虑裁切优化，全部绘制所有字符
func (fr *FontRender) RenderText(pos f32.Vec2, text string) (size f32.Vec2){
	dx, dy := pos[0], pos[1]
	maxWidth := float32(0)

	idxCount := len(text) * 6
	vtxCount := len(text) * 4
	fr.DrawList.PrimReserve(idxCount, vtxCount)

	vtxWriter := fr.DrawList.VtxWriter
	idxWriter := fr.DrawList.IdxWriter

	color := fr.color
	_, gh := fr.font.Bounds()
	scale := fr.fontSize/gh
	lineHeight := fr.fontSize

	lineSpace := fr.lineSpace
	if lineSpace == 0 {
		lineSpace = 0.4 * lineHeight
	}

	bufferUsed := 0

	for i, w := 0, 0; i < len(text); i += w {
		r, width := utf8.DecodeRuneInString(text[i:])
		w = width

		if r == '\n' {
			if dx > maxWidth {
				maxWidth = dx
			}

			dy -= lineHeight + lineSpace
			dx = pos[0]
			continue
		}

		g, _ := fr.font.Glyph(r)

		// Add kerning todo
		// dx += getKerning(preglyph, g)
		x1, y1 := dx, dy - (g.Height+g.YOffset) * scale
		x2, y2 := x1 + g.Width * scale, dy - (g.YOffset) * scale
		u1,v1, u2, v2 := fr.font.Frame(r)

		vi := bufferUsed * 4
		vtxWriter[vi+0] = DrawVert{f32.Vec2{x1, y1}, f32.Vec2{u1, v2}, color}
		vtxWriter[vi+1] = DrawVert{f32.Vec2{x2, y1}, f32.Vec2{u2, v2}, color}
		vtxWriter[vi+2] = DrawVert{f32.Vec2{x2, y2}, f32.Vec2{u2, v1}, color}
		vtxWriter[vi+3] = DrawVert{f32.Vec2{x1, y2}, f32.Vec2{u1, v1}, color}

		ii, offset := bufferUsed * 6, fr.DrawList.vtxIndex
		idxWriter[ii+0] = DrawIdx(offset+0)
		idxWriter[ii+1] = DrawIdx(offset+1)
		idxWriter[ii+2] = DrawIdx(offset+2)
		idxWriter[ii+3] = DrawIdx(offset+0)
		idxWriter[ii+4] = DrawIdx(offset+2)
		idxWriter[ii+5] = DrawIdx(offset+3)

		fr.DrawList.idxIndex += 6
		fr.DrawList.vtxIndex += 4

		dx += float32(g.Advance) * scale
		bufferUsed ++
		// character space & line spacing todo
		// 处理 < 32 的控制符号!!
	}

	if dx > maxWidth {
		maxWidth = dx
	}
	size[0] = maxWidth - pos[0]
	size[1] = pos[1]-dy+fr.fontSize

	fr.DrawList.AddCommand(bufferUsed*6)
	return
}

func (fr *FontRender) RenderWrapped(pos f32.Vec2, text string, wrapWidth float32) (size f32.Vec2){
	wrap  := math.Max(wrapWidth, 0)
	// wrap text
	_, lines := font.Wrap(fr.font, text, wrap, fr.fontSize)
	//log.Println(lines)

	size = fr.RenderText(pos, lines)
	return
}









