package core

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type PlayerBullet struct {
	Position      Vec2
	PositionDelta Vec2
	Type          int
	BaseDamage    float32
}

type PlayerGun struct {
	Timer float32
}

type Player struct {
	Position        Vec2
	Inertia         Vec2
	PositionDelta   Vec2
	MachineGun      PlayerGun
	Cannon          PlayerGun
	CannonFireTimer float32
	Bullets         []PlayerBullet
	Health          float64
	Rotation        float64 // Theta in radians?
	Image           *ebiten.Image
	DragFactor      float64
	ThrusterPower   float64
	IsBreaking      bool
}

func NewPlayer(image *ebiten.Image) Player {
	return Player{
		Position:      Vec2{X: 360, Y: 180},
		Image:         image,
		DragFactor:    0.95,
		ThrusterPower: 0.75,
	}
}

func (player *Player) Update(state *GameState) {
	player.Position.X += player.Inertia.X
	player.Position.Y += player.Inertia.Y

	inputvec := Vec2{0, 0}
	var breakingFactor float64 = 1
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		inputvec.Y -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		inputvec.Y += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		inputvec.X -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		inputvec.X += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		breakingFactor = 0.75
		player.IsBreaking = true
	} else {
		player.IsBreaking = false
	}

	normalized := inputvec.Normalize()
	player.PositionDelta = normalized

	player.Inertia.X += (normalized.X * player.ThrusterPower)
	player.Inertia.Y += (normalized.Y * player.ThrusterPower)
	player.Inertia.X *= player.DragFactor * breakingFactor
	player.Inertia.Y *= player.DragFactor * breakingFactor

	// Solution for Movement or camera jitter when numbers get infinitessimally small
	if math.Abs(player.Inertia.X) < 0.1 {
		player.Inertia.X = 0
	}
	if math.Abs(player.Inertia.Y) < 0.1 {
		player.Inertia.Y = 0
	}

	// cameraBias := player.LookAheadInertiaBias()
	mxi, myi := ebiten.CursorPosition()
	mx, my := float64(mxi), float64(myi)
	mx += state.Camera.X
	my += state.Camera.Y

	targetRotation := math.Atan2(my-player.Position.Y, mx-player.Position.X)

	difference := AngleDifferenceRadians(player.Rotation, targetRotation)

	maxRotationPerTick := math.Pi / 32
	if math.Abs(difference) <= maxRotationPerTick {
		player.Rotation = targetRotation
	} else if difference > 0 {
		player.Rotation += maxRotationPerTick
	} else {
		player.Rotation -= maxRotationPerTick
	}
}

func (player *Player) AdjustedRotation() float64 {
	return player.Rotation + math.Pi/2
}

func (player *Player) Draw(state *GameState, screen *ebiten.Image) {

	camera := state.Camera

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-32, -32)
	subimg := player.Image.SubImage(image.Rect(0, 0, 64, 64))

	spriteAdjustedRotation := player.AdjustedRotation()
	op.GeoM.Rotate(spriteAdjustedRotation)
	op.GeoM.Translate(player.Position.X, player.Position.Y)
	// op.GeoM.
	// screen.DrawImage(img, op)
	// screen.DrawImage(subimg.(*ebiten.Image), op)
	camera.Draw(subimg.(*ebiten.Image), op, screen)
	// Drawing sight cone
	// https://www.bbc.co.uk/bitesize/articles/zyjtfdm#z27hfdm
	// https://www.omnicalculator.com/math/right-triangle-side-angle
	thetaRight := math.Mod(player.Rotation+math.Pi/32, math.Pi*2)
	hypotenuseRight := 1000.0
	adjacentRight := hypotenuseRight * math.Cos(thetaRight) // x
	oppositeRight := hypotenuseRight * math.Sin(thetaRight) // y

	thetaLeft := math.Mod(player.Rotation-math.Pi/32, math.Pi*2)
	hypotenuseLeft := 1000.0
	adjacentLeft := hypotenuseLeft * math.Cos(thetaLeft)
	oppositeLeft := hypotenuseLeft * math.Sin(thetaLeft)

	// https://ebitengine.org/en/examples/vector.html
	var path vector.Path

	px, py := camera.ApplyCameraTransformToPoint(player.Position.X, player.Position.Y)
	path.MoveTo(float32(px), float32(py))
	path.LineTo(float32(px+adjacentLeft), float32(py+oppositeLeft))
	path.LineTo(float32(px+adjacentRight), float32(py+oppositeRight))
	path.Close()

	coneOp := &vector.DrawPathOptions{}

	coneOp.AntiAlias = true
	coneOp.ColorScale.ScaleWithColor(color.NRGBA{0x22, 0xff, 0xaa, 0x22})
	coneOp.Blend = ebiten.BlendSourceOver
	vector.FillPath(screen, &path, nil, coneOp)
	// vector.StrokeLine(screen, float32(player.Position.X), float32(player.Position.Y))

	zoomFactor := 1 - player.Inertia.Magnitude()/100
	state.Camera.ZoomFactor = zoomFactor

	player.drawEnginePlumes(state, screen)
}

// https://stackoverflow.com/a/28037434
func AngleDifferenceRadians(angle1 float64, angle2 float64) float64 {
	half := math.Pi
	full := math.Pi * 2
	difference := math.Mod(angle2-angle1+half, full) - half

	if difference < -1*half {
		return difference + full
	} else {
		return difference
	}
}

func (player *Player) LookAheadInertiaBias() Vec2 {
	peekAhead := 50.0 // 50.0
	inertiaBias := 10.0
	return Vec2{
		X: math.Cos(player.Rotation)*peekAhead + player.Inertia.X*inertiaBias,
		Y: math.Sin(player.Rotation)*peekAhead + player.Inertia.Y*inertiaBias,
	}
}

func (player *Player) drawHitbox(state *GameState, screen *ebiten.Image) {
	hitbox := ebiten.NewImage(32, 24)
	hbOp := &ebiten.DrawImageOptions{}
	hbOp.GeoM.Translate(-16, 0)
	hbOp.GeoM.Rotate(player.AdjustedRotation())

	hbOp.GeoM.Translate(player.Position.X, player.Position.Y)
	hitbox.Fill(color.RGBA{255, 0, 0, 255})
	state.Camera.Draw(hitbox, hbOp, screen)
}
