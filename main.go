package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	tb "github.com/nsf/termbox-go"
)

var which, guess = -1, -1
var results [5][9]struct{ ok, bad int }

func whiteScreen() {
	tb.Clear(tb.ColorWhite, tb.ColorWhite)
	flush()
}

func drawString(s string, x, y int, fg, bg tb.Attribute) int {
	for _, c := range s {
		tb.SetCell(x, y, c, fg, bg)
		x++
	}
	return x
}

func resultsScreen() {
	tb.Clear(tb.ColorBlack, tb.ColorBlack)
	width, height := tb.Size()
	for k, results := range results {
		totOK, totBad := 0, 0
		for i, r := range results {
			totOK += r.ok
			totBad += r.bad
			if height >= 12 {
				y := height - i - 2
				tb.SetCell(0, y, rune('1'+i), tb.ColorWhite, tb.ColorBlack)
				x := drawString(fmt.Sprintf("%2d", r.ok), 2+12*k, y, tb.ColorGreen, tb.ColorBlack)
				x = drawString(fmt.Sprintf("%2d", r.bad), x+1, y, tb.ColorRed, tb.ColorBlack)
				if guess != -1 {
					if i == which {
						if which == guess {
							drawString("***", x+2, y, tb.ColorGreen, tb.ColorBlack)
						} else {
							drawString("*", x+2, y, tb.ColorGreen, tb.ColorBlack)
						}
					} else if i == guess {
						drawString("***", x+2, y, tb.ColorRed, tb.ColorBlack)
					}
				}
			}
		}
		if totOK+totBad > 0 {
			drawString(fmt.Sprintf("%04d", 1000*totOK/(totOK+totBad)), 2+12*k+1, height-1, tb.ColorWhite, tb.ColorBlack)
		}
	}
	if guess != -1 {
		y := (height - 12) / 2
		if y < 0 {
			y = (height - 3) / 2
		}
		if guess != which {
			var r rune = 11014 // up arrow
			if guess < which {
				r++
			}
			tb.SetCell(width/2, y, r, tb.ColorRed, tb.ColorBlack)
		}
		tb.SetCell(width/2, y+1, rune('1'+which), tb.ColorGreen, tb.ColorBlack)
	}
	flush()
}

func flush() {
	if err := tb.Flush(); err != nil {
		tb.Close()
		log.Fatal(err)
	}
}

func isQuit(ev tb.Event) bool {
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
			resultsScreen()
			ev := tb.PollEvent()
			if isQuit(ev) {
				return
			}
			if ev.Type == tb.EventKey {
				break
			}
		}
		which, guess = rand.Intn(9), -1
	flash:
		whiteScreen()
		time.Sleep(time.Duration(which+1) * 100 * time.Millisecond)
		for {
			resultsScreen()
			ev := tb.PollEvent()
			if isQuit(ev) {
				return
			}
			if ev.Key == tb.KeySpace {
				goto flash
			}
			if ev.Ch >= '1' && ev.Ch <= '9' {
				guess = int(ev.Ch - '1')
				break
			}
			resultsScreen()
		}
		d := which - guess
		if d < 0 {
			d = -d
		}
		for k := range results {
			if d <= k {
				results[k][which].ok++
			} else {
				results[k][which].bad++
			}
		}
	}
}
