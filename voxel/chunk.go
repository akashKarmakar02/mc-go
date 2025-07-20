package voxel

type Chunk struct {
	Blocks map[[3]int]Block
}

func NewChunk() *Chunk {
	return &Chunk{
		Blocks: make(map[[3]int]Block),
	}
}

func (c *Chunk) AddBlock(b Block) {
	pos := [3]int{int(b.Position.X), int(b.Position.Y), int(b.Position.Z)}
	c.Blocks[pos] = b
}

func (c *Chunk) HasBlock(x, y, z int) bool {
	_, exists := c.Blocks[[3]int{x, y, z}]
	return exists
}

func (c *Chunk) Render(atlas *TextureAtlas, mesh *VoxelMesh) {
	for _, block := range c.Blocks {
		x, y, z := int(block.Position.X), int(block.Position.Y), int(block.Position.Z)

		// Culling logic
		visibleFaces := []int{}
		if !c.HasBlock(x, y+1, z) {
			visibleFaces = append(visibleFaces, FaceTop)
		}
		if !c.HasBlock(x, y-1, z) {
			visibleFaces = append(visibleFaces, FaceBottom)
		}
		if !c.HasBlock(x, y, z+1) {
			visibleFaces = append(visibleFaces, FaceFront)
		}
		if !c.HasBlock(x, y, z-1) {
			visibleFaces = append(visibleFaces, FaceBack)
		}
		if !c.HasBlock(x-1, y, z) {
			visibleFaces = append(visibleFaces, FaceLeft)
		}
		if !c.HasBlock(x+1, y, z) {
			visibleFaces = append(visibleFaces, FaceRight)
		}

		atlas.RenderBlock(mesh, block, visibleFaces)
	}
}
