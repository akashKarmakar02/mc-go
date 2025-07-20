package main

import (
	"fmt"
	"math"
	"mcgo/voxel"

	"github.com/aquilax/go-perlin"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(800, 600, "Voxel Engine")
	defer rl.CloseWindow()

	// Do NOT call SetConfigFlags with VSyncHint or remove it if you have
	rl.SetTargetFPS(0)
	rl.DisableCursor()
	voxel.InitGlobalVoxelMesh(1.0)

	// Load some textures
	dirtTexture := rl.LoadTexture("texture/dirt.png")
	grassSideTexture := rl.LoadTexture("texture/grass-side.png")
	grassTopTexture := rl.LoadTexture("texture/grass-top.png")
	defer rl.UnloadTexture(dirtTexture)
	defer rl.UnloadTexture(grassTopTexture)
	defer rl.UnloadTexture(grassSideTexture)

	rl.SetTextureFilter(dirtTexture, rl.FilterPoint)

	// Texture atlas
	atlas := voxel.NewTextureAtlas()
	atlas.SetBlockTexture(voxel.BlockGrass, voxel.FaceTop, grassTopTexture)
	atlas.SetBlockTexture(voxel.BlockGrass, voxel.FaceBottom, dirtTexture)
	atlas.SetBlockTexture(voxel.BlockGrass, voxel.FaceFront, grassSideTexture)
	atlas.SetBlockTexture(voxel.BlockGrass, voxel.FaceBack, grassSideTexture)
	atlas.SetBlockTexture(voxel.BlockGrass, voxel.FaceLeft, grassSideTexture)
	atlas.SetBlockTexture(voxel.BlockGrass, voxel.FaceRight, grassSideTexture)

	// Camera
	camera := rl.NewCamera3D(
		rl.NewVector3(10, 20, 10),
		rl.NewVector3(0, 0, 0),
		rl.NewVector3(0, 1, 0),
		45.0, rl.CameraPerspective,
	)

	// Terrain generation using Perlin
	chunk := voxel.NewChunk()
	perlinGen := perlin.NewPerlin(2, 2, 3, 42) // octaves, persistence, lacunarity, seed

	for x := -32; x <= 32; x++ {
		for z := -32; z <= 32; z++ {
			// Get height from Perlin noise
			noise := perlinGen.Noise2D(float64(x)*0.1, float64(z)*0.1)
			height := int(math.Round((noise + 1.0) * 4.0)) // convert [-1,1] to [0,8]
			for y := 0; y <= height; y++ {
				block := voxel.Block{
					Position: rl.NewVector3(float32(x), float32(y), float32(z)),
					Type:     voxel.BlockGrass,
					Size:     1.0,
				}
				chunk.AddBlock(block)
			}
		}
	}

	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()
		updateFreeCamera(&camera, dt)

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		rl.BeginMode3D(camera)

		chunk.Render(atlas, voxel.GlobalVoxelMesh)

		rl.EndMode3D()
		rl.DrawText(fmt.Sprintf("FPS: %d", rl.GetFPS()), 10, 10, 20, rl.Black)
		rl.EndDrawing()
	}
}

var yaw float32 = -90.0 // Start facing forward
var pitch float32 = 0.0 // Looking level

func updateFreeCamera(cam *rl.Camera3D, dt float32) {
	moveSpeed := float32(10.0) * dt
	sensitivity := float32(100.0) * dt

	mouseDelta := rl.GetMouseDelta()
	yaw += mouseDelta.X * sensitivity
	pitch -= mouseDelta.Y * sensitivity

	// Clamp pitch
	if pitch > 89.0 {
		pitch = 89.0
	}
	if pitch < -89.0 {
		pitch = -89.0
	}

	// Calculate new direction
	dirX := float32(math.Cos(float64(rl.Deg2rad*yaw)) * math.Cos(float64(rl.Deg2rad*pitch)))
	dirY := float32(math.Sin(float64(rl.Deg2rad * pitch)))
	dirZ := float32(math.Sin(float64(rl.Deg2rad*yaw)) * math.Cos(float64(rl.Deg2rad*pitch)))
	direction := rl.NewVector3(dirX, dirY, dirZ)
	cam.Target = rl.Vector3Add(cam.Position, rl.Vector3Normalize(direction))

	// WASD movement
	forward := rl.Vector3Normalize(rl.Vector3Subtract(cam.Target, cam.Position))
	right := rl.Vector3Normalize(rl.Vector3CrossProduct(forward, cam.Up))

	if rl.IsKeyDown(rl.KeyW) {
		cam.Position = rl.Vector3Add(cam.Position, rl.Vector3Scale(forward, moveSpeed))
	}
	if rl.IsKeyDown(rl.KeyS) {
		cam.Position = rl.Vector3Subtract(cam.Position, rl.Vector3Scale(forward, moveSpeed))
	}
	if rl.IsKeyDown(rl.KeyA) {
		cam.Position = rl.Vector3Subtract(cam.Position, rl.Vector3Scale(right, moveSpeed))
	}
	if rl.IsKeyDown(rl.KeyD) {
		cam.Position = rl.Vector3Add(cam.Position, rl.Vector3Scale(right, moveSpeed))
	}

	// Recalculate target each frame
	cam.Target = rl.Vector3Add(cam.Position, rl.Vector3Normalize(direction))
}
