package game

/**
	场景
 */
type Scene interface {
	Preload()

	Setup(g *Game)

	Update(dt float32)

	Name() string
}
