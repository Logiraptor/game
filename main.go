package main

import (
	"image"
	"os"

	_ "image/png"

	"fmt"
	"time"

	"github.com/ByteArena/box2d"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1024, 600),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	pic, err := loadPicture("sprites/ground.png")
	if err != nil {
		panic(err)
	}

	ball, err := loadPicture("sprites/ball.png")
	if err != nil {
		panic(err)
	}

	sprite := pixel.NewSprite(pic, pic.Bounds())
	ballSprite := pixel.NewSprite(ball, ball.Bounds())

	win.Clear(colornames.Skyblue)

	ground := pixel.NewBatch(&pixel.TrianglesData{}, pic)
	for i := 0; i*16 < 1024; i++ {
		sprite.Draw(ground, pixel.IM.Moved(pixel.V(float64(i)*16+8, 8)))
	}


	gravity := box2d.MakeB2Vec2(0.0, -10.0)
	world := box2d.MakeB2World(gravity)

	// Ground
	{
		bd := box2d.MakeB2BodyDef()
		ground := world.CreateBody(&bd)

		shape := box2d.MakeB2EdgeShape()
		shape.Set(box2d.MakeB2Vec2(0, 16), box2d.MakeB2Vec2(1024, 16))
		ground.CreateFixture(&shape, 0.0)
	}

	// Circle character
	var ballBody *box2d.B2Body
	{
		bd := box2d.MakeB2BodyDef()
		bd.Position.Set(512, 300)
		bd.Type = box2d.B2BodyType.B2_dynamicBody
		bd.AllowSleep = false

		ballBody = world.CreateBody(&bd)

		shape := box2d.MakeB2CircleShape()
		shape.M_radius = 8

		fd := box2d.MakeB2FixtureDef()
		fd.Shape = &shape
		fd.Density = 1.0
		fd.Friction = 1.0
		ballBody.CreateFixtureFromDef(&fd)
	}

	// Circle character
	var ballBody2 *box2d.B2Body
	{
		bd := box2d.MakeB2BodyDef()
		bd.Position.Set(520, 330)
		bd.Type = box2d.B2BodyType.B2_dynamicBody
		bd.AllowSleep = false

		ballBody2 = world.CreateBody(&bd)

		shape := box2d.MakeB2CircleShape()
		shape.M_radius = 8

		fd := box2d.MakeB2FixtureDef()
		fd.Shape = &shape
		fd.Density = 1.0
		fd.Friction = 1.0
		ballBody2.CreateFixtureFromDef(&fd)
	}

	var (
		frames = 0
		second = time.Tick(time.Second)
	)

	timeStep := 1.0 / 60.0
	velocityIterations := 8
	positionIterations := 3

	for !win.Closed() {

		world.Step(timeStep, velocityIterations, positionIterations)

		win.Clear(colornames.Skyblue)

		ground.Draw(win)

		ballPos := ballBody.GetPosition()
		ballAngle := ballBody.GetAngle()

		ballMat := pixel.IM
		ballMat = ballMat.Rotated(pixel.ZV ,ballAngle)
		ballMat = ballMat.Moved(pixel.V(ballPos.X, ballPos.Y))
		ballSprite.Draw(win, ballMat)

		ballPos = ballBody2.GetPosition()
		ballAngle = ballBody2.GetAngle()

		ballMat = pixel.IM
		ballMat = ballMat.Rotated(pixel.ZV ,ballAngle)
		ballMat = ballMat.Moved(pixel.V(ballPos.X, ballPos.Y))
		ballSprite.Draw(win, ballMat)

		frames++

		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
