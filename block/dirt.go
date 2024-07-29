package block

import rl "github.com/gen2brain/raylib-go/raylib"

type Dirt struct {
	model    rl.Model
	texture  rl.Texture2D
	Position rl.Vector3
}

func NewDirt(x float32, y float32, z float32, model rl.Model, texture rl.Texture2D) Dirt {
	return Dirt{
		model:    model,
		texture:  texture,
		Position: rl.NewVector3(x, y, z),
	}
}

func (d *Dirt) Render(camera rl.Camera) {
	rl.DrawModel(d.model, d.Position, 1.0, rl.White)
}

// Unload releases the resources used by the Dirt object
func (d *Dirt) Unload() {
	rl.UnloadTexture(d.texture)
	rl.UnloadModel(d.model)
}
