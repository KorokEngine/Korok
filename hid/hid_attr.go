package hid

type AttrTable struct {
	file string
	lang string
}

func (at *AttrTable) File(init func() string) string {
	if at.file == "" {
		at.file = init()
	}
	return at.file
}

func (at *AttrTable) Lang(init func() string) string {
	if at.lang == "" {
		at.lang = init()
	}
	return at.lang
}

var deviceAttr AttrTable
