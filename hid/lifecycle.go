package hid

// 窗口生命周期
type WindowCallback interface {
	// 窗口创建
	OnCreate(w, h float32, pixelRatio float32)

	// Resize...
	OnResize(w, h int32)

	// 窗口循环
	OnLoop()

	// 窗口销毁
	OnDestroy()

	// 窗口切回？
	OnResume()

	// 窗口切入后台
	OnPause()

	// 窗口焦点变化
	OnFocusChanged(focused bool)
}

// 输入系统
type InputCallback interface {
	OnKeyEvent(key int, pressed bool)
	OnPointEvent(key int, pressed bool, x, y float32)
}
