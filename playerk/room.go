package main

import (
	"github.com/jakecoffman/cp/v2"
	"github.com/jakecoffman/cpebiten"
)

type Room struct {
	space *cp.Space
}

type Solid struct {
	v1 cp.Vector
	v2 cp.Vector
}

func (s *Solid) IsHorizontal() bool {
	return s.v1.Y == s.v2.Y
}

func (s *Solid) IsVertical() bool {
	return s.v1.X == s.v2.X
}

func NewRoom(space *cp.Space, width, height int) *Room {
	screenWidth, screenHeight := float64(width), float64(height)
	walls := []cp.Vector{
		// left wall
		{X: 0, Y: 0},
		{X: 0, Y: screenHeight},
		// right wall
		{X: screenWidth, Y: 0},
		{X: screenWidth, Y: screenHeight},
		// ceiling
		{X: 0, Y: 0},
		{X: screenWidth, Y: 0},
		// floor
		{X: 0, Y: screenHeight},
		{X: screenWidth, Y: screenHeight},
	}
	for i := 0; i < len(walls)-1; i += 2 {
		shape := space.AddShape(cp.NewSegment(space.StaticBody, walls[i], walls[i+1], 0))
		shape.SetElasticity(1)
		shape.SetFriction(1)
		shape.SetFilter(cpebiten.NotGrabbable)
		shape.SetCollisionType(COLLISION_SOLID)
		shape.UserData = &Solid{v1: walls[i], v2: walls[i+1]}
	}
	return &Room{
		space: space,
	}
}
