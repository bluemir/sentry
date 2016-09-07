#Sentry

Simple command line tool that execute command when file is changed.


## Installation
```
go get github.com/bluemir/sentry
```

## Usage
```
sentry -w "path/to/watch" -c "echo test" # run echo test when file changed 
```

```
Usage:
  sentry [flags]

Flags:
  -c, --command string    command to run
  -d, --delay value       delay that wait events (default 500)
  -x, --exclude string    exclude pattern(regexp)
  -k, --kill-on-restart   kill on restart (default true)
  -s, --shell string      shell to execute command (default "$SHELL")
  -v, --verbose           show detail log
      --version           show version
  -w, --watch value       paths to watch (default [./])
 ```

## Applcations

### go auto rebuild
```
cd /path/to/myapp
sentry -c "go build -o path/to/result"
```

### go auto test when save code
```
cd /path/to/myapp
sentry -c "go test ./..."
```

### dev server
```
cd /path/to/myapp
sentry -c "make server" # if you can run your server with makefile
```

## Contributing & Bug report
Please write issue on https://github.com/bluemir/sentry/issues or
Just send pull request to https://github.com/bluemir/sentry
