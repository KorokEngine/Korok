package gui

import (
	"github.com/go-gl/mathgl/mgl32"
	"unicode/utf8"
	"korok.io/korok/gfx"
	"korok.io/korok/engi/math"
	//"log"
)

// 工具结构，负责把字符串转化为顶点..
// 拥有所有需要的条件属性
type FontRender struct {
	*DrawList
	fontSize float32
	font gfx.FontSystem
	color uint32

	// Glyphs->AdvanceX in a directly way (more cache-friendly, for calcTextSize functions which are often bottleneck in large UI)
	IndexAdvanceX []float32
	FallbackAdvanceX float32
}

// 当前的实现中，不考虑裁切优化，全部绘制所有字符
func (fr *FontRender) RenderText1(pos mgl32.Vec2, text string) {
	dx, dy := pos[0], pos[1]
	maxWidth := float32(0)

	idxCount := len(text) * 6
	vtxCount := len(text) * 4
	fr.DrawList.PrimReserve(idxCount, vtxCount)

	vtxWriter := fr.DrawList.VtxWriter
	idxWriter := fr.DrawList.IdxWriter

	color := fr.color
	_, tex := fr.font.Tex()
	texWidth, texHeight := tex.Width, tex.Height
	glyphA := fr.font.Glyph('A')
	scale := fr.fontSize/float32(glyphA.Height)
	lineHeight := fr.fontSize

	bufferUsed := 0

	for i, w := 0, 0; i < len(text); i += w {
		r, width := utf8.DecodeRuneInString(text[i:])
		w = width

		if r == '\n' {
			if dx > maxWidth {
				maxWidth = dx
			}

			dy -= lineHeight
			dx = pos[0]
			continue
		}

		g := fr.font.Glyph(r)

		// Add kerning todo
		// dx += getKerning(preglyph, g)

		x1, y1 := dx, dy
		x2, y2 := dx + float32(g.Width) * scale, dy + float32(g.Height) * scale
		u1, v1 := float32(g.X)/ texWidth, float32(g.Y)/ texHeight
		u2, v2 := float32(g.X+g.Width)/ texWidth, float32(g.Y+g.Height)/ texHeight

		vi := bufferUsed * 4
		vtxWriter[vi+0] = DrawVert{mgl32.Vec2{x1, y1}, mgl32.Vec2{u1, v2}, color}
		vtxWriter[vi+1] = DrawVert{mgl32.Vec2{x2, y1}, mgl32.Vec2{u2, v2}, color}
		vtxWriter[vi+2] = DrawVert{mgl32.Vec2{x2, y2}, mgl32.Vec2{u2, v1}, color}
		vtxWriter[vi+3] = DrawVert{mgl32.Vec2{x1, y2}, mgl32.Vec2{u1, v1}, color}

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

	if bufferUsed < len(text) {
		// ..
	}
}

func (fr *FontRender) RenderWrapped(pos mgl32.Vec2, text string, wrapWidth float32) {
	wrap  := math.Max(wrapWidth, 0)
	// wrap text 
	_, lines := fr.wrap(text, wrap)
	fr.RenderText1(pos, lines)
}

// 计算分行后的字符串, 算法其实很简单：
// 遍历字符串，如果发现 '\n' 则换行
// 如果超出(>wrap) WrapWidth 则向前回溯到最后一个空格位置...
// 如果没有空格，则删除最后一个字符
func (fr *FontRender) wrap(text string, wrap float32) (n int, lines string) {
	size := len(text)
	line := make([]byte, 0, size)
	buff := make([]byte, 0, size*2)
	glyphA := fr.font.Glyph('A')
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
			lineSize += float32(fr.font.Glyph(r).Advance) * scale
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












