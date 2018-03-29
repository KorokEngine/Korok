package game

type LoaderState struct {
	progress int
	done bool
}

type AsyncHandler chan bool

// Loader
type Loader interface {
	Load()
}

type UnLoader interface {
	Unload()
}

// Scene has lifecycle callback.
// OnEnter is called when the scene is shown.
// Update is called each frame.
// OnExit is called when the scene is hided.
type Scene interface {
	OnEnter(g *Game)
	Update(dt float32)
	OnExit()
}

// SceneManager manages scenes.
type SceneManager struct {
	g *Game

	stack []Scene
	hScene Scene
}

// 设计合理的异步加载方案：TODO
func (*SceneManager) Load(sn Scene) {
	if loader, ok := sn.(Loader); ok {
		loader.Load()
	}
}

func (*SceneManager) UnLoad(sn Scene) {
	if loader, ok := sn.(UnLoader); ok {
		loader.Unload()
	}
}

func (sm *SceneManager) Setup(g *Game) {
	sm.g = g

	// special case for default scene.
	if h := sm.hScene; h != nil {
		if loader, ok := h.(Loader); ok {
			loader.Load()
		}
		h.OnEnter(g)
	}
}

func (sm *SceneManager) Update(dt float32) {
	if h := sm.hScene; h != nil {
		h.Update(dt)
	}
}

// SetDefault sets the default Scene before the Game start.
// It's designed for internal usage. You should not use it.
func (sm *SceneManager) SetDefault(sn Scene) {
	sm.hScene = sn
	sm.stack = append(sm.stack, sn)
}

func (sm *SceneManager) Push(sn Scene) {
	if h := sm.hScene; h != nil {
		h.OnExit()
	}

	sm.hScene = sn
	sm.stack = append(sm.stack, sn)

	// setup
	sn.OnEnter(sm.g)
}

func (sm *SceneManager) Pop() (sn Scene, ok bool) {
	if size := len(sm.stack); size > 0 {
		sn = sm.stack[size-1]
		ok = true
		sm.stack = sm.stack[:size-1]
		sn.OnExit()
		sm.UnLoad(sn)

		if next := size-2; next >= 0 {
			sn := sm.stack[next]
			sm.hScene = sn
			sn.OnEnter(sm.g)
		}
	}
	return
}

func (sm *SceneManager) Peek() (sn Scene, ok bool) {
	if size := len(sm.stack); size > 0 {
		sn = sm.stack[size-1]
		ok = true
	}
	return
}

