package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	tb "github.com/nsf/termbox-go"
)

var which, guess = -1, -1
var results [2][9]struct{ ok, bad int }
var hist [17]int

func pwhite() {
	tb.Clear(tb.ColorWhite, tb.ColorWhite)
	flush()
}

func pstring(s string, x, y int, fg, bg tb.Attribute) int {
	for _, c := range s {
		tb.SetCell(x, y, c, fg, bg)
		x++
	}
	return x
}

func presults() {
	tb.Clear(tb.ColorBlack, tb.ColorBlack)
	width, height := tb.Size()
	for k, results := range results {

		for i, r := range results {
			y := height - i - 1
			tb.SetCell(0, y, rune('1'+i), tb.ColorWhite, tb.ColorBlack)
			x := pstring(fmt.Sprintf("%d", r.ok), 2+12*k, y, tb.ColorGreen, tb.ColorBlack)
			x = pstring(fmt.Sprintf("%d", r.bad), x+1, y, tb.ColorRed, tb.ColorBlack)
			if guess != -1 {
				if i == which {
					if which == guess {
						pstring("***", x+2, y, tb.ColorGreen, tb.ColorBlack)
					} else {
						pstring("*", x+2, y, tb.ColorGreen, tb.ColorBlack)
					}
				} else if i == guess {
					pstring("***", x+2, y, tb.ColorRed, tb.ColorBlack)
				}
			}
		}
	}
	if guess != -1 {
		if guess != which {
			var r rune = 11014 // up arrow
			if guess < which {
				r++
			}
			tb.SetCell(width/2, height/2, r, tb.ColorRed, tb.ColorBlack)
		}
		tb.SetCell(width/2, height/2+1, rune('1'+which), tb.ColorGreen, tb.ColorBlack)
	}
	flush()
}

func flush() {
	if err := tb.Flush(); err != nil {
		tb.Close()
		log.Fatal(err)
	}
}

func quitev(ev tb.Event) bool {
	return ev.Type == tb.EventKey && (ev.Key == tb.KeyEsc || ev.Key == tb.KeyCtrlD || ev.Ch == 'q')
}

func main() {
	if err := tb.Init(); err != nil {
		log.Fatal(err)
	}
	defer tb.Close()
	rand.Seed(int64(time.Now().Nanosecond()))
	for {
		for true {
			presults()
			ev := tb.PollEvent()
			if quitev(ev) {
				return
			}
			if ev.Type == tb.EventKey {
				break
			}
		}
		which, guess = rand.Intn(9), -1
	flash:
		pwhite()
		time.Sleep(time.Duration(which+1) * 100 * time.Millisecond)
		for {
			presults()
			ev := tb.PollEvent()
			if quitev(ev) {
				return
			}
			if ev.Key == tb.KeySpace {
				goto flash
			}
			if ev.Ch >= '1' && ev.Ch <= '9' {
				guess = int(ev.Ch - '1')
				break
			}
			presults()
		}
		if guess == which {
			results[0][which].ok++
			results[1][which].ok++
		} else {
			results[0][which].bad++
			if guess == which+1 || guess == which-1 {
				results[1][which].ok++
			} else {
				results[1][which].bad++
			}
		}
	}
}
