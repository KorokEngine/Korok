package array

const STEP  = 64

type DenseIntMap struct {
	entries []int
}

func (m *DenseIntMap) Grow(size int) {
	slice := make([]int, size)
	copy(slice, m.entries)
	m.entries = slice
}

func (m *DenseIntMap) Compact() {

}

func (m *DenseIntMap) Put(k, v int) {
	if n := len(m.entries); k >= n {
		m.Grow(k + 64)
	}
	m.entries[k] = v
}

func (m *DenseIntMap) Get(k int) (v int, ok bool) {
	if n := len(m.entries); k < n {
		if vv := m.entries[k]; vv != 0 {
			v, ok = vv, true
		}
	}
	return
}

func (m *DenseIntMap) Del(k int) {
	if n := len(m.entries); k < n {
		m.entries[k] = 0
	}
}

func (m *DenseIntMap) Clear() {

}


type SparseIntMap struct {
	keys []int
	vals []int
	used int
}

func (m *SparseIntMap) Grow(size int) {
	newKeys := make([]int, size)
	newVals := make([]int, size)
	copy(newKeys, m.keys)
	copy(newVals, m.vals)
	m.keys = newKeys
	m.vals = newVals
}

func (m *SparseIntMap) Compact() {

}

func (m *SparseIntMap) Put(k, v int) {
	m.keys[m.used] = k
	m.vals[m.used] = v
	m.used ++
}

func (m *SparseIntMap) Get(k int) (v int) {
	v = 0
	return
}

func (m *SparseIntMap) Del(k int) {

}

