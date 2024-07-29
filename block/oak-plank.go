package block

import rl "github.com/gen2brain/raylib-go/raylib"

type OakPlank struct {
	model    rl.Model
	texture  rl.Texture2D
	position rl.Vector3
}

func NewOakPlank(x float32, y float32, z float32) Dirt {
	texture := rl.LoadTexture("texture/oak-plank.png")
	mesh := rl.GenMeshCube(2.0, 2.0, 2.0)
	material := rl.LoadMaterialDefault()
	material.GetMap(rl.MapAlbedo).Texture = texture

	model := rl.LoadModelFromMesh(mesh)
	model.GetMaterials()[0] = material

	return Dirt{
		model:    model,
		texture:  texture,
		Position: rl.NewVector3(x, y, z),
	}
}

func (d *OakPlank) Render() {
	rl.DrawModel(d.model, d.position, 1.0, rl.White)
}

// Unload releases the resources used by the Dirt object
func (d *OakPlank) Unload() {
	rl.UnloadTexture(d.texture)
	rl.UnloadModel(d.model)
}
