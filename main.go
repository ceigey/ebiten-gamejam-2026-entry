package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"mygame/core"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/setanarut/kamera/v2"
)

var img *ebiten.Image
var bg *ebiten.Image

var state core.GameState

var cam *kamera.Camera

func init() {
	var err error

	img, _, err = ebitenutil.NewImageFromFile("sprites/gunshipgame.png")
	if err != nil {
		log.Fatal(err)
	}

	// https://ebitengine.org/en/examples/infinitescroll.html
	// bgRaw, _, err := image.Decode(bytes.NewReader(exampleimages.Tile_png))
	bg, _, err = ebitenutil.NewImageFromFile("3rdparty/scigho-water-plus.png")
	if err != nil {
		log.Fatal(err)
	}
	// bg = ebiten.NewImageFromImage(bgTileset)

	state = core.GameState{
		Player: core.NewPlayer(img),
	}

	cam = kamera.NewCamera(state.Player.Position.X, state.Player.Position.Y, 640, 480)
	// Copied from Kamera library's example
	cam.ShakeEnabled = true
	cam.SmoothType = kamera.None
	// cam.SmoothOptions.SmoothDampTimeX = 0.1
	// cam.SmoothOptions.SmoothDampTimeY = 0.1
	state.Camera = cam
}

type Game struct{}

func (g *Game) Update() error {
	cameraBias := state.Player.LookAheadInertiaBias()
	cameraTargetX := state.Player.Position.X + cameraBias.X
	cameraTargetY := state.Player.Position.Y + cameraBias.Y
	cam.LookAt(cameraTargetX, cameraTargetY)

	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		cam.AddTrauma(1.0)
	}
	state.Player.Update(&state)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 77, G: 155, B: 230})
	// drawTestBg(&state, screen)
	mx, my := ebiten.CursorPosition()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f, Mouse: %d, %d, Azim: %0.2f", ebiten.ActualTPS(), mx, my, state.Player.Rotation*57.2958))
	state.Player.Draw(&state, screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Render an image")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

// Draws a test background for the game because I'm lazy
func drawTestBg(state *core.GameState, screen *ebiten.Image) {
	ticks := ebiten.Tick()
	if ticks%60 == 0 {
		if state.BackgroundCellFlip == 1 {
			state.BackgroundCellFlip = 0
		} else {
			state.BackgroundCellFlip = 1
		}
	}
	cells := []*ebiten.Image{
		// bg.SubImage(image.Rect(0, 96, 0+16, 96+16)).(*ebiten.Image),
		// bg.SubImage(image.Rect(112, 0, 112+16*2, 0+16)).(*ebiten.Image),
		// bg.SubImage(image.Rect(112, 32, 112+16, 32+16)).(*ebiten.Image),
		// bg.SubImage(image.Rect(96, 80, 96+16, 80+16)).(*ebiten.Image),
		// bg.SubImage(image.Rect(144, 112, 144+16, 112+16)).(*ebiten.Image),
		bg.SubImage(image.Rect(32, 32, 80, 48)).(*ebiten.Image),
		// bg.SubImage(image.Rect(112+16, 32, 112+32, 32+16)).(*ebiten.Image),
	}
	repeat := 20
	parallaxFactor := 0.75
	scale := 10.0

	zoomFactor := state.Camera.ZoomFactor / 2

	w := 16 * 1 * zoomFactor //float64(bg.Bounds().Dx()) * zoomFactor
	h := 16 * zoomFactor     //float64(bg.Bounds().Dy()) * zoomFactor

	offsetX := state.Camera.X * parallaxFactor * zoomFactor
	offsetY := state.Camera.Y * parallaxFactor * zoomFactor

	// fmt.Printf("%0.2f %0.2f\n", offsetX, offsetY)

	for j := range repeat {
		for i := range repeat {
			cell := cells[(j+i+state.BackgroundCellFlip)%1]
			op := &ebiten.DrawImageOptions{}
			// op.ColorScale.SetR(0.4)
			// op.ColorScale.SetG(0.7)
			// op.ColorScale.SetB(1.0)
			op.GeoM.Scale(scale, scale)
			op.GeoM.Scale(zoomFactor, zoomFactor)
			op.GeoM.Translate(
				w*scale*float64(i)-offsetX,
				h*scale*float64(j)-offsetY,
			)
			screen.DrawImage(cell, op)
		}
	}
}
