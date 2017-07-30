package inputs

type VButton struct {
	Name string
	Keys []Key
}

func (*VButton) JustPressed() bool {
	return false
}

func (*VButton) JustReleased() bool {
	return false
}

func (*VButton) State() int {
	return 0
}




