package gfx

type ViewState struct {
	viewTmp [2][CONFIG_MAX_VIEWS]Matrix4
	viewPro [2][CONFIG_MAX_VIEWS]Matrix4
	view 	[2]*Matrix4
	rect 	Rect
	invView	Matrix4
	invProj Matrix4
	invViewProj Matrix4
	alphaRef	float32
	invViewCached uint16
	invProjCached uint16
	invViewProjCached uint16
}

// 忽略 HMD
func NewViewState(frame *Frame) *ViewState {
	return nil
}

func (vs *ViewState) reset(frame *Frame) {

}

func (vs *ViewState) setPredefined(rc RendererContext, view uint16, eye uint8, render *Frame, draw *RenderItem) {

}



