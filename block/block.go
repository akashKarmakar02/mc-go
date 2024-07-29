package block

import rl "github.com/gen2brain/raylib-go/raylib"

type Block interface {
	Render(rl.Camera)
	Unload()
}
