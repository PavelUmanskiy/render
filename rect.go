package main

import (
	"image/color"
	"time"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

func Rectangle() {
	var (
		wHeight int32 = 600
		wLength int32 = 600
	)
	raylib.InitWindow(wHeight, wLength, "My Window")
	defer raylib.CloseWindow()
	raylib.SetTargetFPS(60)

	var (
		rectPosX int32 = raylib.GetRandomValue(0, wLength)
		rectPosY int32 = raylib.GetRandomValue(0, wHeight)
		rectH    int32 = 50
		rectW    int32 = 50
		rectStep int32 = 10
	)
	rectColor := color.RGBA{A: 255}

	yDown := true
	xLeft := true

	for !raylib.WindowShouldClose() {
		yHit := true
		xHit := true

		raylib.BeginDrawing()
		raylib.ClearBackground(raylib.RayWhite)
		raylib.DrawRectangle(rectPosX, rectPosY, rectW, rectH, rectColor)

		if raylib.IsMouseButtonPressed(raylib.MouseButtonLeft) {
			time.Sleep(time.Second * 3)
		}

		if yDown && rectPosY+rectH >= wHeight {
			yDown = false
		} else if !yDown && rectPosY <= 0 {
			yDown = true
		} else {
			yHit = false
		}
		if yDown {
			rectPosY = (rectPosY + rectStep) % wHeight
		} else {
			rectPosY = (rectPosY - rectStep) % wHeight
		}

		if xLeft && rectPosX+rectW >= wLength {
			xLeft = false
		} else if !xLeft && rectPosX <= 0 {
			xLeft = true
		} else {
			xHit = false
		}
		if xLeft {
			rectPosX = (rectPosX + rectStep) % wLength
		} else {
			rectPosX = (rectPosX - rectStep) % wLength
		}

		raylib.EndDrawing()

		if xHit || yHit {
			rectColor = color.RGBA{
				R: uint8(raylib.GetRandomValue(0, 255)),
				G: uint8(raylib.GetRandomValue(0, 255)),
				B: uint8(raylib.GetRandomValue(0, 255)),
				A: 255,
			}
		}
	}
}
