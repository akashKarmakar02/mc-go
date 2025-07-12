package voxel

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	FaceTop = iota
	FaceBottom
	FaceFront
	FaceBack
	FaceLeft
	FaceRight
)

const faceSize = 2.0

type BlockType int

const (
	BlockAir BlockType = iota
	BlockGrass
	BlockDirt
)

type VoxelBlock struct {
	Position rl.Vector3
	Type     BlockType
}

// Directional face offsets for neighbor checking
var FaceDirs = map[int]rl.Vector3{
	FaceTop:    {X: 0, Y: 1, Z: 0},
	FaceBottom: {X: 0, Y: -1, Z: 0},
	FaceFront:  {X: 0, Y: 0, Z: 1},
	FaceBack:   {X: 0, Y: 0, Z: -1},
	FaceLeft:   {X: -1, Y: 0, Z: 0},
	FaceRight:  {X: 1, Y: 0, Z: 0},
}

func (b *VoxelBlock) Render(visibleFaces []int) {
	baseColor := rl.Brown
	if b.Type == BlockGrass {
		baseColor = rl.Green
	}

	lightDir := rl.Vector3{X: 1, Y: -1.5, Z: 1}
	lightDir = rl.Vector3Normalize(lightDir)

	for _, face := range visibleFaces {
		normal := getFaceNormal(face)
		intensity := rl.Vector3DotProduct(normal, lightDir)

		// Clamp intensity to [0.2, 1.0] to avoid full darkness
		intensity = float32(math.Max(0.4, float64(intensity)))

		// Apply lighting to base color
		shaded := rl.NewColor(
			uint8(float32(baseColor.R)*intensity),
			uint8(float32(baseColor.G)*intensity),
			uint8(float32(baseColor.B)*intensity),
			255,
		)

		drawFace(b.Position, face, shaded)
	}
}

func getFaceNormal(face int) rl.Vector3 {
	switch face {
	case FaceTop:
		return rl.NewVector3(0, 1, 0)
	case FaceBottom:
		return rl.NewVector3(0, -1, 0)
	case FaceFront:
		return rl.NewVector3(0, 0, 1)
	case FaceBack:
		return rl.NewVector3(0, 0, -1)
	case FaceLeft:
		return rl.NewVector3(-1, 0, 0)
	case FaceRight:
		return rl.NewVector3(1, 0, 0)
	default:
		return rl.NewVector3(0, 1, 0)
	}
}

func drawFace(pos rl.Vector3, face int, color rl.Color) {
	var p1, p2, p3, p4 rl.Vector3

	half := float32(faceSize / 2)

	switch face {
	case FaceBottom:
		p1 = rl.NewVector3(pos.X-half, pos.Y-half, pos.Z+half)
		p2 = rl.NewVector3(pos.X-half, pos.Y-half, pos.Z-half)
		p3 = rl.NewVector3(pos.X+half, pos.Y-half, pos.Z-half)
		p4 = rl.NewVector3(pos.X+half, pos.Y-half, pos.Z+half)
	case FaceTop:
		p1 = rl.NewVector3(pos.X-half, pos.Y+half, pos.Z+half)
		p2 = rl.NewVector3(pos.X+half, pos.Y+half, pos.Z+half)
		p3 = rl.NewVector3(pos.X+half, pos.Y+half, pos.Z-half)
		p4 = rl.NewVector3(pos.X-half, pos.Y+half, pos.Z-half)
	case FaceFront:
		p1 = rl.NewVector3(pos.X-half, pos.Y+half, pos.Z+half)
		p2 = rl.NewVector3(pos.X-half, pos.Y-half, pos.Z+half)
		p3 = rl.NewVector3(pos.X+half, pos.Y-half, pos.Z+half)
		p4 = rl.NewVector3(pos.X+half, pos.Y+half, pos.Z+half)
	case FaceBack:
		p1 = rl.NewVector3(pos.X+half, pos.Y+half, pos.Z-half)
		p2 = rl.NewVector3(pos.X+half, pos.Y-half, pos.Z-half)
		p3 = rl.NewVector3(pos.X-half, pos.Y-half, pos.Z-half)
		p4 = rl.NewVector3(pos.X-half, pos.Y+half, pos.Z-half)
	case FaceLeft:
		p1 = rl.NewVector3(pos.X-half, pos.Y+half, pos.Z-half)
		p2 = rl.NewVector3(pos.X-half, pos.Y-half, pos.Z-half)
		p3 = rl.NewVector3(pos.X-half, pos.Y-half, pos.Z+half)
		p4 = rl.NewVector3(pos.X-half, pos.Y+half, pos.Z+half)
	case FaceRight:
		p1 = rl.NewVector3(pos.X+half, pos.Y+half, pos.Z+half)
		p2 = rl.NewVector3(pos.X+half, pos.Y-half, pos.Z+half)
		p3 = rl.NewVector3(pos.X+half, pos.Y-half, pos.Z-half)
		p4 = rl.NewVector3(pos.X+half, pos.Y+half, pos.Z-half)
	}

	// Ensure winding is CCW
	rl.DrawTriangle3D(p1, p2, p3, color)
	rl.DrawTriangle3D(p1, p3, p4, color)

	// Optional: draw border
	borderColor := rl.Black
	rl.DrawLine3D(p1, p2, borderColor)
	rl.DrawLine3D(p2, p3, borderColor)
	rl.DrawLine3D(p3, p4, borderColor)
	rl.DrawLine3D(p4, p1, borderColor)
}
