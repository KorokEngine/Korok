package gameplay

/**
	游戏对象绑定脚本/行为
 */
type BehaviorComp interface {
	Init(id uint32)

	Update(dt float32)

	Destroy()
}

type BehaviorSystem struct {
}

func (*BehaviorSystem) AddBehavior(behavior BehaviorComp) {
}

func (*BehaviorSystem) Delete(id uint32) {
	// todo
}

func (*BehaviorSystem) GetComp(id uint32) BehaviorComp{
	return nil
}

func (*BehaviorSystem) Update(dt float32) {

}



