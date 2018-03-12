package gfx

import (
	"testing"
	"korok.io/korok/engi"
	"korok.io/korok/math/f32"
)

func TestTransform(t *testing.T) {
	em := engi.NewEntityManager()
	tt := NewTransformTable(1024)

	e := em.New()
	xf := tt.NewComp(e)

	c1, c2 := em.New(), em.New()
	xf1, xf2 := tt.NewComp(c1), tt.NewComp(c2)

	xf.SetPosition(f32.Vec2{100, 100})

	if xy := xf.Local().Position; xy[0] != 100 || xy[1] != 100 {
		t.Error("xf local position is not set correctly")
	}

	if xy := xf.World().Position; xy[0] != 100 || xy[1] != 100 {
		t.Error("xf world is not set correctly")
	}

	// link c1 and c2 to parent e
	xf.LinkChild(xf1)
	xf.LinkChild(xf2)

	xf1.SetPosition(f32.Vec2{50, 50})
	xf2.SetPosition(f32.Vec2{-50, -50})

	if xy := xf1.World().Position; xy[0] != 150 || xy[1] != 150 {
		t.Error("child1 postion:", xy, "expected:", f32.Vec2{150, 150})
	}

	if xy := xf2.World().Position; xy[0] != 50 || xy[1] != 50 {
		t.Error("child2 postion:", xy, "expected:", f32.Vec2{50, 50})
	}

	xf.SetPosition(f32.Vec2{200, 200})

	if xy := xf.World().Position; xy[0] != 200 || xy[1] != 200 {
		t.Error("xf world is not set correctly")
	}

	xf1, xf2 = tt.Comp(c1), tt.Comp(c2)
	if xy := xf1.World().Position; xy[0] != 250 || xy[1] != 250 {
		t.Error("child1 postion err:", xy, "expected:", f32.Vec2{250, 250})
	}

	if xy := xf2.World().Position; xy[0] != 150 || xy[1] != 150 {
		t.Error("child2 postion err:", xy, "expected:", f32.Vec2{150, 150})
	}
}


// Test CRUD operation for TransformTable
func TestTransformTable(t *testing.T) {
	em := &engi.EntityManager{}
	tt := NewTransformTable(1024)

	e1 := em.New()
	xf1 := tt.NewComp(e1)

	if xf := tt.Comp(e1); xf != xf1 {
		t.Error("fail to create Comp")
	}

	tt.Delete(e1)
	if xf := tt.Comp(e1); xf != nil {
		t.Error("fail to delete Comp")
	}

	if size, _ := tt.Size(); size != 0 {
		t.Error("fail to reset Table state")
	}

	// create 10
	eList := make([]engi.Entity, 10)
	for i := 0; i < 10; i++ {
		e := em.New()
		tt.NewComp(e)
		eList[i] = e
	}

	if size, _ := tt.Size(); size != len(eList) {
		t.Error("fail to create 10 Comps")
	}

	// delete 5
	for i, e := range eList {
		if i % 2 == 0 {
			tt.Delete(e)
		}
	}

	if size, _ := tt.Size(); size != len(eList)/2 {
		t.Error("fail to delete Comps")
	}

	// test left
	for i, e := range eList {
		if i % 2 == 1 {
			if comp := tt.Comp(e); comp == nil || comp.Entity != e {
				t.Error("fail to keep entity:", e)
			}
		} else {
			if tt.Comp(e) != nil {
				t.Error("fail to delete Comps:", e)
			}
		}
	}
}

// step=64
func TestTransformTableResize(t *testing.T) {
	em := &engi.EntityManager{}
	tt := NewTransformTable(1024)

	list30 := make([]engi.Entity, 30)
	for i := 0; i < 30; i++ {
		e := em.New()
		tt.NewComp(e)
		list30[i] = e
	}

	// will cause resize
	list100 := make([]engi.Entity, 100)
	for i := 0; i < 100; i++ {
		e := em.New()
		tt.NewComp(e)
		list100[i] = e
	}

	if size, _ := tt.Size(); size != (len(list30) + len(list100)) {
		t.Errorf("fail to create Comps: %d/%d", size, len(list30) + len(list100))
	}

	list := append(list100, list30...)
	for _, e := range list {
		if xf := tt.Comp(e); xf == nil || xf.Entity != e {
			t.Error("comp is not create correctly")
		}
	}
}

// test node-link
func TestParentChildLink(t *testing.T) {
	em := engi.NewEntityManager()
	tt := NewTransformTable(1024)

	car := em.New()
	xf := tt.NewComp(car)

	wheel1, wheel2, wheel3, wheel4 := em.New(), em.New(), em.New(), em.New()
	xf1 := tt.NewComp(wheel1)
	xf2 := tt.NewComp(wheel2)
	xf3 := tt.NewComp(wheel3)
	xf4 := tt.NewComp(wheel4)

	xf.LinkChildren(xf1, xf2, xf3, xf4)

	wheels := []*Transform{xf1, xf2, xf3, xf4}
	for _, xf := range wheels {
		if xf.Parent().Entity != car {
			t.Error("fail to link child to parent")
		}
	}

	if xf.FirstChild().Entity != wheel1 {
		t.Error("fail to link wheel1")
	}

	for i := 0; i < 3; i++ {
		if _, nxt := wheels[i].Sibling(); nxt.Entity != wheels[i+1].Entity {
			t.Error("fail to link next sibling wheels")
		}
	}

	for i := 3; i >= 1; i-- {
		if pre, _ := wheels[i].Sibling(); pre.Entity != wheels[i-1].Entity {
			t.Error("fail to link prev sibling wheels")
		}
	}
	// remove first wheel
	xf.RemoveChild(xf1)
	if xf1.Parent() != nil {
		t.Error("fail to remove wheel1, wheel ref parent")
	}
	if xf.FirstChild().Entity == wheel1 {
		t.Error("fail to remove wheel1, car ref wheel")
	}

	if xf.FirstChild().Entity != wheel2 {
		t.Error("fail to keep wheel2")
	}

	// remove second wheel
	xf.RemoveChild(xf3)
	if xf.FirstChild().Entity != wheel2 {
		t.Error("fail to keep wheel2")
	}
	if pre, nxt := xf2.Sibling(); pre != nil && nxt.Entity != wheel4 {
		t.Error("fail to keep wheel4")
	}

	// re-construct wheels
	xf.LinkChildren(xf1, xf3)
	wheels = []*Transform{xf2, xf4, xf1, xf3}

	for i := 0; i < 3; i++ {
		if _, nxt := wheels[i].Sibling(); nxt.Entity != wheels[i+1].Entity {
			t.Error("fail to link next sibling wheels")
		}
	}

	// try to delete something. danger!
	tt.Delete(wheel1)
	tt.Delete(wheel3)

	// delete will cause ref change
	xf1, xf2, xf3, xf4 = tt.Comp(wheel1), tt.Comp(wheel2), tt.Comp(wheel3), tt.Comp(wheel4)

	if xf1 != nil {
		t.Error("fail to delete wheel1")
	}
	if xf3 != nil {
		t.Error("fail to delete wheel3")
	}

	if xf.FirstChild().Entity != wheel2 {
		t.Error("fail to keep wheel2")
	}
	if pre, nxt := xf.FirstChild().Sibling(); pre != nil || nxt != xf4 {
		t.Error("fail to keep wheel2 and wheel4")
	}
}