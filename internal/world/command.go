package world

type Command any

type CmdTick struct{}
type CmdSpawnAgent struct{}
type CmdSpawnFood struct{}
type CmdSnapshot struct {
	Reply chan WorldSnapshot
}
type CmdStop struct{}
