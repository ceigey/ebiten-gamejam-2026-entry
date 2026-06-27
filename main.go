package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"mygame/core"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var img *ebiten.Image

var gameState core.GameState

func init() {
	var err error

	img, _, err = ebitenutil.NewImageFromFile("sprites/gunshipgame.png")
	if err != nil {
		log.Fatal(err)
	}

	gameState = core.GameState{
		Player: core.NewPlayer(img),
	}
}

type Game struct{}

func (g *Game) Update() error {
	gameState.Player.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 77, G: 155, B: 230})
	mx, my := ebiten.CursorPosition()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f, Mouse: %d, %d, Azim: %0.2f", ebiten.ActualTPS(), mx, my, gameState.Player.Rotation*57.2958))
	gameState.Player.Draw(screen)
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
