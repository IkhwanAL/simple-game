# Tiny Worlds
Simulation World Where Agent Search For Food, Eat, Reproduction

![Cover Tiny World Simulation](./docs/image-cover.png) 

## Why This Project Exists
I want to learn Something But a Normal CRUD is boring. So i try to create game like simulation that involve backend.

Manage learn the difference between Channel, Mutex and Atomic which both same handle concurrency be different ideology

Channel -> is Based On Communication Pipe Like Post Office
Mutex -> is Based on Queue, Where it able to change one by one
Atomic -> is more like special CPU instruction allowed them to change value in memory in atomic level at the hardware level

## How to Run

Need to have `make` command. To install `make`, look at google. (Don't Be Lazy)

And then type `make run`

## Existing Feature

- Debug UI (Pause, Stop/Start, Add New Agent Manual, Add New Food Manual, Speed Up, Slow Down)
- BFS For Path Finding
- Eat, Walk, Reproduction
- Death
- Add Trail For Path Finding (Need A Button To Be Able To Toggle)

## Current Plan

- Convert HTTP Pooling into Websocket
- Smooth Animation

## Next Future Plan

- Each Agent Able to Expand Conquer other Agent
- The Current Right Now is heavily depends on Mutex, Might Convert it Into Channel

