package userinterface

import "github.com/veandco/go-sdl2/sdl"

func DrawCircle(renderer *sdl.Renderer, centerX int32, centerY int32, r int32) error {
	d := r * 2

	x := r - 1
	y := int32(0)
	tx := int32(1)
	ty := int32(1)
	e := tx - d
	var err error

	type op func(int32, int32) int32
	add := func(x, y int32) int32 { return x + y }
	sub := func(x, y int32) int32 { return x - y }
	ops := []op{add, sub}
	for x >= y {
		for _, op1 := range ops {
			for _, op2 := range ops {
				err = renderer.DrawPoint(op1(centerX, x), op2(centerY, y))
				if err != nil {
					return err
				}
				err = renderer.DrawPoint(op1(centerX, y), op2(centerY, x))
				if err != nil {
					return err
				}
			}
		}

		if e <= 0 {
			y++
			e += ty
			ty += 2
		}
		if e <= 0 {
			x--
			tx += 2
			e += tx - d
		}
	}
	return err
}
