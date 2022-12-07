/*
//  Data structure for representing the four runners
//  used in the game. Provided with a few utilitary
//  methods:
//    - CheckArrival
//    - Draw
//    - DrawSelection
//    - ManualChoose
//    - ManualUpdate
//    - RandomChoose
//    - RandomUpdate
//    - Reset
//    - UpdateAnimation
//    - UpdatePos
//    - UpdateSpeed
*/

package main

import "fmt"

// ManualUpdate allows to use the keyboard in order to control a runner
// when the game is in the StateRun state (i.e. during a run)
func (r *Runner) ServerUpdate(b bool) {
	//r.UpdateSpeed(inpututil.IsKeyJustPressed(ebiten.KeySpace))
	r.UpdateSpeed(b)
	r.UpdatePos()
}

// ManualChoose allows to use the keyboard for selecting the appearance of a
// runner when the game is in StateChooseRunner state (i.e. at player selection
// screen)
func (r *Runner) ServerChoose(left bool, right bool, space bool) (done bool) {
	if left {
		fmt.Println("left", left)
	}
	if right {
		fmt.Println("right", right)
	}
	if space {
		fmt.Println("space", space)
	}
	r.colorSelected =
		(!r.colorSelected && space) ||
			(r.colorSelected && !space)
	if !r.colorSelected {
		if right {
			r.colorScheme = (r.colorScheme + 1) % 8
		} else if left {
			r.colorScheme = (r.colorScheme + 7) % 8
		}
	}
	return r.colorSelected
}
