package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type SpaceShip struct {
	Center   rl.Vector2
	Rotation rl.Vector2
	HP       int32
}

type Asteroid struct {
	Center   rl.Vector2
	Rotation rl.Vector2
	HP       int32
	ID       int32
}

const (
	spaceShipMaxHp       int32   = 5
	spaceShipRadius      float32 = 26.
	spaceShipSpeed       float32 = 250
	asteroidsAmount      int32   = 12
	asteroidMaxHP        int32   = 3
	asteroidRadius       float32 = 50.
	asteroidSpeed        float32 = 80
	spaceShipSpawnZone   float32 = 40.
	asteroidCollisionDmg int32   = 1
	wHeight              float32 = 1080
	wWidth               float32 = 1920
)

func SpawnAsteroids(s *SpaceShip, asteroids *[asteroidsAmount]Asteroid) {
	width := int32(wWidth)
	height := int32(wHeight)
	for i := range asteroidsAmount {
		asteroids[i] = Asteroid{
			Center:   rl.Vector2{X: float32(rl.GetRandomValue(0, width)), Y: float32(rl.GetRandomValue(0, height))},
			Rotation: rl.Vector2{X: float32(rl.GetRandomValue(-width, width)), Y: float32(rl.GetRandomValue(-height, height))},
			HP:       asteroidMaxHP,
			ID:       i,
		}
		for {
			if !rl.CheckCollisionCircles(s.Center, spaceShipRadius+spaceShipSpawnZone, asteroids[i].Center, asteroidRadius) {
				break
			}
			asteroids[i].Center = rl.Vector2{X: float32(rl.GetRandomValue(0, width)), Y: float32(rl.GetRandomValue(0, height))}
		}
	}
}

func DrawAsteroids(asteroids *[asteroidsAmount]Asteroid) {
	for i := range asteroids {
		rl.DrawCircleLinesV(asteroids[i].Center, asteroidRadius, rl.White)
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

// Applies direction and velocity to the space ship and checks for collisions.
func MoveSpaceShip(s *SpaceShip, direction rl.Vector2, asteroids *[asteroidsAmount]Asteroid, dt float32) {
	oldCenter := s.Center
	velocity := rl.Vector2Scale(direction, spaceShipSpeed*dt)
	s.Center = rl.Vector2Add(s.Center, velocity)
	ApplyWindowBounds(&s.Center, spaceShipRadius)
	for i := range asteroidsAmount {
		if rl.CheckCollisionCircles(s.Center, spaceShipRadius, asteroids[i].Center, asteroidRadius) {
			s.HP -= asteroidCollisionDmg
			s.Center = oldCenter
			break
		}
	}
}

// Applies direction and velocity to all the asteroids and checks for collisions.
func MoveAsteroids(s *SpaceShip, asteroids *[asteroidsAmount]Asteroid, dt float32) {
	for i := range asteroidsAmount {
		oldCenter := asteroids[i].Center
		direction := rl.Vector2Normalize(asteroids[i].Rotation)
		velocity := rl.Vector2Scale(direction, asteroidSpeed*dt)
		asteroids[i].Center = rl.Vector2Add(asteroids[i].Center, velocity)
		if rl.CheckCollisionCircles(s.Center, spaceShipRadius, asteroids[i].Center, asteroidRadius) {
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

func Asteroids() {
	rl.InitWindow(int32(wWidth), int32(wHeight), "Asteroids")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	s := SpaceShip{
		Center: rl.Vector2{X: float32(wWidth) / 2, Y: float32(wHeight) - float32(wHeight)/3},
		HP:     spaceShipMaxHp,
	}
	asteroids := [asteroidsAmount]Asteroid{}
	SpawnAsteroids(&s, &asteroids)
	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.DrawText(fmt.Sprintf("HP: %d", s.HP), 0, 0, 24, rl.Yellow)
		MoveAsteroids(&s, &asteroids, dt)
		DrawAsteroids(&asteroids)
		direction := HandleControls()
		MoveSpaceShip(&s, direction, &asteroids, dt)
		rl.DrawCircleLinesV(s.Center, spaceShipRadius, rl.Yellow)

		rl.EndDrawing()
	}
}

func main() {
	Asteroids()
}
