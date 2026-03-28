package main

import (
	"fmt"
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type IvState struct {
	IvMaxSec    float32
	IvCurrSec   float32
	IvCDMaxSec  float32
	IvCDCurrSec float32
	IvActive    bool
	IvCDActive  bool
}

type SpaceShip struct {
	Center   rl.Vector2
	Rotation rl.Vector2
	HP       int32
}

type Bullet struct {
	Base  rl.Vector2
	End   rl.Vector2
	Dir   rl.Vector2
	Alive bool
	ID    int32
}

type Asteroid struct {
	Center   rl.Vector2
	Rotation rl.Vector2
	HP       int32
	ID       int32
	Alive    bool
}

const (
	spaceShipMaxHp  int32   = 5
	spaceShipRadius float32 = 26
	spaceShipSpeed  float32 = 250
)

const (
	bulletsAmount int32   = 30
	bulletLength  float32 = 15
	bulletSpeed   float32 = 280
	bulletCD      float32 = 0.5
)

const (
	asteroidsAmount      int32   = 12
	asteroidMaxHP        int32   = 3
	asteroidRadius       float32 = 50
	asteroidSpeed        float32 = 80
	spaceShipSpawnZone   float32 = 40
	asteroidCollisionDmg int32   = 1
)

const (
	ivMaxSec    = 1.5
	ivCurrSec   = 1.5
	ivCDMaxSec  = 3.
	ivCDCurrSec = 3.
)

const (
	wHeight float32 = 1080
	wWidth  float32 = 1920
)

func SpawnAsteroids(s *SpaceShip, asteroids *[asteroidsAmount]Asteroid) {
	width := int32(wWidth)
	height := int32(wHeight)
	for i := range asteroidsAmount {
		if !asteroids[i].Alive {
			asteroids[i] = Asteroid{
				Center:   rl.Vector2{X: float32(rl.GetRandomValue(0, width)), Y: float32(rl.GetRandomValue(0, height))},
				Rotation: rl.Vector2{X: float32(rl.GetRandomValue(-width, width)), Y: float32(rl.GetRandomValue(-height, height))},
				HP:       asteroidMaxHP,
				ID:       i,
				Alive:    true,
			}
			for {
				if !rl.CheckCollisionCircles(s.Center, spaceShipRadius+spaceShipSpawnZone, asteroids[i].Center, asteroidRadius) {
					break
				}
				asteroids[i].Center = rl.Vector2{X: float32(rl.GetRandomValue(0, width)), Y: float32(rl.GetRandomValue(0, height))}
			}
		}
	}
}

func DrawAsteroids(asteroids *[asteroidsAmount]Asteroid) {
	for i := range asteroids {
		if asteroids[i].Alive {
			rl.DrawCircleLinesV(asteroids[i].Center, asteroidRadius, rl.White)
		}
	}
}

func DrawBullets(bullets []Bullet) {
	for i := range len(bullets) {
		if bullets[i].Alive {
			rl.DrawLineV(bullets[i].Base, bullets[i].End, rl.White)
		}
	}
}

func HandleControls() rl.Vector2 {
	direction := rl.Vector2{}
	if rl.IsKeyPressed(rl.KeyUp) || rl.IsKeyPressed(rl.KeyW) || rl.IsKeyDown(rl.KeyUp) || rl.IsKeyDown(rl.KeyW) {
		direction.Y -= 1
	}
	if rl.IsKeyPressed(rl.KeyDown) || rl.IsKeyPressed(rl.KeyS) || rl.IsKeyDown(rl.KeyDown) || rl.IsKeyDown(rl.KeyS) {
		direction.Y += 1
	}
	if rl.IsKeyPressed(rl.KeyLeft) || rl.IsKeyPressed(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft) || rl.IsKeyDown(rl.KeyA) {
		direction.X -= 1
	}
	if rl.IsKeyPressed(rl.KeyRight) || rl.IsKeyPressed(rl.KeyD) || rl.IsKeyDown(rl.KeyRight) || rl.IsKeyDown(rl.KeyD) {
		direction.X += 1
	}
	if direction.X != 0 || direction.Y != 0 {
		direction = rl.Vector2Normalize(direction)
	}
	return direction
}

func HandleMouse(s *SpaceShip) {
	delta := rl.Vector2Subtract(rl.GetMousePosition(), s.Center)
	if rl.Vector2Length(delta) > 0 {
		s.Rotation = rl.Vector2Normalize(delta)
	}
}

// Returns a Vector2 that shows if the window collision happened. If center is greater than max axis,
// the value for that coordinate will be 1, if less than min axis, the value will be -1. 0 if no collision.
func ApplyWindowBounds(center *rl.Vector2, radius float32) rl.Vector2 {
	windowCollision := rl.Vector2{}
	if center.X > wWidth-radius {
		center.X = wWidth - radius
		windowCollision.X = 1
	} else if center.X < radius {
		center.X = radius
		windowCollision.X = -1
	}
	if center.Y > wHeight-radius {
		center.Y = wHeight - radius
		windowCollision.Y = 1
	} else if center.Y < radius {
		center.Y = radius
		windowCollision.Y = -1
	}
	return windowCollision
}

// Applies direction and velocity to all the asteroids and checks for collisions.
func MoveAsteroids(s *SpaceShip, asteroids *[asteroidsAmount]Asteroid, dt float32) {
	for i := range asteroidsAmount {
		oldCenter := asteroids[i].Center
		direction := rl.Vector2Normalize(asteroids[i].Rotation)
		velocity := rl.Vector2Scale(direction, asteroidSpeed*dt)
		asteroids[i].Center = rl.Vector2Add(asteroids[i].Center, velocity)
		if asteroids[i].Alive && rl.CheckCollisionCircles(s.Center, spaceShipRadius, asteroids[i].Center, asteroidRadius) {
			asteroids[i].Rotation = rl.Vector2Scale(asteroids[i].Rotation, -1)
			s.HP -= asteroidCollisionDmg
			asteroids[i].Center = oldCenter
		}
		windowCollision := ApplyWindowBounds(&asteroids[i].Center, asteroidRadius)
		if windowCollision.X != 0 {
			asteroids[i].Rotation.X *= -1
		}
		if windowCollision.Y != 0 {
			asteroids[i].Rotation.Y *= -1
		}
	}
}

// Applies direction and velocity to the space ship and checks for collisions.
func MoveSpaceShip(s *SpaceShip, direction rl.Vector2, asteroids *[asteroidsAmount]Asteroid, dt float32) {
	oldCenter := s.Center
	velocity := rl.Vector2Scale(direction, spaceShipSpeed*dt)
	s.Center = rl.Vector2Add(s.Center, velocity)
	ApplyWindowBounds(&s.Center, spaceShipRadius)
	for i := range asteroidsAmount {
		if asteroids[i].Alive && rl.CheckCollisionCircles(s.Center, spaceShipRadius, asteroids[i].Center, asteroidRadius) {
			s.HP -= asteroidCollisionDmg
			s.Center = oldCenter
			break
		}
	}
}

func DrawSpaceShip(s *SpaceShip, invuln *IvState) {
	rl.DrawCircleLinesV(s.Center, spaceShipRadius, rl.Yellow)
	tipOffset := rl.Vector2Scale(s.Rotation, spaceShipRadius+10)
	tip := rl.Vector2Add(s.Center, tipOffset)
	backwardDir := rl.Vector2Scale(s.Rotation, -1)
	leftWingDir := rl.Vector2Rotate(backwardDir, -70)
	rightWingDir := rl.Vector2Rotate(backwardDir, 70)
	leftWing := rl.Vector2Add(tip, rl.Vector2Scale(leftWingDir, 15))
	rightWing := rl.Vector2Add(tip, rl.Vector2Scale(rightWingDir, 15))
	rl.DrawLineV(tip, leftWing, rl.Green)
	rl.DrawLineV(tip, rightWing, rl.Red)

	if invuln.IvActive {
		rl.DrawCircleLinesV(s.Center, spaceShipRadius*1.5, rl.Green)
	}
}

func HandleHP(state *IvState, prevHP, currHP int32, dt float32) (int32, rl.Color) {
	// iv active
	if state.IvActive && state.IvCurrSec > 0 {
		state.IvCurrSec -= dt
		// iv active but just ran out
	} else if state.IvActive && state.IvCurrSec <= 0 {
		state.IvActive = false
		state.IvCDActive = true
		state.IvCurrSec = ivMaxSec
		// iv on cd
	} else if state.IvCDActive && state.IvCDCurrSec > 0 {
		state.IvCDCurrSec -= dt
		// iv on cd but just ran out
	} else if state.IvCDActive && state.IvCDCurrSec <= 0 {
		state.IvCDActive = false
		state.IvCDCurrSec = ivCDMaxSec
		// iv ready
	} else if !state.IvCDActive && prevHP != currHP {
		state.IvActive = true
	}
	// Determine color
	// Green - iv ready
	// Yellow - iv active
	// Red - iv on cd
	var drawColor rl.Color = rl.Green
	if state.IvActive {
		drawColor = rl.Yellow
	} else if state.IvCDActive {
		drawColor = rl.Red
	}
	// HP logic
	resultHP := currHP
	if state.IvActive || prevHP == currHP {
		resultHP = prevHP
	}
	return resultHP, drawColor
}

func FireBullet(s *SpaceShip, bullets []Bullet, bulletTime *float32, dt float32) []Bullet {
	*bulletTime -= dt
	if *bulletTime > 0 {
		return bullets
	}
	direction := rl.Vector2Normalize(s.Rotation)
	bulletBase := rl.Vector2Add(s.Center, rl.Vector2Scale(direction, spaceShipRadius))
	bulletEnd := rl.Vector2Add(bulletBase, rl.Vector2Scale(direction, bulletLength))
	*bulletTime = bulletCD
	return append(bullets, Bullet{Base: bulletBase, End: bulletEnd, Dir: direction, Alive: true})
}

func MoveBullets(bullets []Bullet, asteroids *[asteroidsAmount]Asteroid, score *int32, dt float32) {
	for i := range len(bullets) {
		velocity := rl.Vector2Scale(bullets[i].Dir, bulletSpeed*dt)
		bullets[i].Base = rl.Vector2Add(bullets[i].Base, velocity)
		bullets[i].End = rl.Vector2Add(bullets[i].End, velocity)
		windowCollision := ApplyWindowBounds(&bullets[i].End, 1)
		if windowCollision.X != 0 || windowCollision.Y != 0 {
			bullets[i].Alive = false
		}
		for j := range asteroidsAmount {
			if asteroids[j].Alive && bullets[i].Alive && rl.CheckCollisionPointCircle(bullets[i].Base, asteroids[j].Center, asteroidRadius) {
				bullets[i].Alive = false
				asteroids[j].Alive = false
				*score++
			}
		}
	}
}

func Asteroids() {
	rl.InitWindow(int32(wWidth), int32(wHeight), "Asteroids")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)
	var score int32

	s := SpaceShip{
		Center:   rl.Vector2{X: float32(wWidth) / 2, Y: float32(wHeight) - float32(wHeight)/3},
		HP:       spaceShipMaxHp,
		Rotation: rl.NewVector2(1, 0),
	}
	invuln := IvState{
		IvMaxSec:    ivMaxSec,
		IvCurrSec:   ivCurrSec,
		IvCDMaxSec:  ivCDMaxSec,
		IvCDCurrSec: ivCDCurrSec,
		IvActive:    false,
		IvCDActive:  false,
	}
	asteroids := [asteroidsAmount]Asteroid{}
	var bulletTime float32 = bulletCD
	bullets := make([]Bullet, 0, 30)
	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()
		SpawnAsteroids(&s, &asteroids)
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.DrawText(fmt.Sprintf("HP: %d", s.HP), 0, 0, 24, rl.Yellow)
		prevHP := s.HP
		MoveAsteroids(&s, &asteroids, dt)
		DrawAsteroids(&asteroids)
		direction := HandleControls()
		MoveSpaceShip(&s, direction, &asteroids, dt)
		HandleMouse(&s)
		currHP := s.HP
		resultHP, invulnColor := HandleHP(&invuln, prevHP, currHP, dt)
		s.HP = resultHP

		bullets = FireBullet(&s, bullets, &bulletTime, dt)
		MoveBullets(bullets, &asteroids, &score, dt)
		bullets = slices.DeleteFunc(bullets, func(b Bullet) bool { return !b.Alive })
		DrawBullets(bullets)
		DrawSpaceShip(&s, &invuln)
		rl.DrawText(fmt.Sprintf("IV: %.2f", invuln.IvCurrSec), 0, 26, 24, invulnColor)
		rl.DrawText(fmt.Sprintf("CD: %.2f", invuln.IvCDCurrSec), 0, 52, 24, rl.Yellow)
		rl.DrawText(fmt.Sprintf("SC: %d", score), 0, 78, 24, rl.White)

		rl.EndDrawing()
	}
}

func main() {
	Asteroids()
}
