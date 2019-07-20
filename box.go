package main

import (
	"image"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/ByteArena/box2d"
)

type boxImage struct {
	dims  pixel.Vec
	color color.Color
}

var _ image.Image = boxImage{}

func (b boxImage) At(x, y int) color.Color {
	return b.color
}

func (b boxImage) Bounds() image.Rectangle {
	return image.Rect(0, 0, int(b.dims.X), int(b.dims.Y))
}

func (b boxImage) ColorModel() color.Model {
	return color.RGBAModel
}

func createBox(game *Game, pos pixel.Vec, dims pixel.Vec) Body {
	bounds := pixel.V(1,1)
	img := boxImage{dims: bounds, color: color.RGBA{255, 0, 0, 255}}
	pic := pixel.PictureDataFromImage(img)
	sprite := pixel.NewSprite(pic, pic.Bounds())
	spriteIndex := game.AddSprite(sprite)

	w := dims.X / 2
	h := dims.Y / 2

	bd := box2d.MakeB2BodyDef()
	bd.Position.Set(pos.X, pos.Y)
	bd.Type = box2d.B2BodyType.B2_dynamicBody
	bd.FixedRotation = true

	body := game.world.CreateBody(&bd)

	shape := box2d.MakeB2PolygonShape()
	shape.SetAsBox(w, h)

	fd := box2d.MakeB2FixtureDef()
	fd.Shape = &shape
	fd.Density = 1.0
	fd.Friction = 1.0
	body.CreateFixtureFromDef(&fd)

	scaledMat := pixel.IM.ScaledXY(pixel.ZV, pixel.V(
		dims.X / bounds.X,
		dims.Y / bounds.Y,
	))
	game.AddBody(body, scaledMat, spriteIndex)

	return body
}

func createBall(game *Game, pos pixel.Vec, dims pixel.Vec) Body {
	bounds := pixel.V(1,1)
	img := boxImage{dims: bounds, color: color.RGBA{255, 0, 0, 255}}
	pic := pixel.PictureDataFromImage(img)
	sprite := pixel.NewSprite(pic, pic.Bounds())
	spriteIndex := game.AddSprite(sprite)

	w := dims.X / 2
	h := dims.Y / 2

	bd := box2d.MakeB2BodyDef()
	bd.Position.Set(pos.X, pos.Y)
	bd.Type = box2d.B2BodyType.B2_dynamicBody

	body := game.world.CreateBody(&bd)

	shape := box2d.MakeB2PolygonShape()
	shape.SetAsBox(w, h)

	fd := box2d.MakeB2FixtureDef()
	fd.Shape = &shape
	fd.Density = 1.0
	fd.Friction = 1.0
	body.CreateFixtureFromDef(&fd)

	scaledMat := pixel.IM.ScaledXY(pixel.ZV, pixel.V(
		dims.X / bounds.X,
		dims.Y / bounds.Y,
	))
	game.AddBody(body, scaledMat, spriteIndex)

	return body
}
