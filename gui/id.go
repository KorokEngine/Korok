package gui

// 每个控件都需要手动传入一个 Id
type ID int

// 返回一个 Id
func GenID(name string) ID {
	return 0
}

type IdMap map[string]ID
