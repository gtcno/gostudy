package main

import (
	"github.com/fatih/color"
	"github.com/marcsauter/cli"
	"os"
	"strconv"
	"bytes"
	"fmt"
	"math/rand"
)

var (
	JackPot int
	SlotColour  map[int]func(string, ...interface{}) string
)

const (
	Prompt = ">>> "
	Histsize = 255
)

func init() {
	JackPot = 0
	SlotColour = map[int]func(string, ...interface{}) string{
		0: color.BlackString,
		1: color.WhiteString,
		2: color.GreenString,
		3: color.YellowString,
	}

}

func draw() []int {
	var s []int
	for i := 0; i < 4; i++ {
		s = append(s, rand.Intn(4))
	}
	return s
}

func display(slots []int, winner bool) {
	var buffer bytes.Buffer
	for _, v := range slots {
		buffer.WriteString(SlotColour[v]("="))
	}
	if (winner) {
		buffer.WriteString(" - You win!")
	} else {
		buffer.WriteString(" - You lose!")
	}
	fmt.Println(buffer.String())
}

func isWinner(slots []int) bool {
	return allSameInts(slots)
}

func allSameInts(ints []int) bool {
	for i := 1; i < len(ints); i++ {
		if ints[i] != ints[0] {
			return false
		}
	}
	return true
}

func main() {
	c := cli.NewCLI(nil, Prompt, Histsize)
	coins, _ := strconv.Atoi(c.InputString("Insert coins", nil))
	for coins > 0 {
		c.Info("Credits: %d  JackPot: %d \n", coins, JackPot)
		choice := c.Choice(map[string]string{
			"p": "play",
			"e": "exit",
		})
		if (choice == "e") {
			os.Exit(0)
		} else {
			JackPot++
			coins--
			result := draw()
			winner := isWinner(result)
			display(result, winner)
			if winner {
				coins = JackPot + coins
				JackPot = 0
			}
		}
	}
}
