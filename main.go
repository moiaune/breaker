package main

import (
	"fmt"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type direction uint8

const (
	DIR_UP = iota
	DIR_RIGHT
	DIR_DOWN
	DIR_LEFT
)

type paddle struct {
	width  int32
	height int32
	pos    pos
}

type ball struct {
	radius float32
	dirX   direction
	dirY   direction
	pos    pos
}

type brick struct {
	width  int32
	height int32
	pos    pos
	color  rl.Color
}

type pos struct {
	x int32
	y int32
}

const (
	winW        int32 = 800
	winH        int32 = 600
	paddleW     int32 = 120
	paddleH     int32 = 10
	brickW      int32 = 50
	brickH      int32 = 30
	bricksPrRow int32 = 10
	brickPad    int32 = 2
)

var (
	gameOver   bool       = false
	points     uint32     = 0
	colors     []rl.Color = []rl.Color{rl.Red, rl.Blue, rl.Green, rl.Yellow, rl.Purple}
	startTimer time.Time  = time.Now()
)

func main() {
	rl.InitWindow(winW, winH, "Breaker")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	b := resetBall()
	p := resetPaddle()
	bricks := generateBricks()

	for !rl.WindowShouldClose() {
		// -----------------------------------------------------------------------------
		//     - HANDLE KEY EVENTS -
		// -----------------------------------------------------------------------------
		if gameOver && rl.IsKeyDown(rl.KeySpace) {
			p = resetPaddle()
			b = resetBall()
			bricks = generateBricks()
			startTimer = time.Now()
			points = 0
			gameOver = false
		}

		if rl.IsKeyDown(rl.KeyA) && p.pos.x > 0 {
			p.pos.x -= 5
		}

		if rl.IsKeyDown(rl.KeyD) && (p.pos.x+p.width) < winW {
			p.pos.x += 5
		}

		// -----------------------------------------------------------------------------
		//     - DETECT PADDLE COLLISIONS -
		// -----------------------------------------------------------------------------
		if p.pos.x < 0 {
			p.pos.x = 0
		}

		if (p.pos.x + p.width) > winW {
			p.pos.x = winW - p.width
		}

		// -----------------------------------------------------------------------------
		//     - DETECT BALL COLLISIONS -
		// -----------------------------------------------------------------------------
		if (b.pos.x - int32(b.radius)) < 0 {
			b.pos.x = int32(b.radius)
			b.dirX = DIR_RIGHT
		}

		if (b.pos.x + int32(b.radius)) > winW {
			b.dirX = DIR_LEFT
		}

		if (b.pos.y - int32(b.radius)) < 0 {
			b.pos.y = int32(b.radius)
			b.dirY = DIR_DOWN
		}

		if b.pos.y > p.pos.y &&
			(b.pos.x > p.pos.x &&
				b.pos.x < (p.pos.x+p.width)) {
			b.pos.y = p.pos.y - int32(b.radius)
			b.dirY = DIR_UP
		}

		if b.pos.y > winH {
			gameOver = true
		}

		// -----------------------------------------------------------------------------
		//     - DETECT BALL COLLISION WITH BRICKS -
		// -----------------------------------------------------------------------------
		for idx, brick := range bricks {
			if b.pos.y < (brick.pos.y+brick.height) &&
				b.pos.y > brick.pos.y &&
				b.pos.x > brick.pos.x &&
				b.pos.x < (brick.pos.x+brick.width) {
				b.dirY = DIR_DOWN
				bricks = append(bricks[:idx], bricks[idx+1:]...)
				points++
			}
		}

		// -----------------------------------------------------------------------------
		//     - UPDATE BALL DIRECTIONS -
		// -----------------------------------------------------------------------------
		if b.dirY == DIR_DOWN {
			b.pos.y += 2
		}

		if b.dirY == DIR_UP {
			b.pos.y -= 2
		}

		if b.dirX == DIR_RIGHT {
			b.pos.x += 2
		}

		if b.dirX == DIR_LEFT {
			b.pos.x -= 2
		}

		// -----------------------------------------------------------------------------
		//     - CALCULATE DURATION -
		// -----------------------------------------------------------------------------
		now := time.Now()
		dur := now.Sub(startTimer).Round(time.Second)

		// -----------------------------------------------------------------------------
		//     - DRAW ENTITIES -
		// -----------------------------------------------------------------------------
		rl.BeginDrawing()

		if gameOver {
			rl.DrawText("GAME OVER", winW/2-100, winH/2, 36, rl.White)
		} else {
			rl.DrawText(fmt.Sprintf("Points: %d", points), 10, 10, 14, rl.White)
			rl.DrawText(fmt.Sprintf("Time: %s", dur.String()), 10, 30, 14, rl.White)
			rl.ClearBackground(rl.Black)
			rl.DrawRectangle(p.pos.x, p.pos.y, p.width, p.height, rl.White)
			rl.DrawCircle(b.pos.x, b.pos.y, b.radius, rl.White)
			for _, brick := range bricks {
				rl.DrawRectangle(brick.pos.x,
					brick.pos.y,
					brick.width,
					brick.height,
					brick.color)
			}
		}

		rl.EndDrawing()
	}
}

func generateBricks() []brick {
	var i int32
	var startY int32 = 50
	var startX int32 = (winW - (brickW * bricksPrRow)) / 2
	var bricks []brick
	for i < (bricksPrRow * 5) {
		if i > 0 && i%bricksPrRow == 0 {
			startY += brickH + brickPad
		}

		randColor := rl.GetRandomValue(0, int32(len(colors)-1))

		brick := brick{
			width:  brickW,
			height: brickH,
			color:  colors[randColor],
			pos: pos{
				startX + (i%bricksPrRow)*(brickW+brickPad),
				startY,
			},
		}

		bricks = append(bricks, brick)

		i++
	}

	return bricks
}

func resetPaddle() paddle {
	p := paddle{
		width:  paddleW,
		height: paddleH,
		pos: pos{
			(winW - paddleW) / 2,
			winH - (paddleH * 3),
		},
	}

	return p
}

func resetBall() ball {
	b := ball{
		radius: 5.0,
		dirX:   DIR_RIGHT,
		dirY:   DIR_DOWN,
		pos: pos{
			winW / 2,
			winH / 2,
		},
	}

	return b
}
