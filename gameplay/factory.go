package gameplay

import (
	"github.com/go-gl/mathgl/mgl32"

	"korok/gfx"
	"korok/ecs"
	"korok/assets"
	"korok/space"
)

func Create(image string) *Sprite{
	// Create Id
	id := ecs.Create()

	// Create RenderComponent
	renderComp := gfx.NewRenderComp(id.Index())
	renderComp.Texture2D = *assets.GetTexture(image)
	renderComp.Model = mgl32.Ident4()

	// Create TransformComponent
	transComp := space.NewTransform(id.Index())
	transComp.Position = mgl32.Vec2{100, 100}

	//

	return &Sprite{}
}
