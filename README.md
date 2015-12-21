# autoreload

This tool provides the functionality to automatically reload a go application upon source code changes.

### About
There are two parts to getting autoreload to work:

1. **autoreload** (inside docker or your code, for example)

	`autoreload` works with your code. It's responsible for rebuilding your binary
	whenever it receives data (what's in the data is irrelevant, it just needs to be something)
	through the **reset port** (default: `12345`).


2. **watcher** (should be run outside docker/your application, on your machine for example)
	
	The `watcher.sh` script watches for `.go` file changes in the project directory, telling the server to
	restart whenever a change is detected. Runs [nc/netcat](http://linux.die.net/man/1/nc) command to send data
	to the **reset port**, which informs `autoreload` to rebuild

	If you're running with Docker: _This should be run outside the docker
	container, because file system events are not properly detected inside the
	boot2docker VM: https://github.com/boot2docker/boot2docker/issues/688_

* `watcher.sh` will continue to watch the file system and contact the **reset port** in a loop until cancelled.

* `autoreload` will continue to listen on the **reset port** and rebuild the binary whenever contacted in a loop until cancelled.

***

### Requirements

##### Golang
Currently tested on `golang 1.4` and `golang 1.5`, but presumably works on other versions

##### [nc/netcat](http://linux.die.net/man/1/nc)

##### File system watchers
* For Macs, [fswatch](https://github.com/emcrisostomo/fswatch) is required to run `watcher.sh`.
* For other systems, anything that can watch the file system and send a packet to port 12345 will work.

### Required Ports

Port `12345` must be open for `watcher.sh` to communicate with `autoreload`.

Port `9000` must be open if you're running the godocs server.

If you're running Docker, see [here](https://docs.docker.com/engine/userguide/networking/default_network/dockerlinks/) for how to open ports between Docker and host.

***

### How it works

1. Install `autoreload`

	```
	go get github.com/graciouseloise/autoreload
	```

2. Set up your application autoreload
	```
	cd app/your-go-app
	$GOPATH/bin/autoreload
	```

3. Run the watcher

	```
	./watcher.sh
	```

4. Change some of your source code and watch as the binary is reloaded

***

### Get Help

```
$GOPATH/bin/autoreload --help

usage: autoreload [<flags>]

Autoreload functionality

Flags:
      --help            Show context-sensitive help (also try --help-long and --help-man).
  -d, --godocs          Run godocs server
  -c, --command=serve   Command arguments for running
  -b, --build=./binary  Target file to build to and run
  -f, --file=main.go    File name to build against
```

***

### Examples

**You can combine any of the below flags**

To run with a local godocs server (default: `False`)	
```
$GOPATH/bin/autoreload -d
```

To specify the filename of the binary being built (default: `./binary`)
```
$GOPATH/bin/autoreload --build=./my-app
```

To specify the file to build against (default: No filename)
```
$GOPATH/bin/autoreload --file=main.go
```

To run with command line arguments for your application (default: None)
```
$GOPATH/bin/autoreload -c run
```

