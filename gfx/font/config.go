package font

type TTFConfig interface {
	FontSize() int
	Runes() []rune
	Direction() Direction
}

// TTF config based on range
type rangeTTFConfig struct {
	Low, High uint32
	Size int
}

func (fc rangeTTFConfig) FontSize() int {
	return fc.Size
}

func (fc rangeTTFConfig) Runes() (runes []rune) {
	if fc.Low == fc.High && fc.Low == 0{
		fc.Low = 32
		fc.High = 127
	}
	runes = make([]rune, fc.High-fc.Low+1)
	for i,j := fc.Low, fc.High; i <= j; i++ {
		runes[i-fc.Low] = rune(i)
	}
	return
}

func (rangeTTFConfig) Direction() Direction {
	return LeftToRight
}

// TTF config based on runes collection
type ttfConfig struct {
	size int
	runes []rune
}

func (ttf ttfConfig) FontSize() int {
	return ttf.size
}

func (ttf ttfConfig) Runes() []rune {
	return ttf.runes
}

func (ttfConfig) Direction() Direction {
	return LeftToRight
}

func NewTTFConfig(size int, runes []rune) TTFConfig {
	return ttfConfig{size:size, runes:runes}
}

func ASCII(size int) TTFConfig {
	return rangeTTFConfig{Size:size, Low:32, High:127}
}










