package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp/v2"
)

const (
	screenWidth  = 600
	screenHeight = 480
)

const (
	PlayerVelocity = 500.0

	PlayerGroundAccelTime = 0.1
	PlayerGroundAccel     = PlayerVelocity / PlayerGroundAccelTime

	PlayerAirAccelTime = 0.25
	PlayerAirAccel     = PlayerVelocity / PlayerAirAccelTime

	JumpHeight      = 50.0
	JumpBoostHeight = 95.0
	FallVelocity    = 900.0
	Gravity         = 2000.0
)

var (
	playerBody  *cp.Body
	playerShape *cp.Shape
)

var (
	remainingBoost          float64
	grounded, lastJumpState bool
)

type Player struct {
	space *cp.Space
	shape *cp.Shape
	joint *cp.Constraint
}

func NewPlayer(space *cp.Space) *Player {
	player := &Player{
		space: space,
	}
	playerBody = space.AddBody(cp.NewBody(1, cp.INFINITY))
	playerBody.SetPosition(cp.Vector{X: 100, Y: 200})
	playerBody.SetVelocityUpdateFunc(player.playerUpdateVelocity)
	playerShape = space.AddShape(cp.NewBox2(playerBody, cp.BB{L: -15, B: -27.5, R: 15, T: 27.5}, 10))
	playerShape.SetElasticity(0)
	playerShape.SetFriction(0)
	playerShape.SetCollisionType(COLLISION_PLAYER)
	playerShape.UserData = player
	player.shape = playerShape
	return player
}

func (p *Player) collistionWithSolid(arb *cp.Arbiter, space *cp.Space, data interface{}) bool {
	plr, solid := arb.Shapes()
	sd := data.(*Solid)
	if sd.IsVertical() && p.joint == nil {
		contacts := arb.ContactPointSet()
		fmt.Println(contacts, plr.Body().Position(), solid.Body().Position())
		joint := cp.NewGrooveJoint(plr.Body(), solid.Body(), sd.v2, sd.v1, cp.Vector{X: 0, Y: 100})
		space.AddPostStepCallback(p.startWallSliding, joint, nil)
		p.joint = joint
		remainingBoost = 0
	}
	return true
}

func (p *Player) separateFromSolid(arb *cp.Arbiter, space *cp.Space, data interface{}) {
	if p.joint != nil {
		space.AddPostStepCallback(p.stopWallSliding, p.joint, nil)
	}
}

func (p *Player) startWallSliding(space *cp.Space, key, _ interface{}) {
	space.AddConstraint(key.(*cp.Constraint))
}

func (p *Player) stopWallSliding(space *cp.Space, key, _ interface{}) {
	space.RemoveConstraint(key.(*cp.Constraint))
	p.joint = nil
}

func (p *Player) playerUpdateVelocity(body *cp.Body, gravity cp.Vector, damping, dt float64) {
	jumpState := ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp)

	// Grab the grounding normal from last frame
	groundNormal := cp.Vector{}
	counter := 0
	playerBody.EachArbiter(func(arb *cp.Arbiter) {
		counter++
		sa, sb := arb.Shapes()
		fmt.Println("sa", sa, "sb", sb)
		n := arb.Normal().Neg()

		if n.Y < groundNormal.Y {
			groundNormal = n
		}
	})
	fmt.Println("counter", counter)

	grounded = groundNormal.Y < 0
	if groundNormal.Y > 0 {
		remainingBoost = 0
	}

	// Do a normal-ish update
	boost := jumpState && remainingBoost > 0
	var g cp.Vector
	if !boost {
		g = gravity
	}
	if p.joint != nil {
		g = gravity
	}
	body.UpdateVelocity(g, damping, dt)

	// Target horizontal speed for air/ground control
	var targetVx float64
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		targetVx -= PlayerVelocity
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		targetVx += PlayerVelocity
	}

	// Update the surface velocity and friction
	// Note that the "feet" move in the opposite direction of the player.
	surfaceV := cp.Vector{X: -targetVx, Y: 0}
	playerShape.SetSurfaceV(surfaceV)
	if grounded {
		playerShape.SetFriction(PlayerGroundAccel / Gravity)
	} else {
		playerShape.SetFriction(0)
	}

	// Apply air control if not grounded
	if !grounded {
		v := playerBody.Velocity()
		playerBody.SetVelocity(cp.LerpConst(v.X, targetVx, PlayerAirAccel*dt), v.Y)
	}

	v := body.Velocity()
	body.SetVelocity(v.X, cp.Clamp(v.Y, -FallVelocity, cp.INFINITY))
}
