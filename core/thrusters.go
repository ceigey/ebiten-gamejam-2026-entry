package core

import (
	"image"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
)

type Thruster struct {
	Position Vec2
	Angle    float64
}

func (player *Player) Thrusters() [8]Thruster {
	leftRev := Thruster{
		Position: Vec2{X: 19, Y: 20},
		Angle:    0.0,
	}
	leftA := Thruster{
		Position: Vec2{X: 11, Y: 38},
		Angle:    270.0,
	}
	leftB := Thruster{
		Position: Vec2{X: 11, Y: 49},
		Angle:    270.0,
	}
	leftMain := Thruster{
		Position: Vec2{X: 19, Y: 64},
		Angle:    180.0,
	}
	rightRev := Thruster{
		Position: Vec2{X: 45, Y: 20},
		Angle:    0.0,
	}
	rightA := Thruster{
		Position: Vec2{X: 53, Y: 38},
		Angle:    90.0,
	}
	rightB := Thruster{
		Position: Vec2{X: 53, Y: 49},
		Angle:    90.0,
	}
	rightMain := Thruster{
		Position: Vec2{X: 45, Y: 64},
		Angle:    180.0,
	}

	thrusters := [8]Thruster{
		leftRev,
		leftA,
		leftB,
		leftMain,
		rightRev,
		rightA,
		rightB,
		rightMain,
	}

	return thrusters
}

func (player *Player) DrawThrusterPlumes(state *GameState, screen *ebiten.Image) {
	thrusters := player.Thrusters()

	for _, outlet := range thrusters {
		outlet.DrawPlume(player, state, screen)
	}

}

func (thruster *Thruster) AngleRadians() float64 {
	return (thruster.Angle - 90) * (math.Pi / 180.0)
}

func (thruster *Thruster) PositionRelativeToPlayer(player *Player) Vec2 {
	localX := thruster.Position.X - 32
	localY := thruster.Position.Y - 32
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

func (thruster *Thruster) AbsoluteRotation(player *Player) float64 {
	return thruster.AngleRadians() + player.Rotation
}

func (thruster *Thruster) PreparePlumeGeometry(player *Player, thrustFactor float64) (*ebiten.Image, *ebiten.DrawImageOptions) {
	jitterFactor := 0.75 + rand.Float64()/4
	plumeFactor := thrustFactor * jitterFactor

	adjustedPosition := thruster.PositionRelativeToPlayer(player)
	absoluteRotation := thruster.AbsoluteRotation(player)

	subimg := player.Image.SubImage(image.Rect(0, 80, 16, 96)).(*ebiten.Image)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(subimg.Bounds().Dx())/2, 0)
	op.GeoM.Scale(1.0, plumeFactor)
	op.GeoM.Rotate(absoluteRotation)
	op.GeoM.Translate(adjustedPosition.Coords())
	return subimg, op
}

// func (outlet *Thruster) Calculate

func (thruster *Thruster) DrawPlume(player *Player, state *GameState, screen *ebiten.Image) {
	reverseEngines := player.Inertia.Magnitude() > 1 && player.IsBreaking
	if !reverseEngines && player.PositionDelta.Magnitude() < 0.01 {
		return
	}

	// Damn angular adjustments needed again because "North is West"
	angleInRadians := thruster.AngleRadians()
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

	subimg, op := thruster.PreparePlumeGeometry(player, thrustFactor*inertiaFactor)
	state.Camera.Draw(subimg, op, screen)
}
