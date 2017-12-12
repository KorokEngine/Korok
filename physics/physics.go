package physics

import (
	"korok.io/box2d"
	"korok.io/korok/engi"
	"korok.io/box2d/common/f32"
	"korok.io/box2d/dynamics/rigid"
	"korok.io/box2d/collision/shapes"

	"log"
)

func init() {

}

// 刚体/质点
type RigidBodyComp struct {

}

func (*RigidBodyComp) MoveTo(x, y float32)  {

}

func (*RigidBodyComp) MoveBy(dx, dy float32)  {

}

type RigidBodyTable struct {
	_comps []RigidBodyComp
	_index uint32
	_map   map[int]uint32
}

func (bt *RigidBodyTable) NewComp(entity engi.Entity) (bc *RigidBodyComp) {
	bc = &bt._comps[bt._index]
	bt._map[int(entity)] = bt._index
	bt._index ++
	return
}

func (bt *RigidBodyTable) Comp(entity engi.Entity) (bc *RigidBodyComp) {
	if v, ok := bt._map[int(entity)]; ok{
		bc = &bt._comps[v]
	}
	return
}

// TODO impl
func (bt *RigidBodyTable) Delete(entity engi.Entity) (bc *RigidBodyComp) {
	return nil
}

// 可碰撞组件
type ColliderComp struct {
	*rigid.Body
}

func (th *ColliderComp) SetPosition(x, y float32) {

}

func (th *ColliderComp) SetSize(w, h float32) {

}

type ColliderTable struct {
	_comps []ColliderComp
	_index uint32
	_map   map[int]uint32
}

func (ct *ColliderTable) NewComp(entity engi.Entity) (cc *ColliderComp){
	cc = &ct._comps[ct._index]
	ct._map[int(entity)] = ct._index
	ct._index ++
	return
}

func (ct *ColliderTable) Comp(entity engi.Entity) (cc *ColliderComp) {
	if v, ok := ct._map[int(entity)]; ok {
		cc = &ct._comps[v]
	}
	return
}

// TODO impl
func (ct *ColliderTable) Delete(entity engi.Entity) (cc *ColliderComp) {
	return nil
}

type CollisionSystem struct {
	B2World *box2d.World
}

func NewCollisionSystem() *CollisionSystem {
	th := new(CollisionSystem)
	gravity := f32.Vec2{0, -10}
	th.B2World = box2d.CreateWorld(gravity)

	// Add Body
	bodyDef := rigid.BodyDef{}
	bodyDef.Initialize()
	bodyDef.Position = f32.Vec2{0, -10}

	groundBody := th.B2World.CreateBody(&bodyDef)
	// groundBody

	// Shape
	groundBox := shapes.PolygonShape{}
	groundBox.Initialize()
	groundBox.SetAsBox(50, 10)

	// Fixture
	groundBody.CreateFixtureWithShape(&groundBox, 0)
	return th
}

func (th *CollisionSystem) NewCollisionComp(id uint32) {

}

// a body <position, shape>
func (th *CollisionSystem) NewBody(x, y float32, w, h float32, id uint32) {
	log.Println("create body for id:", id, "w=", w, ", h=", h, "  x=", x, ", y=", y)

	// Add Body
	bodyDef := rigid.BodyDef{}
	bodyDef.Initialize()
	bodyDef.Type = rigid.DynamicBody
	bodyDef.Position = f32.Vec2{x, y}
	bodyDef.UserData = id
	bodyDef.Bullet = true
	body := th.B2World.CreateBody(&bodyDef)

	// Shape
	shape := shapes.PolygonShape{}
	shape.Initialize()
	shape.SetAsBox(w/2, h/2)

	// Fixture
	fixtureDef := rigid.FixtureDef{}
	fixtureDef.Initialize()
	fixtureDef.Shape = &shape
	fixtureDef.Density = 1
	fixtureDef.Friction = 0.3

	body.CreateFixture(&fixtureDef)
}

func (th *CollisionSystem) Update(dt float32) {
	th.B2World.Step(1.0/30, 6, 2)
}

func (th *CollisionSystem) Destroy() {
	th.B2World.Destroy()
}








