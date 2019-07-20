package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/ByteArena/box2d"
	"github.com/faiface/pixel"
)

const physicsScale = 8.0

type Game struct {
	world         *box2d.B2World
	physicsBodies []*box2d.B2Body
	renderSprites []int
	renderMatZV   []pixel.Matrix

	renderMatMutex       *sync.Mutex
	renderMatBackBuffer  []pixel.Matrix
	renderMatFrontBuffer []pixel.Matrix

	sprites []*pixel.Sprite
	batches []*pixel.Batch
}

func NewGame() *Game {
	gravity := box2d.MakeB2Vec2(0.0, -9.8)
	world := box2d.MakeB2World(gravity)

	return &Game{
		world:          &world,
		renderMatMutex: new(sync.Mutex),
	}
}

func (g *Game) AddSprite(sprite *pixel.Sprite) int {
	g.sprites = append(g.sprites, sprite)
	// TODO: auto-sprite-map and batch
	g.batches = append(g.batches, pixel.NewBatch(&pixel.TrianglesData{}, sprite.Picture()))
	return len(g.sprites) - 1
}

func (g *Game) AddBody(body *box2d.B2Body, baseMat pixel.Matrix, sprite int) {
	if sprite >= len(g.sprites) {
		panic(fmt.Sprintf("Invalid sprite index: %d", sprite))
	}
	g.physicsBodies = append(g.physicsBodies, body)
	g.renderSprites = append(g.renderSprites, sprite)
	g.renderMatZV = append(g.renderMatZV, baseMat)
}

func (g *Game) Draw(target pixel.Target) {
	for _, batch := range g.batches {
		batch.Clear()
	}

	g.renderMatMutex.Lock()
	for i, mat := range g.renderMatFrontBuffer {
		bodyIndex := g.renderSprites[i]
		if bodyIndex < 0 {
			continue
		}
		batch := g.batches[bodyIndex]
		sprite := g.sprites[bodyIndex]

		sprite.Draw(batch, mat)
	}
	g.renderMatMutex.Unlock()

	for _, batch := range g.batches {
		batch.Draw(target)
	}
}

func (g *Game) StartPhysics() func(func()) {
	var workChan = make(chan func())

	go func() {
		for {
			select {
			case <-time.Tick(time.Second / 60):
				g.PhysicsStep()
			case f := <-workChan:
				f()
			}
		}
	}()

	return func(f func()) {
		workChan <- f
	}
}

func (g *Game) PhysicsStep() {
	g.world.Step(1/60.0, 8, 3)
	if len(g.renderMatBackBuffer) <= len(g.physicsBodies) {
		g.renderMatBackBuffer = make([]pixel.Matrix, len(g.physicsBodies))
	}
	for i, body := range g.physicsBodies {
		var angle = body.GetAngle()
		var position = body.GetPosition()

		var mat = g.renderMatZV[i]
		mat = mat.Rotated(pixel.ZV, angle)
		mat = mat.Moved(pixel.V(position.X, position.Y))

		g.renderMatBackBuffer[i] = mat
	}

	g.renderMatMutex.Lock()
	g.renderMatFrontBuffer, g.renderMatBackBuffer = g.renderMatBackBuffer, g.renderMatFrontBuffer
	g.renderMatMutex.Unlock()
}
