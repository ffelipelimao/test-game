package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	BULLET_DRUM = 5
	LIFE_POINTS = 2
)

type Game struct {
	bulletsLeft int
	lifePlayer  int
	lifeNpc     int
	message     string
	feed        []string
	playerTurn  bool
	gameOver    bool
	timer       time.Time
	waiting     bool
}

func (g *Game) Update() error {
	if g.bulletsLeft == 0 {
		g.endGame()
	}

	if g.lifeNpc == 0 || g.lifePlayer == 0 {
		g.endGame()
	}

	if g.gameOver {
		g.message = "Press R to play again"
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			g.resetGame()
		}
		return nil
	}

	if g.waiting {
		if time.Since(g.timer) >= 3*time.Second {
			g.waiting = false
			g.npcTurn()
		}
		return nil
	}

	if g.playerTurn {
		g.message = "Player: Your turn! Press Space to pull the trigger."
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.playerTurn = false
			g.checkShot(true)
		}
	} else if !g.waiting {
		g.message = "NPC: Thinking..."
		g.waiting = true
		g.timer = time.Now()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	position := fmt.Sprintf("Bullets left: %d/%d", g.bulletsLeft, BULLET_DRUM)
	ebitenutil.DebugPrintAt(screen, position, 10, 10)

	x, y := 10, 40
	for _, msg := range g.feed {
		ebitenutil.DebugPrintAt(screen, msg, x, y)
		y += 20
	}

	ebitenutil.DebugPrintAt(screen, g.message, 10, y+20)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func (g *Game) resetGame() {
	g.lifeNpc = LIFE_POINTS
	g.lifePlayer = LIFE_POINTS
	g.bulletsLeft = randomizes(LIFE_POINTS, BULLET_DRUM)
	g.playerTurn = randomizes(1, 2) == 1
	g.gameOver = false
	g.feed = []string{}
	g.message = "New game started!"
	g.addMessageOnFeed(g.message)
}

func randomizes(start, end int) int {
	return rand.Intn(end-start+1) + start
}

func (g *Game) addMessageOnFeed(message string) {
	if len(g.feed) == 0 || g.feed[len(g.feed)-1] != message {
		g.feed = append(g.feed, message)
	}
}

func (g *Game) checkShot(isPlayer bool) {
	if g.bulletsLeft > 0 {
		indexShot := randomizes(1, BULLET_DRUM)
		if indexShot <= g.bulletsLeft {
			if isPlayer {
				g.message = "BANG! You got hit..."
				g.lifePlayer--
				g.bulletsLeft--
			} else {
				g.message = "BANG! NPC got hit..."
				g.lifeNpc--
				g.bulletsLeft--
			}
		} else {
			if isPlayer {
				g.message = "Click! You survived..."
			} else {
				g.message = "Click! NPC survived..."
			}
		}
		g.addMessageOnFeed(g.message)
	}
}

func (g *Game) npcTurn() {
	g.checkShot(false)
	g.playerTurn = true
}

func (g *Game) endGame() {
	g.gameOver = true
	if g.lifeNpc == g.lifePlayer {
		g.message = "Game over... It's a draw"
		g.addMessageOnFeed(g.message)
		return
	}

	var winnerMessage string
	if g.lifeNpc > g.lifePlayer {
		winnerMessage = "NPC wins!"
	} else {
		winnerMessage = "You wins!"
	}
	g.message = fmt.Sprintf("Game over... %s", winnerMessage)
	g.addMessageOnFeed(g.message)
}

func main() {
	game := &Game{}
	game.resetGame()

	ebiten.SetWindowSize(1080, 720)
	ebiten.SetWindowTitle("Face Death!")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
