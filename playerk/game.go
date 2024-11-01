package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp/v2"
	"github.com/jakecoffman/cpebiten"
)

const (
	COLLISION_PLAYER = iota + 1
	COLLISION_SOLID
)

type Game struct {
	*cpebiten.Game
	room   *Room
	player *Player
}

func NewGame() *Game {
	space := cp.NewSpace()
	space.Iterations = 10
	space.SetGravity(cp.Vector{X: 0, Y: Gravity})

	// room
	room := NewRoom(space, screenWidth, screenHeight)

	// player
	player := NewPlayer(space)

	AddBoxes(space)

	// handler := space.NewCollisionHandler(COLLISION_PLAYER, COLLISION_SOLID)
	// handler.BeginFunc = func(arb *cp.Arbiter, space *cp.Space, data interface{}) bool {
	// 	fmt.Println("begin")
	// 	plrs, solid := arb.Shapes()
	// 	sd := solid.UserData.(*Solid)
	// 	plr := plrs.UserData.(*Player)
	// 	return plr.collistionWithSolid(arb, space, sd)
	// }
	// handler.SeparateFunc = func(arb *cp.Arbiter, space *cp.Space, data interface{}) {
	// 	fmt.Println("separate")
	// 	plrs, _ := arb.Shapes()
	// 	plr := plrs.UserData.(*Player)
	// 	plr.separateFromSolid(arb, space, data)
	// }

	return &Game{
		Game:   cpebiten.NewGame(space, 180),
		room:   room,
		player: player,
	}
}

func (g *Game) Update() error {
	jumpState := ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp)

	// If the jump key was just pressed this frame, jump!
	if jumpState && !lastJumpState && grounded {
		jumpV := math.Sqrt(2.0 * JumpHeight * Gravity)
		playerBody.SetVelocityVector(playerBody.Velocity().Add(cp.Vector{0, -jumpV}))

		remainingBoost = JumpBoostHeight / jumpV
	}

	if err := g.Game.Update(); err != nil {
		return err
	}

	remainingBoost -= 1. / 60.
	lastJumpState = jumpState

	return nil
}
