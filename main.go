package main

// autoreload builds and starts the go program, and waits for connections on
// resetPort. If a connection is opened on the port, the server is killed,
// rebuilt, and started up again. This is intended to be used with watcher.sh,
// which watches for file system changes and sends the reset signal on
// resetPort.

import (
	"fmt"
	"net"
	"os"
	"os/exec"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app    = kingpin.New("autoreload", "Autoreload functionality")
	godocs = app.Flag("godocs", "Run godocs server").Short('d').Bool()

	cmd         = app.Flag("command", "Command arguments for running").PlaceHolder("serve").Short('c').String()
	targetBuild = app.Flag("build", "Target file to build to and run").PlaceHolder("./binary").Default("./build").String()
	targetFile  = app.Flag("file", "File name to build against").PlaceHolder("main.go").String()
	resetPort   = app.Flag("reset-port", "Port to listen for reset signal").PlaceHolder("12345").Short('r').Default("12345").Int()
)

const godocPort = 9000

func serveGodoc() *exec.Cmd {
	fmt.Println("Starting godoc server on port", godocPort)
	cmd := exec.Command("godoc", "-http", fmt.Sprintf(":%d", godocPort))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error starting godoc server!", err)
		return nil
	}
	return cmd
}

func buildAndServe() (*exec.Cmd, error) {
	fmt.Printf("Rebuilding %s\n", *targetBuild)
	build := exec.Command("go", "build", "-v", "-o", *targetBuild)
	if *targetFile != "" {
		build = exec.Command("go", "build", "-v", "-o", *targetBuild, *targetFile)
	}
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr
	err := build.Run()
	if err != nil {
		fmt.Println("Error building application!", err)
		return nil, err
	}

	fmt.Println("Serving")
	app := exec.Command(*targetBuild)
	if *cmd != "" {
		app = exec.Command(*targetBuild, *cmd)
	}
	app.Stdout = os.Stdout
	app.Stderr = os.Stderr
	err = app.Start()
	if err != nil {
		fmt.Println("Error starting application!", err)
		return nil, err
	}

	return app, nil
}

func checkError(err error) {
	if err != nil {
		fmt.Errorf("%v", err)
	}
}

func main() {
	_ = kingpin.MustParse(app.Parse(os.Args[1:]))
	if *godocs {
		godoc := serveGodoc()
		if godoc != nil {
			defer godoc.Process.Kill()
		}
	}
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", *resetPort))
	checkError(err)
	defer l.Close()

	for {
		app, err := buildAndServe()
		checkError(err)

		fmt.Printf("Waiting for reset signal on %s ...\n", l.Addr())
		conn, err := l.Accept()
		checkError(err)
		defer conn.Close()
		fmt.Println("Received reset signal")

		if app != nil {
			fmt.Println("Killing the application")
			err := app.Process.Kill()
			checkError(err)
		}
	}
}
