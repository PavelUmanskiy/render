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
	spaceShipSpeed       float32 = 5
	asteroidsAmount      int32   = 12
	asteroidMaxHP        int32   = 3
	asteroidRadius       float32 = 50.
	asteroidSpeed        float32 = .005
	asteroidSpawnZone            = 40
	asteroidCollisionDmg int32   = 1
	wHeight              float32 = 800
	wWidth               float32 = 800
)

func SpawnAsteroids(asteroids *[asteroidsAmount]Asteroid) {
	width := int32(wWidth)
	height := int32(wHeight)
	for i := range asteroidsAmount {
		asteroids[i] = Asteroid{
			Center:   rl.Vector2{X: float32(rl.GetRandomValue(0, width)), Y: float32(rl.GetRandomValue(0, height))},
			Rotation: rl.Vector2{X: float32(rl.GetRandomValue(-width, width)), Y: float32(rl.GetRandomValue(-height, height))},
			HP:       asteroidMaxHP,
			ID:       i,
		}
	}
}

func DrawAsteroids(asteroids *[asteroidsAmount]Asteroid) {
	for i := range asteroids {
		asteroid := asteroids[i]
		rl.DrawCircleLinesV(asteroid.Center, asteroidRadius, rl.White)
	}
}

func HandleControls() rl.Vector2 {
	desiredVector := rl.Vector2{}
	if rl.IsKeyPressed(rl.KeyDown) || rl.IsKeyPressed(rl.KeyS) || rl.IsKeyDown(rl.KeyDown) || rl.IsKeyDown(rl.KeyS) {
		// s.Center.Y = min(s.Center.Y+spaceShipSpeed, wHeight-spaceShipRadius)
		desiredVector.Y += spaceShipSpeed
	}
	if rl.IsKeyPressed(rl.KeyUp) || rl.IsKeyPressed(rl.KeyW) || rl.IsKeyDown(rl.KeyUp) || rl.IsKeyDown(rl.KeyW) {
		// s.Center.Y = max(s.Center.Y-spaceShipSpeed, spaceShipRadius)
		desiredVector.Y -= spaceShipSpeed
	}
	if rl.IsKeyPressed(rl.KeyLeft) || rl.IsKeyPressed(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft) || rl.IsKeyDown(rl.KeyA) {
		// s.Center.X = max(s.Center.X-spaceShipSpeed, spaceShipRadius)
		desiredVector.X -= spaceShipSpeed
	}
	if rl.IsKeyPressed(rl.KeyRight) || rl.IsKeyPressed(rl.KeyD) || rl.IsKeyDown(rl.KeyRight) || rl.IsKeyDown(rl.KeyD) {
		// s.Center.X = min(s.Center.X+spaceShipSpeed, wWidth-spaceShipRadius)
		desiredVector.X += spaceShipSpeed
	}
	return desiredVector
}

func ApplyWindowBounds(center *rl.Vector2, radius float32) rl.Vector2 {
	// Returns a Vector2 that shows if the window collision happened. If center is greater than max axis,
	// the value for that coordinate will be 1, if less than min axis, the value will be -1. 0 if no collision.
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

func MoveSpaceShip(s *SpaceShip, desiredVector rl.Vector2, asteroids *[asteroidsAmount]Asteroid) {
	// Applies desiredVector to the space ship and checks for collisions.
	oldCenter := s.Center
	s.Center = rl.Vector2Add(s.Center, desiredVector)
	ApplyWindowBounds(&s.Center, spaceShipRadius)
	for i := range asteroidsAmount {
		if rl.CheckCollisionCircles(s.Center, spaceShipRadius, asteroids[i].Center, asteroidRadius) {
			s.HP -= asteroidCollisionDmg
			s.Center = oldCenter
		}
	}
}

func MoveAsteroids(s *SpaceShip, asteroids *[asteroidsAmount]Asteroid) {
	for i := range asteroidsAmount {
		oldCenter := asteroids[i].Center
		asteroids[i].Center = rl.Vector2Add(
			asteroids[i].Center,
			rl.Vector2{X: asteroids[i].Rotation.X * asteroidSpeed, Y: asteroids[i].Rotation.Y * asteroidSpeed},
		)
		if rl.CheckCollisionCircles(s.Center, spaceShipRadius, asteroids[i].Center, asteroidRadius) {
			asteroids[i].Rotation = rl.Vector2{X: -asteroids[i].Rotation.X, Y: -asteroids[i].Rotation.Y}
			s.HP -= asteroidCollisionDmg
			asteroids[i].Center = oldCenter
		}
		windowCollision := ApplyWindowBounds(&asteroids[i].Center, asteroidRadius)
		if windowCollision.X != 0 {
			asteroids[i].Rotation.X = -asteroids[i].Rotation.X
		}
		if windowCollision.Y != 0 {
			asteroids[i].Rotation.Y = -asteroids[i].Rotation.Y
		}
	}
}

func Asteroids() {
	rl.InitWindow(int32(wWidth), int32(wHeight), "Asteroids")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	s := SpaceShip{
		Center: rl.Vector2{X: float32(wWidth) / 2, Y: float32(wHeight) / 2},
		HP:     spaceShipMaxHp,
	}
	asteroids := [asteroidsAmount]Asteroid{}
	SpawnAsteroids(&asteroids)
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.DrawText(fmt.Sprintf("HP: %d", s.HP), 0, 0, 24, rl.Yellow)
		MoveAsteroids(&s, &asteroids)
		DrawAsteroids(&asteroids)
		desiredVector := HandleControls()
		MoveSpaceShip(&s, desiredVector, &asteroids)
		rl.DrawCircleLinesV(s.Center, spaceShipRadius, rl.Yellow)

		rl.EndDrawing()
	}
}

func main() {
	Asteroids()
}
