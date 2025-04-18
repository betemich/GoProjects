package main

import "fmt"

type Ranked interface {
	CalcCost()
}

type Gold struct {
	Nickname    string
	Cost        int32
	PlayedGames int64
}

func (g *Gold) CalcCost() {
	g.PlayedGames++
	fmt.Printf("Gold player: %s, cost of game: %d, played games: %d\n", g.Nickname, g.Cost, g.PlayedGames)
}

type Platinum struct {
	Nickmame    string
	Cost        int32
	PlayedGames int64
}

func (p *Platinum) CalcCost() {
	p.PlayedGames++
	fmt.Printf("Platinum player: %s, cost of game: %d, played games: %d\n", p.Nickmame, p.Cost, p.PlayedGames)
}

func PrintInfo(r Ranked, Counter func()) {
	r.CalcCost()
	Counter()
}

func HoursCounter() func() {
	var Hours int32 = 0
	return func() {
		Hours += 2
		fmt.Printf("Hours: %d\n", Hours)
	}
}

func main() {
	p := Platinum{"average 483 enjoyer", 40, 0}
	g := Gold{"abobus228", 30, 0}
	HrsCounter := HoursCounter()
	PrintInfo(&g, HrsCounter)
	PrintInfo(&p, HrsCounter)
}
