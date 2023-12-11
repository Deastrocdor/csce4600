package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/Deastrocdor/csce4600/builtins"
)

func main() {
	exit := make(chan struct{}, 2) // buffer this so there's no deadlock.
	runLoop(os.Stdin, os.Stdout, os.Stderr, exit)
}

func runLoop(r io.Reader, w, errW io.Writer, exit chan struct{}) {
	var (
		input    string
		err      error
		readLoop = bufio.NewReader(r)
	)
	for {
		select {
		case <-exit:
			_, _ = fmt.Fprintln(w, "exiting gracefully...")
			return
		default:
			if err := printPrompt(w); err != nil {
				_, _ = fmt.Fprintln(errW, err)
				continue
			}
			if input, err = readLoop.ReadString('\n'); err != nil {
				_, _ = fmt.Fprintln(errW, err)
				continue
			}
			if err = handleInput(w, input, exit); err != nil {
				_, _ = fmt.Fprintln(errW, err)
			}
		}
	}
}

func printPrompt(w io.Writer) error {
	// Get current user.
	// Don't prematurely memoize this because it might change due to `su`?
	u, err := user.Current()
	if err != nil {
		return err
	}
	// Get current working directory.
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	// /home/User [Username] $
	_, err = fmt.Fprintf(w, "%v [%v] $ ", wd, u.Username)

	return err
}

func handleInput(w io.Writer, input string, exit chan<- struct{}) error {
	// Remove trailing spaces.
	input = strings.TrimSpace(input)

	// Split the input separate the command name and the command arguments.
	args := strings.Split(input, " ")
	name, args := args[0], args[1:]

	// Check for built-in commands.
	// New builtin commands should be added here. Eventually this should be refactored to its own func.
	switch name {
	case "cd":
		return builtins.ChangeDirectory(args...)
	case "env":
		return builtins.EnvironmentVariables(w, args...)
	case "exit":
		exit <- struct{}{}
		return nil
	case "echo": // display text
		return builtins.Echo(w, args...)
	case "ls":
		//assumes directory path needs handling
		var dirPath string
		if len(args) > 0 {
			dirPath = args[0]
		}
		return builtins.ListDirectories(w, dirPath)
	case "pwd": // print
		return builtins.PrintWorkingDirectory(w)
	case "mkdir":
		// assumes directory path is first arg
		var mode os.FileMode
		if len(args) > 1 {
			modeInt, err := strconv.ParseUint(args[1], 8, 32)
			if err != nil {
				return fmt.Errorf("invalid mode: %s", args[1])
			}
			mode = os.FileMode(modeInt)
		}
		if len(args) > 0 {
			return builtins.MakeDirectory(args[0], mode)
		} else {
			return fmt.Errorf("mkdir: missing operand")
		}
	case "rmdir"://assumes directory path is the first arg
		if len(args) > 0 {
			return builtins.RemoveDirectory(args[0])
		} else {
			return fmt.Errorf("rmdir: missing operand")
		}
	default:
		return executeCommand(name, args...)
	}
}	

func executeCommand(name string, arg ...string) error {
	// Otherwise prep the command
	cmd := exec.Command(name, arg...)

	// Set the correct output device.
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	// Execute the command and return the error.
	return cmd.Run()
}
