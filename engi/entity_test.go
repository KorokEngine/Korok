package engi

import "testing"

func TestEntityCreate(t *testing.T) {
	em := NewEntityManager()
	e := em.New()

	if !em.Alive(e) {
		t.Error("fail to create entity")
	}

	em.Destroy(e)

	if em.Alive(e) {
		t.Error("fail to destroy entity")
	}

	e1 := em.New()
	if e.Index() != e1.Index() {
		t.Error("fail to reuse index")
	}

	if e1.Gene() != (e.Gene() + 1) {
		t.Error("fail to compute generation")
	}
}