package effect

import (
	"testing"
	"korok.io/korok/engi"
)

// Test CRUD operation for MeshTable
func TestMeshTable(t *testing.T) {
	em := &engi.EntityManager{}
	et := NewEffectTable(1024)

	e1 := em.New()
	xf1 := et.NewComp(e1)

	if xf := et.Comp(e1); xf != xf1 {
		t.Error("fail to create Comp")
	}

	et.Delete(e1)
	if xf := et.Comp(e1); xf != nil {
		t.Error("fail to delete Comp")
	}

	if size, _ := et.Size(); size != 0 {
		t.Error("fail to reset Table state")
	}

	// create 10
	eList := make([]engi.Entity, 10)
	for i := 0; i < 10; i++ {
		e := em.New()
		et.NewComp(e)
		eList[i] = e
	}

	if size, _ := et.Size(); size != len(eList) {
		t.Error("fail to create 10 Comps")
	}

	// delete 5
	for i, e := range eList {
		if i % 2 == 0 {
			et.Delete(e)
		}
	}

	if size, _ := et.Size(); size != len(eList)/2 {
		t.Error("fail to delete Comps")
	}

	// test left
	for i, e := range eList {
		if i % 2 == 1 {
			if comp := et.Comp(e); comp == nil || comp.Entity != e {
				t.Error("fail to keep entity:", e)
			}
		} else {
			if et.Comp(e) != nil {
				t.Error("fail to delete Comps:", e)
			}
		}
	}


}

// step=64
func TestMeshTableResize(t *testing.T) {
	em := &engi.EntityManager{}
	et := NewEffectTable(1024)

	list30 := make([]engi.Entity, 30)
	for i := 0; i < 30; i++ {
		e := em.New()
		et.NewComp(e)
		list30[i] = e
	}

	// will cause resize
	list100 := make([]engi.Entity, 100)
	for i := 0; i < 100; i++ {
		e := em.New()
		et.NewComp(e)
		list100[i] = e
	}

	if size, _ := et.Size(); size != (len(list30) + len(list100)) {
		t.Errorf("fail to create Comps: %d/%d", size, len(list30) + len(list100))
	}

	list := append(list100, list30...)
	for _, e := range list {
		if xf := et.Comp(e); xf == nil || xf.Entity != e {
			t.Error("comp is not create correctly")
		}
	}
}