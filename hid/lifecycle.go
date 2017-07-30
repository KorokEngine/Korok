package hid

type WindowCallback interface {
	// 窗口创建
	OnCreate()

	// 窗口循环
	OnLoop()

	// 窗口销毁
	OnDestroy()
}
