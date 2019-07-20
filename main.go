package main

import (
	"image"
	"os"

	_ "image/png"

	"fmt"
	"math"
	"time"

	"github.com/ByteArena/box2d"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type Body = *box2d.B2Body

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
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	ball, err := loadPicture("sprites/ball.png")
	if err != nil {
		panic(err)
	}

	ballSprite := pixel.NewSprite(ball, ball.Bounds())

	win.Clear(colornames.Skyblue)

	game := NewGame()

	var (
		frames = 0
		second = time.Tick(time.Second)
	)
	createGround(game)

	ballSpriteIndex := game.AddSprite(ballSprite)

	player := createBox(game, pixel.V(100, 40), pixel.V(1.5, 3))
	for i := 5; i < 75; i += 4 {
		for j := 5; j < 75; j += 4 {
			createBall(game, pixel.V(float64(i), float64(j)), pixel.V(1, 1))
		}
	}

	workChan := game.StartPhysics()

	for !win.Closed() {

		if win.Pressed(pixelgl.KeyA) {
			player.ApplyForceToCenter(box2d.MakeB2Vec2(-100, 0), true)
		}

		if win.Pressed(pixelgl.KeyD) {
			player.ApplyForceToCenter(box2d.MakeB2Vec2(100, 0), true)
		}

		if win.JustPressed(pixelgl.KeySpace) {
			player.ApplyLinearImpulseToCenter(box2d.MakeB2Vec2(0, 100), true)
		}

		if win.JustPressed(pixelgl.KeyF) {
			workChan(func() {
				explode(game, player, 3, ballSpriteIndex)
			})
		}

		playerPos := player.GetPosition()
		win.SetMatrix(
			pixel.IM.
				Moved(pixel.V(-playerPos.X, -playerPos.Y)).
				Scaled(pixel.ZV, physicsScale*2).
				Moved(win.Bounds().Center()),
		)
		win.Clear(colornames.Skyblue)

		game.Draw(win)

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

func createGround(g *Game) Body {
	bd := box2d.MakeB2BodyDef()
	ground := g.world.CreateBody(&bd)

	shape := box2d.MakeB2EdgeShape()
	shape.Set(
		box2d.MakeB2Vec2(-1000000, 0),
		box2d.MakeB2Vec2(1000000, 0),
	)
	ground.CreateFixture(&shape, 0.0)

	g.AddBody(ground, pixel.IM, -1)

	return ground
}

func createPlayer(g *Game, pos pixel.Vec, spriteIndex int) Body {
	sprite := g.sprites[spriteIndex]

	w := sprite.Frame().W() / 2
	h := sprite.Frame().H() / 2

	bd := box2d.MakeB2BodyDef()
	bd.Position.Set(pos.X, pos.Y)
	bd.Type = box2d.B2BodyType.B2_dynamicBody
	bd.FixedRotation = true

	body := g.world.CreateBody(&bd)

	shape := box2d.MakeB2PolygonShape()
	shape.SetAsBox(w / physicsScale, h / physicsScale)

	fd := box2d.MakeB2FixtureDef()
	fd.Shape = &shape
	fd.Density = 1.0
	fd.Friction = 1.0
	body.CreateFixtureFromDef(&fd)

	scaledMat := pixel.IM.Scaled(pixel.ZV, 1/physicsScale)
	g.AddBody(body, scaledMat, spriteIndex)

	return body
}

func explode(g *Game, source Body, scale float64, spriteIndex int) {
	sourcePos := source.GetPosition()
	sourceVel := source.GetLinearVelocity()
	for i := 0.0; i < math.Pi*2; i += (math.Pi * 2) / 10 {
		x := math.Cos(i) * scale
		y := math.Sin(i) * scale

		ball := createBall(g, pixel.V(x+sourcePos.X, y+sourcePos.Y), pixel.V(1, 1))
		ball.SetLinearVelocity(sourceVel)
		ball.ApplyLinearImpulseToCenter(box2d.MakeB2Vec2(x*100, y*100), true)
	}
}

func main() {
	pixelgl.Run(run)
}
