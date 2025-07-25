/*
//  New file for our functions to manage other runners
*/

package main

// Will update the runner in function of a boolean, it will be call with boolean receive by the server
func (r *Runner) ServerUpdate(b bool) {
	r.UpdateSpeed(b)
	r.UpdatePos()
}

// Manage the selection to the player en function of boolean, it will be call with booleans receive by the server
func (r *Runner) ServerChoose(right bool, left bool, space bool) (done bool) {
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
