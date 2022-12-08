/*
//  Implementation of the Update method for the Game structure
//  This method is called once at every frame (60 frames per second)
//  by ebiten, juste before calling the Draw method (game-draw.go).
//  Provided with a few utilitary methods:
//    - CheckArrival
//    - ChooseRunners
//    - HandleLaunchRun
//    - HandleResults
//    - HandleWelcomeScreen
//    - Reset
//    - UpdateAnimation
//    - UpdateRunners
*/

package main

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	//"net"
	//"log"
	//"sync"
	"strconv"
)

// HandleWelcomeScreen waits for the player to push SPACE in order to
// start the game
func (g *Game) HandleWelcomeScreen() bool {
	return g.done
}

// ChooseRunners loops over all the runners to check which sprite each
// of them selected
func (g *Game) ChooseRunners() (done bool) {
	done = true
	for i := range g.runners {
		if i == g.myRunner {
			done = g.runners[i].ManualChoose() && done
			myRunner := strconv.Itoa(g.myRunner)
			var selection_failed = false
			if g.runners[i].colorSelected {
				for _, runner := range g.runners {
					if runner.colorSelected && runner.colorScheme == g.runners[i].colorScheme && g.runners[i] != runner {
						g.runners[i].colorSelected = false
						selection_failed = true
						done = false
					}
				}
			}
			if !selection_failed {
				fmt.Fprintf(g.conn, ":key"+","+myRunner+","+strconv.FormatBool(inpututil.IsKeyJustPressed(ebiten.KeyRight))+","+strconv.FormatBool(inpututil.IsKeyJustPressed(ebiten.KeyLeft))+","+strconv.FormatBool(inpututil.IsKeyJustPressed(ebiten.KeySpace))+","+"\n")
			}
		} else {
			// done = g.runners[i].RandomChoose() && done
			done = g.runners[i].ServerChoose(g.keys_bool[i][0], g.keys_bool[i][1], g.keys_bool[i][2]) && done
			g.keys_bool[i] = [3]bool{false, false, false}
		}
	}
	return done
}

// HandleLaunchRun countdowns to the start of a run
func (g *Game) HandleLaunchRun() bool {
	if time.Since(g.f.chrono).Milliseconds() > 1000 {
		g.launchStep++
		g.f.chrono = time.Now()
	}
	if g.launchStep >= 5 {
		g.launchStep = 0
		return true
	}
	return false
}

// UpdateRunners loops over all the runners to update each of them
func (g *Game) UpdateRunners() {
	for i := range g.runners {
		if i == g.myRunner {
			if g.runners[i].ManualUpdate() {
				fmt.Fprintf(g.conn, ":space"+strconv.Itoa(g.myRunner)+"\n")
			}
		} else {
			g.runners[i].ServerUpdate(g.counter_space[i])
			g.counter_space[i] = false
		}
	}
}

// CheckArrival loops over all the runners to check which ones are arrived
func (g *Game) CheckArrival() (finished bool) {
	finished = true
	for i := range g.runners {
		if i == g.myRunner && g.runners[i].arrived {
			fmt.Fprintf(g.conn, strconv.Itoa(g.myRunner)+":r"+strconv.Itoa(int(g.runners[i].runTime))+"\n")
		}
		g.runners[i].CheckArrival(&g.f)
		finished = finished && g.runners[i].arrived
	}
	return finished
}

// Reset resets all the runners and the field in order to start a new run
func (g *Game) Reset() {
	for i := range g.runners {
		g.runners[i].Reset(&g.f)
	}
	g.f.Reset()
}

// UpdateAnimation loops over all the runners to update their sprite
func (g *Game) UpdateAnimation() {
	for i := range g.runners {
		g.runners[i].UpdateAnimation(g.runnerImage)
	}
}

// HandleResults computes the resuls of a run and prepare them for
// being displayed
func (g *Game) HandleResults() bool {
	if time.Since(g.f.chrono).Milliseconds() > 1000 || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.resultStep++
		g.f.chrono = time.Now()
	}
	if g.resultStep >= 4 && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.resultStep = 0
		return true
	}
	return false
}

// Update is the main update function of the game. It is called by ebiten
// at each frame (60 times per second) just before calling Draw (game-draw.go)
// Depending of the current state of the game it calls the above utilitary
// function and then it may update the state of the game
func (g *Game) Update() error {
	switch g.state {
	case StateWelcomeScreen:
		done := g.HandleWelcomeScreen()
		if done {
			g.state++
			g.done = false
			g.nbPlayer = 0
		}
	case StateChooseRunner:
		done := g.ChooseRunners()
		if done {
			fmt.Fprintf(g.conn, "Player "+strconv.Itoa(g.myRunner)+" choose his skin"+"\n")
		}
		if done {
			fmt.Fprintf(g.conn, ":skins"+"\n")
			g.done = false
			g.UpdateAnimation()
			g.state++
		}
	case StateLaunchRun:
		done := g.HandleLaunchRun()
		if done {
			g.state++
			g.done = false
		}
	case StateRun:
		g.UpdateRunners()
		finished := g.CheckArrival()
		g.UpdateAnimation()
		if finished && g.done {
			g.state++
			g.done = false
		}
	case StateResult:
		done := g.HandleResults()
		if done {
			fmt.Fprintf(g.conn, "Player "+strconv.Itoa(g.myRunner)+" want to restart"+"\n")
		}
		if g.nbPlayer == 4 {
			g.Reset()
			g.state = StateLaunchRun
			g.done = false
			g.nbPlayer = 0
			g.resultStep = 0
		}
	}
	return nil
}
