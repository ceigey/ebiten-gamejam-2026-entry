package core

import (
	"image"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
)

type EngineOutlet struct {
	Position Vec2
	Angle    float64
}

func (player *Player) EngineOutlets() [8]EngineOutlet {
	leftRev := EngineOutlet{
		Position: Vec2{X: 19, Y: 20},
		Angle:    0.0,
	}
	leftA := EngineOutlet{
		Position: Vec2{X: 11, Y: 38},
		Angle:    270.0,
	}
	leftB := EngineOutlet{
		Position: Vec2{X: 11, Y: 49},
		Angle:    270.0,
	}
	leftMain := EngineOutlet{
		Position: Vec2{X: 19, Y: 64},
		Angle:    180.0,
	}
	rightRev := EngineOutlet{
		Position: Vec2{X: 45, Y: 20},
		Angle:    0.0,
	}
	rightA := EngineOutlet{
		Position: Vec2{X: 53, Y: 38},
		Angle:    90.0,
	}
	rightB := EngineOutlet{
		Position: Vec2{X: 53, Y: 49},
		Angle:    90.0,
	}
	rightMain := EngineOutlet{
		Position: Vec2{X: 45, Y: 64},
		Angle:    180.0,
	}

	engineOutlets := [8]EngineOutlet{
		leftRev,
		leftA,
		leftB,
		leftMain,
		rightRev,
		rightA,
		rightB,
		rightMain,
	}

	return engineOutlets
}

func (player *Player) drawEnginePlumes(state *GameState, screen *ebiten.Image) {
	engineOutlets := player.EngineOutlets()

	for _, outlet := range engineOutlets {
		outlet.DrawPlume(player, state, screen)
	}

}

func (outlet *EngineOutlet) AngleRadians() float64 {
	return (outlet.Angle - 90) * (math.Pi / 180.0)
}

func (outlet *EngineOutlet) PositionRelativeToPlayer(player *Player) Vec2 {
	localX := outlet.Position.X - 32
	localY := outlet.Position.Y - 32
	// I had to look up this trigonometry, I need to study harder
	rotatedX := localX*math.Cos(player.AdjustedRotation()) - localY*math.Sin(player.AdjustedRotation())
	rotatedY := localX*math.Sin(player.AdjustedRotation()) + localY*math.Cos(player.AdjustedRotation())

	finalX := player.Position.X + rotatedX
	finalY := player.Position.Y + rotatedY
	return Vec2{
		X: finalX,
		Y: finalY,
	}
}

func (outlet *EngineOutlet) AbsoluteRotation(player *Player) float64 {
	return outlet.AngleRadians() + player.Rotation
}

func (outlet *EngineOutlet) PreparePlumeGeometry(player *Player, thrustFactor float64) (*ebiten.Image, *ebiten.DrawImageOptions) {
	jitterFactor := 0.75 + rand.Float64()/4
	plumeFactor := thrustFactor * jitterFactor

	adjustedPosition := outlet.PositionRelativeToPlayer(player)
	absoluteRotation := outlet.AbsoluteRotation(player)

	subimg := player.Image.SubImage(image.Rect(0, 80, 16, 96)).(*ebiten.Image)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(subimg.Bounds().Dx())/2, 0)
	op.GeoM.Scale(1.0, plumeFactor)
	op.GeoM.Rotate(absoluteRotation)
	op.GeoM.Translate(adjustedPosition.Coords())
	return subimg, op
}

func (outlet *EngineOutlet) DrawPlume(player *Player, state *GameState, screen *ebiten.Image) {
	reverseEngines := player.Inertia.Magnitude() > 1 && player.IsBreaking
	if !reverseEngines && player.PositionDelta.Magnitude() < 0.01 {
		return
	}

	// Damn angular adjustments needed again because "North is West"
	angleInRadians := outlet.AngleRadians()
	angleForComparison := angleInRadians
	if reverseEngines {
		angleForComparison -= math.Pi
	}

	thrustIndicator := player.PositionDelta
	if player.PositionDelta.Magnitude() < 0.01 {
		thrustIndicator = player.Inertia
	}

	expectedThrustAngleFromDelta := math.Atan2(thrustIndicator.Y, thrustIndicator.X) - math.Pi

	angleRelativeToPlayer := angleForComparison + player.AdjustedRotation()
	absoluteDelta := math.Abs(AngleDifferenceRadians(angleRelativeToPlayer, expectedThrustAngleFromDelta))
	similarity := math.Max(0, 1-absoluteDelta*2/math.Pi)

	thrustFactor := similarity
	inertiaFactor := player.Inertia.Magnitude() / 10
	if reverseEngines {
		inertiaFactor = 1.5
	}

	subimg, op := outlet.PreparePlumeGeometry(player, thrustFactor*inertiaFactor)
	state.Camera.Draw(subimg, op, screen)
}
