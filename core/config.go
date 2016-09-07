package core

type Config struct {
	Command       string
	Exclude       string
	Delay         int32
	WatchPaths    []string
	Shell         string
	KillOnRestart bool
}
