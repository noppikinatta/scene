package scene

import (
	"errors"

	"github.com/hajimehoshi/ebiten/v2"
)

// Parallel runs multiple Games in parallel.
type Parallel struct {
	games []ebiten.Game
	errs  []error // errs is error cache for Update()
}

// NewParallel creates a new Parallel instance.
func NewParallel(games ...ebiten.Game) *Parallel {
	return &Parallel{games: games}
}

func (p *Parallel) Update() error {
	if len(p.errs) < len(p.games) {
		p.errs = make([]error, len(p.games))
	}
	p.errs = p.errs[:len(p.games)]

	for i := range p.games {
		p.errs[i] = p.games[i].Update()
	}

	return errors.Join(p.errs...)
}

func (p *Parallel) Draw(screen *ebiten.Image) {
	for _, g := range p.games {
		g.Draw(screen)
	}
}

func (p *Parallel) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	var maxW, maxH int = 0, 0
	for _, g := range p.games {
		w, h := g.Layout(outsideWidth, outsideHeight)
		if w > maxW {
			maxW = w
		}
		if h > maxH {
			maxH = h
		}
	}

	return maxW, maxH
}

func (p *Parallel) OnStart() {
	for _, g := range p.games {
		callIfImpl(g, func(o OnStarter) { o.OnStart() })
	}
}

func (p *Parallel) OnEnd() {
	for _, g := range p.games {
		callIfImpl(g, func(o OnEnder) { o.OnEnd() })
	}
}

func (p *Parallel) OnDeparture() {
	for _, g := range p.games {
		callIfImpl(g, func(o OnDeparturer) { o.OnDeparture() })
	}
}

func (p *Parallel) OnArrival() {
	for _, g := range p.games {
		callIfImpl(g, func(o OnArrivaler) { o.OnArrival() })
	}
}

func (p *Parallel) drawFinalScreenFunc() func(screen ebiten.FinalScreen, offScreen *ebiten.Image, geoM ebiten.GeoM) {
	return func(screen ebiten.FinalScreen, offScreen *ebiten.Image, geoM ebiten.GeoM) {
		handled := false

		for _, g := range p.games {
			if f, ok := g.(ebiten.FinalScreenDrawer); ok {
				f.DrawFinalScreen(screen, offScreen, geoM)
				handled = true
			}
		}

		if !handled {
			defaultDrawFinalScreenTemporaryImplRemoveItWhenEbitengineV290Released(screen, offScreen, geoM)
		}
	}
}

func (p *Parallel) layoutFFunc() func(outsideWidth, outsideHeight float64) (screenWidth, screenHeight float64) {
	return func(outsideWidth, outsideHeight float64) (screenWidth float64, screenHeight float64) {
		var maxW, maxH float64 = 0, 0

		layoutF := func(g ebiten.Game) (w, h float64) {
			if l, ok := g.(ebiten.LayoutFer); ok {
				return l.LayoutF(outsideWidth, outsideHeight)
			}

			wi, hi := g.Layout(int(outsideWidth), int(outsideHeight))
			return float64(wi), float64(hi)
		}

		for _, g := range p.games {
			w, h := layoutF(g)
			if w > maxW {
				maxW = w
			}
			if h > maxH {
				maxH = h
			}
		}

		return maxW, maxH
	}
}
