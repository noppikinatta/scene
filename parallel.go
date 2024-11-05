package bamenn

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

// Update is ebiten.Game implementation.
// it calls all ebiten.Game.Update in the order of index. All Updates are called and errors are joined.
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

// Draw is ebiten.Game implementation.
// it calls all ebiten.Game.Draw in the order of index.
func (p *Parallel) Draw(screen *ebiten.Image) {
	for _, g := range p.games {
		g.Draw(screen)
	}
}

// Layout is ebiten.Game implementation.
// It returns the largest width and height of all Layouts.
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

// DrawFinalScreen is ebiten.FinalScreenDrawer implementation.
// It calls only the DrawFinalScreen of the lowest index that implements FinalScreenDrawer among the Games it holds.
func (p *Parallel) DrawFinalScreen(screen ebiten.FinalScreen, offScreen *ebiten.Image, geoM ebiten.GeoM) {
	for _, g := range p.games {
		if f, ok := g.(ebiten.FinalScreenDrawer); ok {
			f.DrawFinalScreen(screen, offScreen, geoM)
			return
		}
	}

	defaultDrawFinalScreenTemporaryImplRemoveItWhenEbitengineV290Released(screen, offScreen, geoM)
}

// LayoutF is ebiten.LayoutFer implementation.
// It returns the largest width and height of all LayoutFs.
func (p *Parallel) LayoutF(outsideWidth, outsideHeight float64) (screenWidth, screenHeight float64) {
	var maxW, maxH float64 = 0, 0

	layoutF := func(g ebiten.Game) (w, h float64) {
		if l, ok := g.(ebiten.LayoutFer); ok {
			return l.LayoutF(outsideWidth, outsideHeight)
		}

		owi := int(outsideWidth)
		ohi := int(outsideHeight)

		if owi < 1 {
			owi = 1
		}
		if ohi < 1 {
			ohi = 1
		}

		wi, hi := g.Layout(owi, ohi)
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

// OnStart is OnStarter implementation.
// it calls all OnStarter.OnStart in the order of index if implemented.
func (p *Parallel) OnStart() {
	for _, g := range p.games {
		callIfImpl(g, func(o OnStarter) { o.OnStart() })
	}
}

// OnEnd is OnEnder implementation.
// it calls all OnEnder.OnEnd in the order of index if implemented.
func (p *Parallel) OnEnd() {
	for _, g := range p.games {
		callIfImpl(g, func(o OnEnder) { o.OnEnd() })
	}
}

// OnArrival is OnArrivaler implementation.
// it calls all OnArrivaler.OnArrival in the order of index if implemented.
func (p *Parallel) OnArrival() {
	for _, g := range p.games {
		callIfImpl(g, func(o OnArrivaler) { o.OnArrival() })
	}
}

// OnDeparture is OnDeparturer implementation.
// it calls all OnDeparturer.OnDeparture in the order of index if implemented.
func (p *Parallel) OnDeparture() {
	for _, g := range p.games {
		callIfImpl(g, func(o OnDeparturer) { o.OnDeparture() })
	}
}
