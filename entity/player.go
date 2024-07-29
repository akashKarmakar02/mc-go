package entity

import (
	"fmt"
	"math"

	"mcgo/block"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const rotationSpeed = 20.0
const numWorkers = 4

type Player struct {
	model        rl.Model
	camera       *rl.Camera
	position     rl.Vector3
	velocity     rl.Vector3
	onGround     bool
	rotation     float32
	cameraAngleY float32
	cameraAngleX float32
	distance     float32
}

func updateCamera(camera *rl.Camera, player *Player) {
	camera.Target = player.position

	// Calculate the camera position based on angles and distance
	camera.Position.X = player.position.X + player.distance*float32(math.Cos(float64(player.cameraAngleY)))*float32(math.Cos(float64(player.cameraAngleX)))
	camera.Position.Y = player.position.Y + player.distance*float32(math.Sin(float64(player.cameraAngleX)))
	camera.Position.Z = player.position.Z + player.distance*float32(math.Sin(float64(player.cameraAngleY)))*float32(math.Cos(float64(player.cameraAngleX)))

	camera.Up = rl.Vector3{X: 0, Y: 1, Z: 0}
}

func NewPlayer(x, y, z float32, camera *rl.Camera) Player {
	model := rl.LoadModel("model/mc_player.glb")

	return Player{
		model:        model,
		position:     rl.Vector3{X: x, Y: y, Z: z},
		camera:       camera,
		velocity:     rl.Vector3{X: 0, Y: 0, Z: 0},
		onGround:     false,
		rotation:     0,
		cameraAngleY: 0,
		cameraAngleX: 0,
		distance:     10, // Initial camera distance from player
	}
}

func (p *Player) Update(dirtBlocks []block.Grass) {
	handleFall(p, dirtBlocks)

	var targetRotation float32 = p.rotation

	handleInput(p, dirtBlocks, &targetRotation)

	rotationDiff := targetRotation - p.rotation
	if rotationDiff > 180 {
		rotationDiff -= 360
	} else if rotationDiff < -180 {
		rotationDiff += 360
	}

	p.rotation += rotationDiff * rotationSpeed * rl.GetFrameTime()
	p.rotation = float32(math.Mod(float64(p.rotation)+360, 360))

	p.cameraAngleY += rl.GetMouseDelta().X * rl.GetFrameTime() / 40
	p.cameraAngleX += rl.GetMouseDelta().Y * rl.GetFrameTime() / 40

	if p.cameraAngleX > 1.0 {
		p.cameraAngleX = 1.0
	} else if p.cameraAngleX < -1.0 {
		p.cameraAngleX = -1.0
	}

	updateCamera(p.camera, p)
}

func (p *Player) Render() {
	rl.DrawModelEx(p.model, p.position, rl.Vector3{X: 0, Y: 1, Z: 0}, p.rotation, rl.Vector3{X: 2.0, Y: 2.0, Z: 2.0}, rl.RayWhite)
	rl.DrawText(fmt.Sprintf("X: %0.2f, Y: %0.2f, Z: %0.2f", p.position.X, p.position.Y, p.position.Z), 50, 50, 20, rl.RayWhite)
}

func (p *Player) Unload() {
	rl.UnloadModel(p.model)
}

func getSpeed() float32 {
	if rl.IsKeyDown(rl.KeyLeftControl) {
		return 5 * rl.GetFrameTime()
	}
	return 2 * rl.GetFrameTime()
}

func handleInput(p *Player, dirtBlocks []block.Grass, targetRotation *float32) {
	moveSpeed := getSpeed()
	moveVector := rl.Vector3{}

	if rl.IsKeyDown(rl.KeyW) {
		moveVector.X -= moveSpeed
		fmt.Println(p.cameraAngleY)
		*targetRotation = p.cameraAngleY*180/float32(math.Pi) - 90 // Facing left
	}
	if rl.IsKeyDown(rl.KeyS) {
		moveVector.X += moveSpeed
		*targetRotation = p.cameraAngleY*180/float32(math.Pi) + 180
	}
	if rl.IsKeyDown(rl.KeyA) {
		moveVector.Z += moveSpeed
	}
	if rl.IsKeyDown(rl.KeyD) {
		moveVector.Z -= moveSpeed
		*targetRotation = p.cameraAngleY*180/float32(math.Pi) + 90
	}
	if rl.IsKeyDown(rl.KeySpace) && p.onGround {
		p.velocity.Y = 13
	}

	// Rotate the movement vector by the camera's Y angle
	sinAngle := float32(math.Sin(float64(p.cameraAngleY)))
	cosAngle := float32(math.Cos(float64(p.cameraAngleY)))

	transformedMoveVector := rl.Vector3{
		X: moveVector.X*cosAngle - moveVector.Z*sinAngle,
		Y: moveVector.Y,
		Z: moveVector.X*sinAngle + moveVector.Z*cosAngle,
	}

	// Apply movement with collision detection
	if !detectCollision(p.position, transformedMoveVector, dirtBlocks) {
		p.position.X += transformedMoveVector.X
		p.position.Z += transformedMoveVector.Z
	}
}

func detectCollision(position rl.Vector3, moveVector rl.Vector3, dirtBlocks []block.Grass) bool {
	nextPosition := rl.Vector3{
		X: position.X + moveVector.X,
		Y: position.Y,
		Z: position.Z + moveVector.Z,
	}

	collisionChan := make(chan bool, numWorkers)
	blockChunks := len(dirtBlocks) / numWorkers

	for i := 0; i < numWorkers; i++ {
		start := i * blockChunks
		end := start + blockChunks

		if i == numWorkers-1 {
			end = len(dirtBlocks)
		}

		go func(blocks []block.Grass) {
			for _, block := range blocks {
				if nextPosition.X+1 > block.Position.X && block.Position.X+1 > nextPosition.X {
					if nextPosition.Z+1 > block.Position.Z && block.Position.Z+1 > nextPosition.Z {
						if nextPosition.Y > block.Position.Y-1 && nextPosition.Y < block.Position.Y+1 {
							rl.DrawText(fmt.Sprintf("X: %0.2f, Y: %0.2f, Z: %0.2f", block.Position.X, block.Position.Y, block.Position.Z), 50, 50, 20, rl.RayWhite)
							rl.DrawText(fmt.Sprintf("X: %0.2f, Y: %0.2f, Z: %0.2f", nextPosition.X, nextPosition.Y, nextPosition.Z), 50, 80, 20, rl.Brown)
							collisionChan <- true
							return
						}
					}
				}
			}
			collisionChan <- false
		}(dirtBlocks[start:end])
	}

	for i := 0; i < numWorkers; i++ {
		if <-collisionChan {
			return true
		}
	}

	return false
}

func handleFall(p *Player, dirtBlocks []block.Grass) {
	if !p.onGround {
		p.velocity.Y -= 35 * rl.GetFrameTime()
	}

	p.position.X += p.velocity.X * rl.GetFrameTime()
	p.position.Y += p.velocity.Y * rl.GetFrameTime()
	p.position.Z += p.velocity.Z * rl.GetFrameTime()

	p.onGround = false
	for _, block := range dirtBlocks {
		if p.position.X > block.Position.X-2 && p.position.X < block.Position.X+2 &&
			p.position.Z > block.Position.Z-2 && p.position.Z < block.Position.Z+2 &&
			p.position.Y > block.Position.Y-2 && p.position.Y < block.Position.Y+2 {
			p.velocity.Y = 0
			p.onGround = true
		}
	}
}
