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

	// Glyphs->AdvanceX in a directly way (more cache-friendly, for calcTextSize functions which are often bottleneck in large UI)
	IndexAdvanceX []float32
	FallbackAdvanceX float32
}

// 当前的实现中，不考虑裁切优化，全部绘制所有字符
func (fr *FontRender) RenderText1(pos f32.Vec2, text string) (size f32.Vec2){
	dx, dy := pos[0], pos[1]
	maxWidth := float32(0)

	idxCount := len(text) * 6
	vtxCount := len(text) * 4
	fr.DrawList.PrimReserve(idxCount, vtxCount)

	vtxWriter := fr.DrawList.VtxWriter
	idxWriter := fr.DrawList.IdxWriter

	color := fr.color
	_, tex := fr.font.Tex2D()
	texWidth, texHeight := tex.Width, tex.Height
	glyphA, _ := fr.font.Glyph('A')
	scale := fr.fontSize/float32(glyphA.Height)
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
		x1, y1 := dx, dy - float32(g.Height+g.YOffset) * scale
		x2, y2 := x1 + float32(g.Width) * scale, y1 + float32(g.Height) * scale
		u1, v1 := float32(g.X)/ texWidth, float32(g.Y)/ texHeight
		u2, v2 := float32(g.X+g.Width)/ texWidth, float32(g.Y+g.Height)/ texHeight

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
	_, lines := fr.wrap(text, wrap)
	//log.Println(lines)

	size = fr.RenderText1(pos, lines)
	return
}

// 计算分行后的字符串, 算法其实很简单：
// 遍历字符串，如果发现 '\n' 则换行
// 如果超出(>wrap) WrapWidth 则向前回溯到最后一个空格位置...
// 如果没有空格，则删除最后一个字符
func (fr *FontRender) wrap(text string, wrap float32) (n int, lines string) {
	size := len(text)
	line := make([]byte, 0, size)
	buff := make([]byte, 0, size*2)
	glyphA, _ := fr.font.Glyph('A')
	scale := fr.fontSize/float32(glyphA.Height)


	for i, w := 0, 0; i < size;{
		var lineSize float32
		var lastSpace = -1

		for j := 0 ; i < size; i, j = i+w, j+w {
			r, width := utf8.DecodeRuneInString(text[i:])
			w = width

			if r == '\n' {
				i, j = i+w, j+w
				goto NEW_LINE
			}

			line = append(line, text[i:i+w]...)
			if g, ok :=  fr.font.Glyph(r); ok {
				lineSize += float32(g.Advance) * scale
			} else {
				// todo fallback ?
			}
			if r == ' ' || r == '\t' {
				lastSpace = j
			}

			if lineSize > wrap {
				i, j = i+w, j+w
				break
			}
		}

		// reach the end
		if lineSize < wrap {
			buff = append(buff, line...)
			n += 1
			break
		}

		if len(line) == 0 {

		}

		// if has space, break! or remove last char to fit line-width
		if lastSpace > 0 {
			i -= len(line) - lastSpace

			//log.Println("trim last space i:", i, "lastSpace:", lastSpace, "lineSize:", len(line))
			line = line[:lastSpace]
		} else {
			i -= w
			line = line[:len(line)-w]
		}

		NEW_LINE:
			line = append(line, '\n')
			buff = append(buff, line...)
			line = line[:0]
			n += 1
	}
	return n, string(buff)
}

func (fr *FontRender) CalculateTextSize1(text string) f32.Vec2{
	glyphA, _ := fr.font.Glyph('A')
	scale := fr.fontSize/float32(glyphA.Height)
	size := f32.Vec2{0, fr.fontSize}

	for i, w := 0, 0; i < len(text); i += w {
		r, width := utf8.DecodeRuneInString(text[i:])
		w = width
		g, _ := fr.font.Glyph(r)
		if r >= 32 {
			size[0] += float32(g.Advance) * scale
		}
	}
	return size
}










