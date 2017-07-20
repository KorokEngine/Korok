package spine

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

type Atlas struct {
	Pages   []*AtlasPage
	Regions []*AtlasRegion
	loader  TextureLoader
}

func NewAtlas(r io.Reader, loader TextureLoader) (*Atlas, error) {
	if loader == nil {
		return nil, errors.New("spine: texture loader cannot be nil")
	}
	var atlas Atlas
	atlas.loader = loader

	scanner := bufio.NewScanner(r)
	var page *AtlasPage
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			page = nil
			continue
		}
		if page == nil {
			var tuple []string
			var terr error

			formatName, err := readValue("format", scanner)
			if err != nil {
				return nil, errors.New("spine: failed to parse format: " + err.Error())
			}
			format, err := formatFromName(formatName)
			if err != nil {
				return nil, errors.New("spine: failed to parse format: " + err.Error())
			}

			size, err := readTuple("size", scanner)
			if err != nil {
				return nil, errors.New("spine: failed to read size" + err.Error())
			}

			width,  _ := strconv.Atoi(size[0])
			height, _ := strconv.Atoi(size[1])

			tuple, terr = readTuple("filter", scanner)
			if terr != nil {
				return nil, errors.New("spine: failed to read page: " + err.Error())
			}

			minFilter, err := filterFromName(tuple[0])
			if err != nil {
				return nil, errors.New("spine: failed to parse min filter: " + err.Error())
			}
			magFilter, err := filterFromName(tuple[1])
			if err != nil {
				return nil, errors.New("spine: failed to parse mag filter: " + err.Error())
			}
			page = &AtlasPage{
				Name:      line,
				Format:    format,
				MinFilter: minFilter,
				MagFilter: magFilter,
				UWrap:     ClampToEdge,
				VWrap:     ClampToEdge,
				Width:	   width,
				Height:	   height,
			}

			direction, err := readValue("repeat", scanner)
			if err != nil {
				return nil, errors.New("spine: failed to parse direction: " + err.Error())
			}
			switch direction {
			case "x":
				page.UWrap = Repeat
			case "y":
				page.VWrap = Repeat
			case "xy":
				page.UWrap = Repeat
				page.VWrap = Repeat
			}

			if err := atlas.loader.Load(page); err != nil {
				return nil, errors.New("spine: failed to load texture: " + err.Error())
			}

			atlas.Pages = append(atlas.Pages, page)
		} else {
			var tuple []string
			var terr error
			rotateVal, err := readValue("rotate", scanner)
			if err != nil {
				return nil, errors.New("spine: failed to parse rotate: " + err.Error())
			}
			rotate, err := strconv.ParseBool(rotateVal)
			if err != nil {
				return nil, errors.New("spine: failed to parse rotate: " + err.Error())
			}
			tuple, terr = readTuple("xy", scanner)
			if terr != nil {
				return nil, errors.New("spine: failed to read x, y tuple: " + err.Error())
			}
			x, err := strconv.Atoi(tuple[0])
			if err != nil {
				return nil, errors.New("spine: failed to parse x: " + err.Error())
			}
			y, err := strconv.Atoi(tuple[1])
			if err != nil {
				return nil, errors.New("spine: failed to parse y: " + err.Error())
			}

			tuple, terr = readTuple("size", scanner)
			if terr != nil {
				return nil, errors.New("spine: failed to read width, height tuple: " + err.Error())
			}
			width, err := strconv.Atoi(tuple[0])
			if err != nil {
				return nil, errors.New("spine: failed to parse width: " + err.Error())
			}
			height, err := strconv.Atoi(tuple[1])
			if err != nil {
				return nil, errors.New("spine: failed to parse height: " + err.Error())
			}

			region := &AtlasRegion{
				Name:   line,
				Page:   page,
				Rotate: rotate,
				X:      x,
				Y:      y,
				U:      float32(x) / float32(page.Width),
				V:      float32(y) / float32(page.Height),
			}

			if region.Rotate {
				region.U2 = float32(x+height) / float32(page.Width)
				region.V2 = float32(y+width) / float32(page.Height)
			} else {
				region.U2 = float32(x+width) / float32(page.Width)
				region.V2 = float32(y+height) / float32(page.Height)
			}
			region.X, region.Y = x, y
			region.Width, region.Height = width, height
			if region.Width < 0 {
				region.Width = -region.Width
			}
			if region.Height < 0 {
				region.Height = -region.Height
			}

			tuple, terr = readTuple("", scanner)
			if terr != nil {
				return nil, errors.New("spine: failed to read tuple: " + err.Error())
			}
			if len(tuple) == 4 { // split is optional
				for idx, sval := range tuple {
					val, err := strconv.Atoi(sval)
					if err != nil {
						return nil, errors.New("spine: failed to read split: " + err.Error())
					}
					region.Splits[idx] = val
				}

				tuple, terr = readTuple("", scanner)
				if terr != nil {
					return nil, errors.New("spine: failed to read tuple: " + err.Error())
				}
				if len(tuple) == 4 { // pad is optional, but only present with splits
					for idx, sval := range tuple {
						val, err := strconv.Atoi(sval)
						if err != nil {
							return nil, errors.New("spine: failed to read split: " + err.Error())
						}
						region.Pads[idx] = val
					}

					tuple, terr = readTuple("orig", scanner)
					if terr != nil {
						return nil, errors.New("spine: failed to read tuple: " + err.Error())
					}
				}
			}

			origWidth, err := strconv.Atoi(tuple[0])
			if err != nil {
				return nil, errors.New("spine: failed to parse original width: " + err.Error())
			}
			origHeight, err := strconv.Atoi(tuple[1])
			if err != nil {
				return nil, errors.New("spine: failed to parse original height: " + err.Error())
			}

			region.OriginalWidth, region.OriginalHeight = origWidth, origHeight

			tuple, terr = readTuple("offset", scanner)
			if terr != nil {
				return nil, errors.New("spine: failed to read x offset, y offset tuple: " + err.Error())
			}
			offX, err := strconv.Atoi(tuple[0])
			if err != nil {
				return nil, errors.New("spine: failed to parse x offset: " + err.Error())
			}
			offY, err := strconv.Atoi(tuple[1])
			if err != nil {
				return nil, errors.New("spine: failed to parse y offset: " + err.Error())
			}
			region.OffsetX, region.OffsetY = float32(offX), float32(offY)

			regIdxVal, err := readValue("index", scanner)
			if err != nil {
				return nil, errors.New("spine: failed to parse region index: " + err.Error())
			}
			regIdx, err := strconv.Atoi(regIdxVal)
			if err != nil {
				return nil, errors.New("spine: failed to parse region index: " + err.Error())
			}
			region.Index = regIdx

			atlas.Regions = append(atlas.Regions, region)
		}
	}
	if scanner.Err() != nil {
		return nil, errors.New("spine: failed to read atlas file: " + scanner.Err().Error())
	}
	return &atlas, nil
}

func readValue(expectedKey string, s *bufio.Scanner) (string, error) {
	if err := nextLine(s); err != nil {
		return "", err
	}
	line := s.Text()
	keyVal := strings.SplitN(line, ":", 2)
	if len(keyVal) < 2 {
		return "", errors.New("spine: invalid key/value: " + line)
	}
	key := strings.TrimSpace(keyVal[0])
	if expectedKey != "" && expectedKey != key {
		return "", errors.New("spine: expected " + expectedKey + ", got " + key)
	}
	return strings.TrimSpace(keyVal[1]), nil
}

func nextLine(s *bufio.Scanner) error {
	if s.Scan() {
		return nil
	}
	if s.Err() != nil {
		return errors.New("spine: failed to read file: " + s.Err().Error())
	} else {
		return errors.New("spine: unexpected EOF encountered")
	}
}

/** Returns the number of tuple values read (2 or 4). */
func readTuple(expectedKey string, s *bufio.Scanner) ([]string, error) {
	val, err := readValue(expectedKey, s)
	if err != nil {
		return nil, err
	}
	tuple := strings.Split(val, ",")
	for idx, v := range tuple {
		tuple[idx] = strings.TrimSpace(v)
	}
	return tuple, nil
}

/** Returns the first region found with the specified name. This method uses string comparison to find the region, so the result
* should be cached rather than calling this method multiple times.
* @return The region, or null. */
func (a *Atlas) FindRegion(name string) *AtlasRegion {
	for _, region := range a.Regions {
		if region.Name == name {
			return region
		}
	}
	return nil
}

func (a *Atlas) Dispose() error {
	for _, page := range a.Pages {
		if err := a.loader.Unload(page); err != nil {
			return err
		}
	}
	return nil
}

func formatFromName(name string) (TextureFormat, error) {
	var f TextureFormat
	switch name {
	case "Alpha":
		f = Alpha
	case "Intensity":
		f = Intensity
	case "LuminanceAlpha":
		f = LuminanceAlpha
	case "RGB565":
		f = RGB565
	case "RGBA4444":
		f = RGBA4444
	case "RGB888":
		f = RGB888
	case "RGBA8888":
		f = RGBA8888
	default:
		return 0, errors.New("unknown texture format: " + name)
	}
	return f, nil
}

func filterFromName(name string) (TextureFilter, error) {
	var f TextureFilter
	switch name {
	case "Nearest":
		f = Nearest
	case "Linear":
		f = Linear
	case "MipMap":
		f = MipMap
	case "MipMapNearestNearest":
		f = MipMapNearestNearest
	case "MipMapLinearNearest":
		f = MipMapLinearNearest
	case "MipMapNearestLinear":
		f = MipMapNearestLinear
	case "MipMapLinearLinear":
		f = MipMapLinearLinear
	default:
		return 0, errors.New("unknown filter: " + name)
	}
	return f, nil
}

type TextureFormat int

const (
	Alpha TextureFormat = iota
	Intensity
	LuminanceAlpha
	RGB565
	RGBA4444
	RGB888
	RGBA8888
)

type TextureFilter int

const (
	Nearest TextureFilter = iota
	Linear
	MipMap
	MipMapNearestNearest
	MipMapLinearNearest
	MipMapNearestLinear
	MipMapLinearLinear
)

type TextureWrap int

const (
	MirroredRepeat TextureWrap = iota
	ClampToEdge
	Repeat
)

type AtlasPage struct {
	Name           string
	Format         TextureFormat
	MinFilter      TextureFilter
	MagFilter      TextureFilter
	UWrap          TextureWrap
	VWrap          TextureWrap
	RendererObject interface{}
	Width, Height  int
}

type AtlasRegion struct {
	Page                          *AtlasPage
	Name                          string
	X, Y, Width, Height           int
	U, V, U2, V2                  float32
	OffsetX, OffsetY              float32
	OriginalWidth, OriginalHeight int
	Index                         int
	Rotate                        bool
	Splits                        [4]int
	Pads                          [4]int
}

type TextureLoader interface {
	Load(page *AtlasPage) error
	Unload(page *AtlasPage) error
}
