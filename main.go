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
	indexesOfBullets int
	indexShot        int
	lifePlayer       int
	lifeNpc          int
	message          string
	feed             []string
	playerTurn       bool
	gameOver         bool
	timer            time.Time
	waiting          bool
}

func (g *Game) Update() error {
	if g.lifeNpc == 0 || g.lifePlayer == 0 {
		g.gameOver = true
		winnerMessage := "You!"
		if g.lifeNpc > g.lifePlayer {
			winnerMessage = "NPC!"
		}
		g.message = fmt.Sprintf("Game over... %s wins", winnerMessage)
		g.addMessageOnFeed(g.message)
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
	position := fmt.Sprintf("Bullets in gun: %d/%d", g.indexesOfBullets, BULLET_DRUM)
	ebitenutil.DebugPrintAt(screen, position, 10, 10)

	x, y := 10, 40
	for i := 0; i < len(g.feed); i++ {
		ebitenutil.DebugPrintAt(screen, g.feed[i], x, y)
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
	g.indexesOfBullets = randomizes(1, BULLET_DRUM)
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

/*
TODO: Preciso fazer a lógica de contagem de balas
Ex: Se tem 3 balas e o NPC dispara uma, sobram 2 de 5 para o Player
*/
func (g *Game) checkShot(isPlayer bool) {
	g.indexShot = randomizes(1, BULLET_DRUM)
	if g.indexShot > g.indexesOfBullets {
		if isPlayer {
			g.message = "BANG! You got hit..."
			g.lifePlayer--
		} else {
			g.message = "BANG! NPC got hit..."
			g.lifeNpc--
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

func (g *Game) npcTurn() {
	g.checkShot(false)
	g.playerTurn = true
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
