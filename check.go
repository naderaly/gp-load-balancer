package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

var defaultFailedCode = 1

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func main() {
	//argsWithProg := os.Args
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 0 {
		fmt.Println("error : args not found")
	} else {
		//arg := os.Args[3]
		Check(argsWithoutProg)

	}

}

// Check icap tacke 4 args
func Check(arg []string) {
	// 192.168.100.100 80 192.168.100.0 80

	//fmt.Println(arg[2])
	//fmt.Println(arg[3])

	servername := arg[2]
	port := arg[3]
	lastrg := len(arg) - 1

	filename := "test.pdf"
	rebfile := "reb_" + filename
	if fileExists(filename) {
		os.Remove(rebfile)
	}

	_ = filename

	runcmd := []string{"c-icap-client", "-i", servername, "-p", port, "-f", filename, "-s", "gw_rebuild", "-o", "reb_" + filename, "-v"}
	//runcmd := []string{"c-icap-client", "-i", servername, "-p", port, "-v"}
	//c-icap-client -i eu.icap.glasswall-icap.com -p 1344 -f test.pdf -s gw_rebuild -o reh.pdf
	s := run(10, arg[lastrg], "time", runcmd...)
	if s == "0" {
		if arg[lastrg] == "-v" {
			fmt.Printf("exitcode 0")
		}
		os.Exit(0)

	} else {
		if arg[lastrg] == "-v" {
			fmt.Printf(s)
		}
		os.Exit(1)
	}

}
func run(timeout int, lastarg string, command string, args ...string) string {

	// instantiate new command

	cmd := exec.Command(command, args...)

	// get pipe to standard output
	stdout, err := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err != nil {
		return "cmd.StdoutPipe() error: " + err.Error()
	}

	// start process via command
	if err := cmd.Start(); err != nil {
		return "cmd.Start() error: " + err.Error()
	}

	slurp, _ := ioutil.ReadAll(stderr)
	if lastarg == "-v" {
		fmt.Printf("%s\n", slurp)
	}
	outlines := strings.Split(string(slurp), "\n")
	l := len(outlines)
	req := false
	for _, line := range outlines[1 : l-1] {
		//parsedLine := strings.Fields(line)
		if strings.Contains(line, "HTTP/1.0 200 OK") == true {
			req = true
		}

	}
	if req == false {
		return "Service not run "
	}

	// setup a buffer to capture standard output
	var buf bytes.Buffer

	// create a channel to capture any errors from wait
	done := make(chan error)
	go func() {
		if _, err := buf.ReadFrom(stdout); err != nil {
			panic("buf.Read(stdout) error: " + err.Error())
		}

		done <- cmd.Wait()
	}()

	// block on select, and switch based on actions received
	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		if err := cmd.Process.Kill(); err != nil {
			return "failed to kill: " + err.Error()
		}
		return "1"
	case err := <-done:
		if err != nil {
			close(done)
			return "exitcode: " + err.Error()
		}

		return "0" + buf.String()
	}
	return ""
}

//RunCommand  run cmd icap
func RunCommand(name string, args []string) (stdout string, stderr string, exitCode int) {

	log.Println("run command:", name, args)
	var outbuf, errbuf bytes.Buffer
	//cmd := exec.Command(name)
	cmd := exec.Command(name, args...)
	//cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	_, err := cmd.Output()
	stdout = outbuf.String()
	stderr = errbuf.String()

	if err != nil {
		// try to get the exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
		} else {
			// This will happen (in OSX) if `name` is not available in $PATH,
			// in this situation, exit code could not be get, and stderr will be
			// empty string very likely, so we use the default fail code, and format err
			// to string and set to stderr
			log.Printf("Could not get exit code for failed program: %v, %v", name, args)

			exitCode = defaultFailedCode
			if stderr == "" {
				stderr = err.Error()
			}
		}
	} else {
		// success, exitCode should be 0 if go is ok
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
	}
	log.Printf("command result, stdout: %v, stderr: %v, exitCode: %v", stdout, stderr, exitCode)
	done := make(chan error)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		log.Println("error:", err)
		// exited
	case <-time.After(10 * time.Second):
		log.Println("error:", "time out")
		// timed out
	}
	return
}
