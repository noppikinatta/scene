# Scene

Simple scene library for Ebitengine.

This library represents a game scene by a Scene interface and processes multiple Scenes in sequence using a Chain structure; the order in which the Scenes are processed can be changed dynamically.

## Other Features

### Parallel type

Multiple Scenes can be run in parallel.

### ToGame function

Convert Scene to ebiten.Game.

## How to add to your project

Add the dependency with `go get` command.

`go get github.com/noppikinatta/scene`

## How to use

See examples.

https://github.com/noppikinatta/scene/tree/main/examples
