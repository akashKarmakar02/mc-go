package main

import (
	"math"
	"math/rand"
	"runtime"
	"sync"

	"github.com/aquilax/go-perlin"
	rl "github.com/gen2brain/raylib-go/raylib"

	"mcgo/voxel"
)

const (
	gridSize     = 30
	blockSpacing = 2.0

	perlinScale = 0.1
	perlinAlpha = 2.0
	perlinBeta  = 2.0
	perlinN     = 3
)

type RenderData struct {
	Block        voxel.VoxelBlock
	VisibleFaces []int
}

func main() {
	rl.InitWindow(0, 0, "Voxel Terrain Optimized")
	defer rl.CloseWindow()
	rl.ToggleFullscreen()
	rl.DisableCursor()
	rl.SetTargetFPS(0) // uncapped

	camera := rl.Camera{
		Position: rl.NewVector3(30, 40, 30),
		Target:   rl.NewVector3(0, 0, 0),
		Up:       rl.NewVector3(0, 1, 0),
		Fovy:     60,
	}

	world := make(map[[3]int]voxel.VoxelBlock)

	// Perlin noise terrain generation
	seed := rand.Int63()
	noise := perlin.NewPerlin(perlinAlpha, perlinBeta, perlinN, seed)

	for x := -gridSize / 2; x < gridSize/2; x++ {
		for z := -gridSize / 2; z < gridSize/2; z++ {
			raw := noise.Noise2D(float64(x)*perlinScale, float64(z)*perlinScale)
			height := int(math.Floor((raw + 1.0) * 0.5 * 6)) // Clamp to 0–6

			for y := 0; y <= height; y++ {
				blockType := voxel.BlockDirt
				if y == height {
					blockType = voxel.BlockGrass
				}

				gridPos := [3]int{x, y, z}
				world[gridPos] = voxel.VoxelBlock{
					Position: rl.NewVector3(float32(x)*blockSpacing, float32(y)*blockSpacing, float32(z)*blockSpacing),
					Type:     blockType,
				}
			}
		}
	}

	// Pre-allocate slices
	renderQueue := make([]RenderData, 0, len(world))
	// visibleFaces := make([]int, 0, 6)

	for !rl.WindowShouldClose() {
		rl.UpdateCamera(&camera, rl.CameraFree)

		renderQueue = renderQueue[:0] // reuse slice

		// MULTITHREADED face culling
		var wg sync.WaitGroup
		jobs := make(chan [3]int, len(world))
		results := make(chan RenderData, len(world))

		numWorkers := runtime.NumCPU()
		wg.Add(numWorkers)
		for i := 0; i < numWorkers; i++ {
			go func() {
				defer wg.Done()
				for pos := range jobs {
					block := world[pos]
					if block.Type == voxel.BlockAir {
						continue
					}

					var faces []int
					for face := 0; face < 6; face++ {
						dir := voxel.FaceDirs[face]
						nx, ny, nz := pos[0]+int(dir.X), pos[1]+int(dir.Y), pos[2]+int(dir.Z)
						neighbor, exists := world[[3]int{nx, ny, nz}]
						if !exists || neighbor.Type == voxel.BlockAir {
							faces = append(faces, face)
						}
					}

					// Frustum culling: skip far blocks
					blockPos := block.Position
					if rl.Vector3Distance(camera.Position, blockPos) < 80 {
						results <- RenderData{Block: block, VisibleFaces: faces}
					}
				}
			}()
		}

		for pos := range world {
			jobs <- pos
		}
		close(jobs)

		// Wait in background
		go func() {
			wg.Wait()
			close(results)
		}()

		for r := range results {
			renderQueue = append(renderQueue, r)
		}

		// DRAW
		rl.BeginDrawing()
		rl.ClearBackground(rl.SkyBlue)
		rl.BeginMode3D(camera)

		for _, data := range renderQueue {
			data.Block.Render(data.VisibleFaces)
		}

		rl.EndMode3D()
		rl.DrawFPS(10, 10)
		rl.EndDrawing()
	}
}
