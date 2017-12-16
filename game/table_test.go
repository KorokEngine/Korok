package game

import (
	"testing"
	"korok.io/korok/engi"
)

// Test CRUD operation for MeshTable
func TestScriptTable(t *testing.T) {
	em := &engi.EntityManager{}
	st := NewScriptTable(1024)

	e1 := em.New()
	xf1 := st.NewComp(e1, nil)

	if xf := st.Comp(e1); xf != xf1 {
		t.Error("fail to create Comp")
	}

	st.Delete(e1)
	if xf := st.Comp(e1); xf != nil {
		t.Error("fail to delete Comp")
	}

	if size, _ := st.Size(); size != 0 {
		t.Error("fail to reset Table state")
	}

	// create 10
	eList := make([]engi.Entity, 10)
	for i := 0; i < 10; i++ {
		e := em.New()
		st.NewComp(e, nil)
		eList[i] = e
	}

	if size, _ := st.Size(); size != len(eList) {
		t.Error("fail to create 10 Comps")
	}

	// delete 5
	for i, e := range eList {
		if i % 2 == 0 {
			st.Delete(e)
		}
	}

	if size, _ := st.Size(); size != len(eList)/2 {
		t.Error("fail to delete Comps")
	}

	// test left
	for i, e := range eList {
		if i % 2 == 1 {
			if comp := st.Comp(e); comp == nil || comp.Entity != e {
				t.Error("fail to keep entity:", e)
			}
		} else {
			if st.Comp(e) != nil {
				t.Error("fail to delete Comps:", e)
			}
		}
	}


}

// step=64
func TestScriptTableResize(t *testing.T) {
	em := &engi.EntityManager{}
	st := NewScriptTable(1024)

	list30 := make([]engi.Entity, 30)
	for i := 0; i < 30; i++ {
		e := em.New()
		st.NewComp(e, nil)
		list30[i] = e
	}

	// will cause resize
	list100 := make([]engi.Entity, 100)
	for i := 0; i < 100; i++ {
		e := em.New()
		st.NewComp(e, nil)
		list100[i] = e
	}

	if size, _ := st.Size(); size != (len(list30) + len(list100)) {
		t.Errorf("fail to create Comps: %d/%d", size, len(list30) + len(list100))
	}

	list := append(list100, list30...)
	for _, e := range list {
		if xf := st.Comp(e); xf == nil || xf.Entity != e {
			t.Error("comp is not create correctly")
		}
	}
}
