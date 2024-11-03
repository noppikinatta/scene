# Scene

Simple scene library for Ebitengine.

This library provides the ability to run multiple `ebiten.Game` s in sequence, with each `ebiten.Game` as a single scene.

## Key features: Sequence and Transition types

`Sequence` structure runs multiple `ebiten.Game` s in sequence. You can switch `ebiten.Game` s with `Sequence.Switch` or `Sequence.SwitchWithTransition`. `Transition` is used to draw (e.g. fade in and out) when switching between `ebiten.Game`.

`Sequence` also implements `ebiten.Game`.

## Event functions

If `ebiten.Game` implements some or all of the `OnStarter`, `OnArrivaler`, `OnDeparturer` and `OnEnder` interfaces, they are called at the following times:

### OnStarter

`OnStarter.OnStart` is called immediately after `ebiten.Game` is switched.

 It is called just before the first `ebiten.Game.Update` after the game starts.

### OnArrivaler

`OnArrivaler.OnArrival` is called when the switched `ebiten.Game` starts and the `Transition` process is complete. For example, when the fade-in is complete. This function is useful, for example, to enable player input upon completion of a scene transition.

 It is called just before the first `ebiten.Game.Update` after the game starts and after `OnStarter.OnStart`.

### OnDeparturer

`OnDeparturer.OnDeparture` is called when a `Transition` is started to switch scenes. This is useful to disable player input at the beginning of a scene transition.

It is not called when the game is terminated by `ebiten.Termination`.

### OnEnder

`OnEnder.OnEnd` is called just before `ebiten.Game` switches.

It is not called when the game is terminated by `ebiten.Termination`.

## Parallel type

`Parallel` structure handles multiple `ebiten.Game`s in parallel. The order of processing is constant.

`Parallel.Layout` calls all `ebiten.Game.Layout` and returns the largest return value.

## Limitations

- `OnStarter`, `OnArrivaler`, `OnDeparturer` and `OnEnder` are called by `Sequence`. Use `Sequence` to enable them.
- When switching scenes, `ebiten.Game.Layout` of the previous scene may be called. If the sizes returned by the previous and following `ebiten.Game.Layout` are different, the size of the screen passed to `ebiten.Game.Draw` after the changeover may be unintended.

## How to add to your project

Add the dependency with `go get` command.

`go get github.com/noppikinatta/scene`

## How to use

See examples.

https://github.com/noppikinatta/scene/tree/main/examples
