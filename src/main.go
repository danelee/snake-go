package main

import (
	"bytes"
	"fmt"
	"time"
	"image/color"
	"log"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"math/rand"

)

const (
	//grid sizes
	SCREEN_WIDTH int = 32 
	SCREEN_HEIGHT int = 32 
	GRID_SIZE int = 16 
)


type Game struct{
	snake *Snake
	food *Food
	gameOver bool
}

func(g *Game) Update() error {
	if g.gameOver {
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			g.gameOver = false
		}
		return nil
	}
	g.snake.moveSnake()
	g.updateSnake()
	g.collisionDetection()
	if g.gameOver {
		g.resetGame()
	}
	return nil
}

func(g *Game) Draw(screen *ebiten.Image) {
	ops := &text.DrawOptions{}
	//draw game over
	if g.gameOver {
		mainText := "GAME OVER"
		w, h := text.Measure(mainText, largeFont, largeFont.Size)
		ops.GeoM.Translate(float64(SCREEN_WIDTH/2 * GRID_SIZE) - w/2,
		float64(SCREEN_HEIGHT/2 * GRID_SIZE) - h/2)
		ops.ColorScale.ScaleWithColor(color.White)
		text.Draw(screen, mainText, largeFont, ops)
		ops.GeoM.Reset()
		ops.GeoM.Translate(float64(SCREEN_WIDTH/2 * GRID_SIZE) - w/2,
		float64(SCREEN_HEIGHT/2 * GRID_SIZE) + h/2)
		text.Draw(screen, "Press enter to continue...", normalFont, ops)
		ops.GeoM.Reset()
	}else {
		//draw snake
		for _, s := range g.snake.body {
			vector.DrawFilledRect(screen, float32(s.x * GRID_SIZE), float32(s.y * GRID_SIZE),
			float32(GRID_SIZE), float32(GRID_SIZE), color.White, false)
		}

		//draw food
		vector.DrawFilledRect(screen, float32(g.food.x * GRID_SIZE), float32(g.food.y * GRID_SIZE),
			float32(GRID_SIZE), float32(GRID_SIZE), color.RGBA{255, 0, 0, 0}, false)
		
		//draw score
		ops.GeoM.Translate(0, 0)
		text.Draw(screen, fmt.Sprintf("Score: %v", g.snake.score), normalFont, ops) 
		ops.GeoM.Reset()
	}
}

func(g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREEN_WIDTH * GRID_SIZE, SCREEN_HEIGHT * GRID_SIZE
}

func main() {
	ebiten.SetWindowSize(SCREEN_WIDTH * GRID_SIZE, SCREEN_HEIGHT * GRID_SIZE)
	ebiten.SetWindowTitle("Snake, or Dragon... who knows")
	ebiten.SetTPS(10)	
	game := initGame()

	err := ebiten.RunGame(game)
	if err != nil {
		log.Fatal(err)
	}
}

type Point struct {
	x, y int
}

type BaseSprite struct {
	Point
	image *ebiten.Image
}

type Snake struct {
	body []BaseSprite
	dir Point
	score int
}

type Food struct {
	BaseSprite
}

var (
//directions
dirUp = Point{0, -1}
dirDown = Point{0, 1}
dirLeft = Point{-1, 0}
dirRight = Point{1, 0}
//fonts
mplusFaceSource *text.GoTextFaceSource
normalFont *text.GoTextFace
largeFont *text.GoTextFace
)

func initGame() *Game {
	//load font
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s	
	normalFont = &text.GoTextFace{
		Source: mplusFaceSource,
		Size: float64(GRID_SIZE),
	}
	largeFont = &text.GoTextFace{
		Source: mplusFaceSource,
		Size: float64(GRID_SIZE * 2),
	}
	
	snake := spawnSnake()
	food := spawnFood()

	return &Game {
		snake: snake,
		food: food,
		gameOver: false,
	} 
}


func (s *Snake) moveSnake() {
	if ebiten.IsKeyPressed(ebiten.KeyUp) && s.dir != dirDown {
		s.dir = dirUp
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) && s.dir != dirUp{
		s.dir = dirDown
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) && s.dir != dirRight{
		s.dir = dirLeft
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) && s.dir != dirLeft{
		s.dir = dirRight
	}
}

func (g *Game) updateSnake() {
	head := g.snake.body[0]
	
	//eat food else move
	if head.Point == g.food.Point {
		//create new head and reattach body
		head.Point = g.food.Point
		g.snake.body = append([]BaseSprite{head}, g.snake.body[:len(g.snake.body)]...)
		g.snake.score++	
		//spawn new point for food
		g.food.Point = spawnRandomPoint()
	} else{
		g.snake.body = append([]BaseSprite{head}, g.snake.body[:len(g.snake.body) - 1]...)
	}
	g.snake.body[0].x += g.snake.dir.x
	g.snake.body[0].y += g.snake.dir.y
}

func spawnSnake() *Snake {
	return &Snake{
		body: []BaseSprite{
			{
				Point: spawnRandomPoint(), 
				image: nil, 
			},
		},
		dir: Point{0, 0},
		score: 0,
	}
}

func spawnFood() *Food {
	return &Food{
		BaseSprite: BaseSprite{
				Point: spawnRandomPoint(), 
				image: nil, 
		},
	}
}

func spawnRandomPoint() Point {
	random := rand.New(rand.NewSource(int64(time.Now().UnixNano())))
	return Point{random.Intn(SCREEN_WIDTH), random.Intn(SCREEN_HEIGHT)}
}

func (g *Game) collisionDetection() {
	//collision with wall
	if g.snake.body[0].x < 0 || g.snake.body[0].x > SCREEN_WIDTH  - 1 || 
		g.snake.body[0].y < 0 || g.snake.body[0].y > SCREEN_HEIGHT - 1{
		g.gameOver = true
	}

	//collide with body
	for _, b := range g.snake.body[1:] {
		if g.snake.body[0].Point == b.Point {
			g.gameOver = true
		}
	}

}

func (g *Game) resetGame() {
	g.snake = spawnSnake() 
	g.food = spawnFood() 
}
