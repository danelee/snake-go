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
}

func(g *Game) Update() error {
	g.snake.moveSnake()
	g.updateSnake()
	return nil
}

func(g *Game) Draw(screen *ebiten.Image) {
	
	//draw snake
	for _, s := range g.snake.body {
		vector.DrawFilledRect(screen, float32(s.x * GRID_SIZE), float32(s.y * GRID_SIZE),
		float32(GRID_SIZE), float32(GRID_SIZE), color.White, false)
	}

	//draw food
	vector.DrawFilledRect(screen, float32(g.food.x * GRID_SIZE), float32(g.food.y * GRID_SIZE),
		float32(GRID_SIZE), float32(GRID_SIZE), color.RGBA{255, 0, 0, 0}, false)
	
	//draw score
	ops := &text.DrawOptions{}
	ops.GeoM.Translate(0, 0)
	ops.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, fmt.Sprintf("Score: %v", g.snake.score), fontFace, ops) 
	
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
fontFace *text.GoTextFace
)

func initGame() *Game {
	//load font
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s	
	fontFace = &text.GoTextFace{
		Source: mplusFaceSource,
		Size: float64(GRID_SIZE),
	};
	//init snake
	snake := &Snake{
		body: []BaseSprite{
			{
				Point: spawnRandomPoint(), 
				image: nil, 
			},
		},
		dir: Point{0, 0},
		score: 0,
	}
	//init food
	food := &Food{
		BaseSprite: BaseSprite{
				Point: spawnRandomPoint(), 
				image: nil, 
		},

	}
	
	return &Game {
		snake: snake,
		food: food,
	} 
}


func (s *Snake) moveSnake() {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		s.dir = dirUp
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
		s.dir = dirDown
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		s.dir = dirLeft
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
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


func spawnRandomPoint() Point {
	random := rand.New(rand.NewSource(int64(time.Now().UnixNano())))
	return Point{random.Intn(SCREEN_WIDTH), random.Intn(SCREEN_HEIGHT)}
}
