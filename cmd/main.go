package main

import (
	"math"
	"math/rand"
	"mcgo/block"
	"mcgo/entity"

	"github.com/aquilax/go-perlin"
	rl "github.com/gen2brain/raylib-go/raylib"

	gui "github.com/gen2brain/raylib-go/raygui"
)

const (
	gridSize    = 70
	gridSpacing = 2.0
	perlinScale = 0.1
	perlinAlpha = 3.2
	perlinBeta  = 1.0
	perlinN     = 5
	numWorkers  = 4
	maxDistance = 70.0
)

type Chunk struct {
	Blocks []block.Grass
}

func generateTerrain(grassModel rl.Model, grassTexture rl.Texture2D, dirtModel rl.Model, dirtTexture rl.Texture2D) ([]block.Grass, []block.Dirt) {
	seed := rand.Int63()
	p := perlin.NewPerlin(perlinAlpha, perlinBeta, perlinN, seed)

	var grassBlocks []block.Grass
	var dirtBlocks []block.Dirt

	for x := -gridSize / 2; x < gridSize/2; x++ {
		for z := -gridSize / 2; z < gridSize/2; z++ {
			originalHeight := p.Noise2D(float64(x)*perlinScale, float64(z)*perlinScale) * 10
			height := math.Max(0, math.Floor(originalHeight))

			for y := 0; y <= int(height); y++ {
				if y == int(height) {
					grassBlock := block.NewGrass(float32(x)*gridSpacing, float32(y)*2, float32(z)*gridSpacing, grassModel, grassTexture)
					grassBlocks = append(grassBlocks, grassBlock)
				} else {
					dirtBlock := block.NewDirt(float32(x)*gridSpacing, float32(y)*2+1, float32(z)*gridSpacing, dirtModel, dirtTexture)
					dirtBlocks = append(dirtBlocks, dirtBlock)
				}
			}
		}
	}

	return grassBlocks, dirtBlocks
}

func renderBlocks(grassBlocks []block.Grass, dirtBlocks []block.Dirt, camera rl.Camera, visibleGrassBlocks *[]block.Grass, visibleDirtBlocks *[]block.Dirt) {
	*visibleGrassBlocks = (*visibleGrassBlocks)[:0]
	*visibleDirtBlocks = (*visibleDirtBlocks)[:0]

	for _, grassBlock := range grassBlocks {
		if block.IsWithinRenderDistance(grassBlock.Position, camera.Position, maxDistance) {
			*visibleGrassBlocks = append(*visibleGrassBlocks, grassBlock)
		}
	}

	for _, dirtBlock := range dirtBlocks {
		if block.IsWithinRenderDistance(dirtBlock.Position, camera.Position, maxDistance) {
			*visibleDirtBlocks = append(*visibleDirtBlocks, dirtBlock)
		}
	}

	for _, grassBlock := range *visibleGrassBlocks {
		grassBlock.Render(camera)
	}

	for _, dirtBlock := range *visibleDirtBlocks {
		dirtBlock.Render(camera)
	}
}

func main() {
	rl.InitWindow(0, 0, "Minecraft")
	rl.ToggleFullscreen()
	defer rl.CloseWindow()

	grassModel := rl.LoadModel("model/grass.obj")
	grassTexture := rl.LoadTexture("model/grass.png")

	dirtTexture := rl.LoadTexture("texture/dirt.png")
	mesh := rl.GenMeshCube(2.0, 2.0, 2.0)
	dirtMaterial := rl.LoadMaterialDefault()
	dirtMaterial.GetMap(rl.MapAlbedo).Texture = dirtTexture

	model := rl.LoadModelFromMesh(mesh)
	model.GetMaterials()[0] = dirtMaterial

	isGameStarted := false

	grassBlocks, dirtBlocks := generateTerrain(grassModel, grassTexture, model, dirtTexture)

	camera := rl.Camera{
		Position: rl.NewVector3(10.0, 10.0, 10.0),
		Up:       rl.NewVector3(0.0, 1.0, 0.0),
		Fovy:     80.0,
	}
	player := entity.NewPlayer(10, 10, 10, &camera)

	image := rl.LoadImage("assets/bg.jpg")
	defer rl.UnloadImage(image)

	background := rl.LoadTextureFromImage(image)
	defer rl.UnloadTexture(background)

	shader := rl.LoadShader("", "assets/blur.fs")
	defer rl.UnloadShader(shader)

	var visibleGrassBlocks []block.Grass
	var visibleDirtBlocks []block.Dirt

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.EnableDepthMask()
		rl.EnableBackfaceCulling()

		if !isGameStarted {
			var button bool
			rl.ClearBackground(rl.RayWhite)

			rl.BeginShaderMode(shader)
			rl.DrawTexture(background, 0, 0, rl.White)
			rl.EndShaderMode()

			rl.DrawRectangle(0, 0, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()), rl.Fade(rl.Black, 0.5))

			gui.SetStyle(gui.DEFAULT, gui.TEXT_SIZE, 30)
			button = gui.Button(rl.NewRectangle(float32(rl.GetRenderWidth())/2-200, float32(rl.GetRenderHeight())/2-70, 400, 70), "Start Game")
			if button {
				isGameStarted = true
				rl.DisableCursor()
			}
		} else {
			player.Update(grassBlocks)
			rl.ClearBackground(rl.SkyBlue)

			rl.BeginMode3D(camera)
			player.Render()

			renderBlocks(grassBlocks, dirtBlocks, camera, &visibleGrassBlocks, &visibleDirtBlocks)

			rl.EndMode3D()

			rl.DrawFPS(10, 10)
		}

		rl.EndDrawing()
	}
}
