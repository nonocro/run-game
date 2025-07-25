/*
//  Data structure for representing a game. Implements the ebiten.Game
//  interface (Update in game-update.go, Draw in game-draw.go, Layout
//  in game-layout.go). Provided with a few utilitary functions:
//    - initGame
*/

package main

import (
	"bytes"
	"course/assets"
	"image"
	"log"
	"net"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	state         int           // Current state of the game
	runnerImage   *ebiten.Image // Image with all the sprites of the runners
	runners       [4]Runner     // The four runners used in the game
	f             Field         // The running field
	launchStep    int           // Current step in StateLaunchRun state
	resultStep    int           // Current step in StateResult state
	getTPS        bool          // Help for debug
	conn          net.Conn      // your connexion
	done          bool          // Used for moving the state depending of the server
	myRunner      int           // number of your player
	nbPlayer      int           // counter of player
	counter_space []bool        // 1 boolean for each runner that is true if the player use the space bar during the race
	keys_bool     [][3]bool     // For each runner, their is a tab with 3 booleans, first is for left key, second is for right and the last is for space
}

// These constants define the five possible states of the game
const (
	StateWelcomeScreen int = iota // Title screen (iota = 0 et les autre constante sont incrémenté automatiquement de 1)
	StateChooseRunner             // Player selection screen
	StateLaunchRun                // Countdown before a run
	StateRun                      // Run
	StateResult                   // Results announcement
)

// InitGame builds a new game ready for being run by ebiten
func InitGame() (g Game) {
	// Initialisation of the new attribute
	g.counter_space = make([]bool, 4)
	g.keys_bool = make([][3]bool, 4)
	g.done = false
	g.nbPlayer = 1
	
	// Open the png image for the runners sprites
	img, _, err := image.Decode(bytes.NewReader(assets.RunnerImage))
	if err != nil {
		log.Fatal(err)
	}
	g.runnerImage = ebiten.NewImageFromImage(img)
	// Define game parameters
	start := 50.0
	finish := float64(screenWidth - 50)
	frameInterval := 20

	// Create the runners
	for i := range g.runners {
		interval := frameInterval
		// interval := 0
		// if i == 0 {
		// 	interval = frameInterval
		// }
		g.runners[i] = Runner{
			xpos: start, ypos: 50 + float64(i*20),
			maxFrameInterval: interval,
			colorScheme:      0,
		}
	}

	

	// Create the field
	g.f = Field{
		xstart:   start,
		xarrival: finish,
		chrono:   time.Now(),
	}

	return g
}
