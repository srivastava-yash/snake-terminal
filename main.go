package main

import (
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/mattn/go-tty"
)

type direction int

type position [2]int

type snake struct {
    direction direction
    body []position
}

type game struct {
    score int
    snake *snake
    food position
}

const (
	north direction = iota
	east
	south
	west
)

func main() {
    game := new_game()
    game.init_game()

    for {
        width, height := get_terminal_window_size()

        head_position := game.snake.body[0]

        switch game.snake.direction {
        case north: head_position[1]--
        case east: head_position[0]++
        case west: head_position[0]--
        case south: head_position[1]++
        }

        if head_position[0] < 1 || head_position[0] > width ||
            head_position[1] < 1 || head_position[1] > height {
            game.end_game()
        }

        for _, pos := range game.snake.body {
			if is_posiion_overlap(head_position, pos) {
				game.end_game()
			}
		}

        game.snake.body = append([]position{head_position}, game.snake.body...)

        if is_posiion_overlap(head_position, game.food) {
            game.score++
            game.update_food()
        } else {
            game.snake.body = game.snake.body[:len(game.snake.body)-1]
        }

        game.draw()
    }
}

func get_random_position() position {
    width, height := get_terminal_window_size()

    x := rand.Intn(width) + 1
	y := rand.Intn(height) + 2 // because we will display score at y=1

	return [2]int{x, y}
}

func get_random_direction() direction {
    directions := []direction{north, east, south, west}
    randomIndex := rand.Intn(len(directions))

    return directions[randomIndex]
}

func get_new_snake() *snake {
    width, height := get_terminal_window_size()
    start_position := position{ width/2, height/2 }

    return &snake {
        body: []position{start_position},
        direction: get_random_direction(),
    }
}

func new_game() *game {
	snake := get_new_snake()

	game := &game{
		score: 0,
		snake: snake,
		food:  get_random_position(),
	}

	return game
}

func (g *game) init_game() {
    hideCursor()

	go g.listen_for_key_press()

    // handle CTRL C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			g.end_game()
		}
	}()
}

func (g *game) end_game() {
    clear()
	showCursor()

	moveCursor(position{1, 1})
	print_string("game over. score: " + strconv.Itoa(g.score) + "\n")

	display()

	os.Exit(0)
}

func is_posiion_overlap(pos1, pos2 position) bool {
    return pos1[0] == pos2[0] && pos1[1] == pos2[1]
}

func (g *game) update_food() {
    for {
        new_pos := get_random_position()

        if(is_posiion_overlap(new_pos, g.food)) {
            continue
        }

        for _, pos := range g.snake.body {
            if is_posiion_overlap(new_pos, pos) {
                continue
            }
        }

        g.food = new_pos

        break
    }
}

func (g *game) draw() {
    clear()
	width, _ := get_terminal_window_size()

	status := "score: " + strconv.Itoa(g.score)
	statusXPos := width/2 - len(status)/2

	moveCursor(position{statusXPos, 0})
	print_string(status)

	moveCursor(g.food)
	print_string("*")

	for i, pos := range g.snake.body {
		moveCursor(pos)

		if i == 0 {
			print_string("O")
		} else {
			print_string("o")
		}
	}

	display()
	time.Sleep(time.Millisecond * 75)
}

func (g *game) listen_for_key_press() {
	tty, err := tty.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer tty.Close()

	for {
		char, err := tty.ReadRune()
		if err != nil {
			panic(err)
		}

		// UP, DOWN, RIGHT, LEFT == [A, [B, [C, [D
		// we ignore the escape character [
		switch char {
		case 'A':
			g.snake.direction = north
		case 'B':
			g.snake.direction = south
		case 'C':
			g.snake.direction = east
		case 'D':
			g.snake.direction = west
        case 'q':
            g.end_game()
		}
	}
}
