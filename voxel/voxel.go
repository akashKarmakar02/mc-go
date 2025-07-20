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

type BlockType int

const (
	BlockAir BlockType = iota
	BlockGrass
	BlockDirt
	BlockStone
)

// Face represents a single face of a voxel block
type Face struct {
	ID       int
	Normal   rl.Vector3
	Vertices [4]rl.Vector3 // 4 corners of the face
	UVs      [4]rl.Vector2 // UV coordinates for each vertex
}

// VoxelMesh is a reusable mesh for rendering voxel blocks
type VoxelMesh struct {
	faces       map[int]Face
	size        float32
	VAO         uint32
	VBO         uint32
	initialized bool
}

// RenderOptions defines how to render a voxel block
type RenderOptions struct {
	Position     rl.Vector3
	Size         float32
	Color        rl.Color
	Texture      *rl.Texture2D // Optional texture
	Textures     []*rl.Texture2D
	VisibleFaces []int      // Which faces to render
	Lighting     bool       // Enable/disable lighting
	LightDir     rl.Vector3 // Light direction
	Wireframe    bool       // Render wireframe
}

// Global reusable voxel mesh
var GlobalVoxelMesh *VoxelMesh

// Initialize creates a reusable voxel mesh
func InitVoxelMesh(size float32) *VoxelMesh {
	mesh := &VoxelMesh{
		faces: make(map[int]Face),
		size:  size,
	}

	half := size / 2.0

	// Define all 6 faces with their vertices and UVs
	mesh.faces[FaceTop] = Face{
		ID:     FaceTop,
		Normal: rl.NewVector3(0, 1, 0),
		Vertices: [4]rl.Vector3{
			{X: -half, Y: half, Z: half},  // Front-left
			{X: half, Y: half, Z: half},   // Front-right
			{X: half, Y: half, Z: -half},  // Back-right
			{X: -half, Y: half, Z: -half}, // Back-left
		},
		UVs: [4]rl.Vector2{
			{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}, {X: 0, Y: 1},
		},
	}

	mesh.faces[FaceBottom] = Face{
		ID:     FaceBottom,
		Normal: rl.NewVector3(0, -1, 0),
		Vertices: [4]rl.Vector3{
			{X: -half, Y: -half, Z: -half}, // Back-left
			{X: half, Y: -half, Z: -half},  // Back-right
			{X: half, Y: -half, Z: half},   // Front-right
			{X: -half, Y: -half, Z: half},  // Front-left
		},
		UVs: [4]rl.Vector2{
			{X: 0, Y: 1}, {X: 1, Y: 1}, {X: 1, Y: 0}, {X: 0, Y: 0},
		},
	}

	mesh.faces[FaceFront] = Face{
		ID:     FaceFront,
		Normal: rl.NewVector3(0, 0, 1),
		Vertices: [4]rl.Vector3{
			{X: -half, Y: -half, Z: half}, // Bottom-left
			{X: half, Y: -half, Z: half},  // Bottom-right
			{X: half, Y: half, Z: half},   // Top-right
			{X: -half, Y: half, Z: half},  // Top-left
		},
		UVs: [4]rl.Vector2{
			{X: 0, Y: 1}, {X: 1, Y: 1}, {X: 1, Y: 0}, {X: 0, Y: 0},
		},
	}

	mesh.faces[FaceBack] = Face{
		ID:     FaceBack,
		Normal: rl.NewVector3(0, 0, -1),
		Vertices: [4]rl.Vector3{
			{X: half, Y: -half, Z: -half},  // Bottom-left
			{X: -half, Y: -half, Z: -half}, // Bottom-right
			{X: -half, Y: half, Z: -half},  // Top-right
			{X: half, Y: half, Z: -half},   // Top-left
		},
		UVs: [4]rl.Vector2{
			{X: 0, Y: 1}, {X: 1, Y: 1}, {X: 1, Y: 0}, {X: 0, Y: 0},
		},
	}

	mesh.faces[FaceLeft] = Face{
		ID:     FaceLeft,
		Normal: rl.NewVector3(-1, 0, 0),
		Vertices: [4]rl.Vector3{
			{X: -half, Y: -half, Z: -half}, // Bottom-left
			{X: -half, Y: -half, Z: half},  // Bottom-right
			{X: -half, Y: half, Z: half},   // Top-right
			{X: -half, Y: half, Z: -half},  // Top-left
		},
		UVs: [4]rl.Vector2{
			{X: 0, Y: 1}, {X: 1, Y: 1}, {X: 1, Y: 0}, {X: 0, Y: 0},
		},
	}

	mesh.faces[FaceRight] = Face{
		ID:     FaceRight,
		Normal: rl.NewVector3(1, 0, 0),
		Vertices: [4]rl.Vector3{
			{X: half, Y: -half, Z: half},  // Bottom-left
			{X: half, Y: -half, Z: -half}, // Bottom-right
			{X: half, Y: half, Z: -half},  // Top-right
			{X: half, Y: half, Z: half},   // Top-left
		},
		UVs: [4]rl.Vector2{
			{X: 0, Y: 1}, {X: 1, Y: 1}, {X: 1, Y: 0}, {X: 0, Y: 0},
		},
	}

	return mesh
}

// RenderVoxel renders a voxel block with the specified options
func (mesh *VoxelMesh) RenderVoxel(options RenderOptions) {
	if len(options.VisibleFaces) == 0 {
		return
	}

	// Set up rendering state
	rl.PushMatrix()
	rl.Translatef(options.Position.X, options.Position.Y, options.Position.Z)

	// Scale if size is different from mesh size
	if options.Size != mesh.size && options.Size > 0 {
		scale := options.Size / mesh.size
		rl.Scalef(scale, scale, scale)
	}

	// Render each visible face
	for _, faceID := range options.VisibleFaces {
		face, exists := mesh.faces[faceID]
		if !exists {
			continue
		}

		// Calculate lighting if enabled
		color := options.Color
		if options.Lighting {
			intensity := rl.Vector3DotProduct(face.Normal, rl.Vector3Normalize(options.LightDir))
			intensity = float32(math.Max(0.3, math.Min(1.0, float64(-intensity))))

			color = rl.NewColor(
				uint8(float32(options.Color.R)*intensity),
				uint8(float32(options.Color.G)*intensity),
				uint8(float32(options.Color.B)*intensity),
				options.Color.A,
			)
		}

		if options.Wireframe {
			mesh.renderFaceWireframe(face, color)
		} else {
			tex := options.Textures[faceID]
			mesh.renderFaceSolid(face, color, tex)
		}
	}

	rl.PopMatrix()
}

// renderFaceSolid renders a face as solid triangles
func (mesh *VoxelMesh) renderFaceSolid(face Face, color rl.Color, texture *rl.Texture2D) {
	if texture != nil {
		// Enable texture and set it
		rl.SetTexture(texture.ID)

		rl.Begin(rl.Quads)
		rl.Color4ub(color.R, color.G, color.B, color.A)

		// Render as a quad for better texture mapping
		rl.TexCoord2f(face.UVs[0].X, face.UVs[0].Y)
		rl.Vertex3f(face.Vertices[0].X, face.Vertices[0].Y, face.Vertices[0].Z)

		rl.TexCoord2f(face.UVs[1].X, face.UVs[1].Y)
		rl.Vertex3f(face.Vertices[1].X, face.Vertices[1].Y, face.Vertices[1].Z)

		rl.TexCoord2f(face.UVs[2].X, face.UVs[2].Y)
		rl.Vertex3f(face.Vertices[2].X, face.Vertices[2].Y, face.Vertices[2].Z)

		rl.TexCoord2f(face.UVs[3].X, face.UVs[3].Y)
		rl.Vertex3f(face.Vertices[3].X, face.Vertices[3].Y, face.Vertices[3].Z)

		rl.End()
		rl.SetTexture(0) // Disable texture
	} else {
		// Render without texture
		rl.Begin(rl.Triangles)
		rl.Color4ub(color.R, color.G, color.B, color.A)

		// First triangle (0, 1, 2)
		rl.Vertex3f(face.Vertices[0].X, face.Vertices[0].Y, face.Vertices[0].Z)
		rl.Vertex3f(face.Vertices[1].X, face.Vertices[1].Y, face.Vertices[1].Z)
		rl.Vertex3f(face.Vertices[2].X, face.Vertices[2].Y, face.Vertices[2].Z)

		// Second triangle (0, 2, 3)
		rl.Vertex3f(face.Vertices[0].X, face.Vertices[0].Y, face.Vertices[0].Z)
		rl.Vertex3f(face.Vertices[2].X, face.Vertices[2].Y, face.Vertices[2].Z)
		rl.Vertex3f(face.Vertices[3].X, face.Vertices[3].Y, face.Vertices[3].Z)

		rl.End()
	}
}

// renderFaceWireframe renders a face as wireframe
func (mesh *VoxelMesh) renderFaceWireframe(face Face, color rl.Color) {
	rl.Begin(rl.Lines)
	rl.Color4ub(color.R, color.G, color.B, color.A)

	// Draw the 4 edges of the face
	for i := 0; i < 4; i++ {
		next := (i + 1) % 4

		rl.Vertex3f(face.Vertices[i].X, face.Vertices[i].Y, face.Vertices[i].Z)
		rl.Vertex3f(face.Vertices[next].X, face.Vertices[next].Y, face.Vertices[next].Z)
	}

	rl.End()
}

// Convenience functions for common use cases

// RenderColoredVoxel renders a voxel with a solid color
func (mesh *VoxelMesh) RenderColoredVoxel(position rl.Vector3, size float32, color rl.Color, visibleFaces []int) {
	options := RenderOptions{
		Position:     position,
		Size:         size,
		Color:        color,
		VisibleFaces: visibleFaces,
		Lighting:     true,
		LightDir:     rl.NewVector3(1, -1, 1),
	}
	mesh.RenderVoxel(options)
}

// RenderTexturedVoxel renders a voxel with a texture
func (mesh *VoxelMesh) RenderTexturedVoxel(position rl.Vector3, size float32, textures []*rl.Texture2D, tint rl.Color, visibleFaces []int) {
	options := RenderOptions{
		Position:     position,
		Size:         size,
		Color:        tint,
		Textures:     textures,
		VisibleFaces: visibleFaces,
		Lighting:     true,
		LightDir:     rl.NewVector3(1, -1, 1),
	}
	mesh.RenderVoxel(options)
}

// RenderWireframeVoxel renders a voxel as wireframe
func (mesh *VoxelMesh) RenderWireframeVoxel(position rl.Vector3, size float32, color rl.Color, visibleFaces []int) {
	options := RenderOptions{
		Position:     position,
		Size:         size,
		Color:        color,
		VisibleFaces: visibleFaces,
		Wireframe:    true,
	}
	mesh.RenderVoxel(options)
}

// Helper functions for face selection

// AllFaces returns all 6 faces
func AllFaces() []int {
	return []int{FaceTop, FaceBottom, FaceFront, FaceBack, FaceLeft, FaceRight}
}

// VisibleFaces returns faces that should be visible based on neighbor presence
func VisibleFaces(neighbors map[int]bool) []int {
	var visible []int
	faces := []int{FaceTop, FaceBottom, FaceFront, FaceBack, FaceLeft, FaceRight}

	for _, face := range faces {
		if !neighbors[face] { // If no neighbor on this side, face is visible
			visible = append(visible, face)
		}
	}

	return visible
}

// ExcludeFaces returns all faces except the specified ones
func ExcludeFaces(excludedFaces []int) []int {
	excluded := make(map[int]bool)
	for _, face := range excludedFaces {
		excluded[face] = true
	}

	var result []int
	allFaces := AllFaces()
	for _, face := range allFaces {
		if !excluded[face] {
			result = append(result, face)
		}
	}

	return result
}

// OnlyFaces returns only the specified faces
func OnlyFaces(faces []int) []int {
	return faces
}

// Block represents a voxel block with position and type
type Block struct {
	Position rl.Vector3
	Type     BlockType
	Size     float32
}

// TextureAtlas manages textures for different block types
type TextureAtlas struct {
	textures map[BlockType]map[int]rl.Texture2D // BlockType -> Face -> Texture
	colors   map[BlockType]rl.Color             // Default colors for block types
}

// NewTextureAtlas creates a new texture atlas
func NewTextureAtlas() *TextureAtlas {
	atlas := &TextureAtlas{
		textures: make(map[BlockType]map[int]rl.Texture2D),
		colors:   make(map[BlockType]rl.Color),
	}

	// Set default colors
	atlas.colors[BlockGrass] = rl.Green
	atlas.colors[BlockDirt] = rl.Brown
	atlas.colors[BlockStone] = rl.Gray

	return atlas
}

// SetBlockTexture sets a texture for a specific block type and face
func (atlas *TextureAtlas) SetBlockTexture(blockType BlockType, face int, texture rl.Texture2D) {
	if atlas.textures[blockType] == nil {
		atlas.textures[blockType] = make(map[int]rl.Texture2D)
	}
	atlas.textures[blockType][face] = texture
}

// SetBlockColor sets the default color for a block type
func (atlas *TextureAtlas) SetBlockColor(blockType BlockType, color rl.Color) {
	atlas.colors[blockType] = color
}

// GetBlockTexture gets the texture for a block type and face
func (atlas *TextureAtlas) GetBlockTexture(blockType BlockType, face int) (rl.Texture2D, bool) {
	if faceTextures, exists := atlas.textures[blockType]; exists {
		if texture, faceExists := faceTextures[face]; faceExists {
			return texture, true
		}
	}
	return rl.Texture2D{}, false
}

// GetBlockColor gets the color for a block type
func (atlas *TextureAtlas) GetBlockColor(blockType BlockType) rl.Color {
	if color, exists := atlas.colors[blockType]; exists {
		return color
	}
	return rl.White
}

// RenderBlock renders a block using the texture atlas
func (atlas *TextureAtlas) RenderBlock(mesh *VoxelMesh, block Block, visibleFaces []int) {
	// Check if we have textures for any faces
	hasTexture := false
	var textures []*rl.Texture2D

	// Try to get texture for the first visible face
	for _, face := range visibleFaces {
		if _, exists := atlas.GetBlockTexture(block.Type, face); exists {
			hasTexture = true
		}
	}

	for i := 0; i < 6; i++ {
		tex, _ := atlas.GetBlockTexture(block.Type, i)
		textures = append(textures, &tex)
	}

	if hasTexture {
		mesh.RenderTexturedVoxel(block.Position, block.Size, textures, rl.White, visibleFaces)
	} else {
		color := atlas.GetBlockColor(block.Type)
		mesh.RenderColoredVoxel(block.Position, block.Size, color, visibleFaces)
	}
}

// Initialize the global voxel mesh (call this once in your main function)
func InitGlobalVoxelMesh(size float32) {
	GlobalVoxelMesh = InitVoxelMesh(size)
}

// Quick render functions using the global mesh
func RenderVoxel(position rl.Vector3, color rl.Color, visibleFaces []int) {
	if GlobalVoxelMesh != nil {
		GlobalVoxelMesh.RenderColoredVoxel(position, GlobalVoxelMesh.size, color, visibleFaces)
	}
}

func RenderTexturedVoxel(position rl.Vector3, texture rl.Texture2D, visibleFaces []int) {
	textures := make([]*rl.Texture2D, 0)
	for i := 0; i < 6; i++ {
		textures = append(textures, &texture)
	}
	if GlobalVoxelMesh != nil {
		GlobalVoxelMesh.RenderTexturedVoxel(position, GlobalVoxelMesh.size, textures, rl.White, visibleFaces)
	}
}
