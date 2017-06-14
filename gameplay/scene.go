package gameplay

/**
	场景
 */
type Scene interface {
	Preload()

	Setup()

	Update(dt float32)

	Name() string
}
