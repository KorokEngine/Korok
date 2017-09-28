package bx

type HandleAlloc struct {
	handles [10]uint16
	_map 	[10]uint16
}

func NewHandleAlloc(cap uint32) *HandleAlloc {
	//ha := new(HandleAlloc)
	//ha.handles =
	return nil
}

func (ha *HandleAlloc) GetHandles() uint16 {
	return 0
}

func (ha *HandleAlloc) GetHandleAt(at uint16) uint16 {
	return 0
}

func (ha *HandleAlloc) GetNumHandles() uint16 {
	return 0
}

func (ha *HandleAlloc) GetMaxHandles() uint16 {
	return 0
}

func (ha *HandleAlloc) Alloc() uint16 {
	return 0
}

func (ha *HandleAlloc) IsValid(handle uint16) bool {
	return false
}

func (ha *HandleAlloc) Free(handle uint16) {

}

func (ha *HandleAlloc) Reset() {

}


type HandleHashMap struct {

}

func (m *HandleHashMap) Insert(key uint32, handle uint16) bool{
	return false
}

func (m *HandleHashMap) RemoveByKey(key uint32) bool{
	return false
}

func (m *HandleHashMap) RemoveByHandle(handle uint16) bool{
	return false
}

func (m *HandleHashMap) Find(key uint32) uint16 {
	return 0
}

func (m *HandleHashMap) Reset() {

}

func (m *HandleHashMap) GetNumElements() uint32 {
	return 0
}

func (m *HandleHashMap) GetMaxCapacity() uint32 {
	return 0
}



