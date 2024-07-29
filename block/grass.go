package block

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Grass struct {
	Model    rl.Model
	Texture  rl.Texture2D
	Position rl.Vector3
}

func NewGrass(x, y, z float32, model rl.Model, texture rl.Texture2D) Grass {
	material := model.GetMaterials()[0]
	material.GetMap(rl.MapDiffuse).Texture = texture
	model.GetMaterials()[0] = material

	return Grass{
		Model:    model,
		Texture:  texture,
		Position: rl.NewVector3(x, y, z),
	}
}

func (c *Grass) Render(camera rl.Camera) {
	maxDistance := 70.0

	if IsWithinRenderDistance(c.Position, camera.Position, maxDistance) {
		rl.DrawModel(c.Model, c.Position, 1.0, rl.RayWhite)
	}
}

func IsWithinRenderDistance(blockPos, cameraPos rl.Vector3, maxDistance float64) bool {
	distance := rl.Vector3Distance(blockPos, cameraPos)
	return distance <= float32(maxDistance)
}

func (c *Grass) Unload() {
	rl.UnloadTexture(c.Texture)
	rl.UnloadModel(c.Model)
}
