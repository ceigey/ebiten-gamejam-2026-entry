package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"mygame/core"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	exampleimages "github.com/hajimehoshi/ebiten/v2/examples/resources/images"
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
	bgRaw, _, err := image.Decode(bytes.NewReader(exampleimages.Tile_png))
	bg = ebiten.NewImageFromImage(bgRaw)

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
	peekAhead := 50.0 // 50.0
	inertiaBias := 10.0
	cameraTargetX := state.Player.Position.X + math.Cos(state.Player.Rotation)*peekAhead + state.Player.Inertia.X*inertiaBias
	cameraTargetY := state.Player.Position.Y + math.Sin(state.Player.Rotation)*peekAhead + state.Player.Inertia.Y*inertiaBias
	cam.LookAt(cameraTargetX, cameraTargetY)

	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		cam.AddTrauma(1.0)
	}
	state.Player.Update(&state)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 77, G: 155, B: 230})
	drawTestBg(&state, screen)
	mx, my := ebiten.CursorPosition()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f, Mouse: %d, %d, Azim: %0.2f", ebiten.ActualTPS(), mx, my, state.Player.Rotation*57.2958))
	state.Player.Draw(state, screen)
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

// I'm not 100% thrilled with what's going on here
func drawTestBg(state *core.GameState, screen *ebiten.Image) {
	repeat := 20
	parallaxFactor := 0.75

	zoomFactor := state.Camera.ZoomFactor / 2

	w := float64(bg.Bounds().Dx()) * zoomFactor
	h := float64(bg.Bounds().Dy()) * zoomFactor

	offsetX := state.Camera.X * parallaxFactor * zoomFactor
	offsetY := state.Camera.Y * parallaxFactor * zoomFactor

	for j := range 10 {
		for i := range repeat {
			op := &ebiten.DrawImageOptions{}
			op.ColorScale.SetR(0.7)
			op.ColorScale.SetG(0.8)
			op.ColorScale.SetB(0.9)
			op.GeoM.Scale(zoomFactor, zoomFactor)
			op.GeoM.Translate(
				w*float64(i-repeat/2)-offsetX,
				h*float64(j-repeat/2)-offsetY,
			)
			screen.DrawImage(bg, op)
		}
	}
}
