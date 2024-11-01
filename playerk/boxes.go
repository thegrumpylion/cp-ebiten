package main

import "github.com/jakecoffman/cp/v2"

func AddBoxes(space *cp.Space) {
	for i := 0; i < 6; i++ {
		for j := 0; j < 3; j++ {
			body := space.AddBody(cp.NewBody(4, cp.INFINITY))
			body.SetPosition(cp.Vector{X: float64(400 + j*60), Y: float64(200 + i*60)})

			shape := space.AddShape(cp.NewBox(body, 50, 50, 0))
			shape.SetElasticity(0)
			shape.SetFriction(0.7)
		}
	}
}
