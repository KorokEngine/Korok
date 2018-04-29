package hid

// 窗口生命周期
type WindowCallback interface {
	// 窗口创建
	OnCreate(pixelRatio float32)

	// Resize...
	OnResize(w, h int32)

	// 窗口循环
	OnLoop()

	// 窗口销毁
	OnDestroy()
}

// 输入系统
type InputCallback interface {
	OnKeyEvent(key int, pressed bool)
	OnPointEvent(key int, pressed bool, x, y float32)
}
