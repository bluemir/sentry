package core

type Config struct {
	Command       string
	Exclude       string
	Delay         int32
	Path          string
	Shell         string
	KillOnRestart bool
}
