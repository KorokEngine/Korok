package gameplay

import (
	"github.com/go-gl/mathgl/mgl32"

	"korok/ecs"
	"korok/space"

	"fmt"
)

func Create(image string) Sprite {
	fmt.Println("Create srpite:" + image)

	// Create Id
	id := ecs.Create()

	// Create RenderComponent
//	renderComp := G.NewRenderComp(id.Index())
//	renderComp.Texture2D = assets.GetTexture(image)
//	renderComp.Model = mgl32.Ident4()

	// Create TransformComponent
	transComp := space.NewTransform(id.Index())
	transComp.Position = mgl32.Vec2{100, 100}

	//
	return Sprite(id)
}
