package francis

import (
	"bytes"
	"fmt"
	"github.com/bnagy/crashwalk/crash"
	"go/build"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func (e *Engine) Run(command []string) (crash.Info, error) {

	pkg, err := build.Import("github.com/bnagy/francis", ".", build.FindOnly)
	if err != nil {
		return crash.Info{}, fmt.Errorf("Couldn't find import path: %s", err)
	}
	tool := filepath.Join(pkg.Dir, "exploitaben/exploitaben.py")

	// Construct the command array
	// TODO LINUX - we don't have MallocScribble, is there an easy equivalent?
	cmdSlice := []string{tool}
	if e.Timeout > 0 {
		cmdSlice = append(cmdSlice, []string{"-t", strconv.Itoa(e.Timeout)}...)
	}
	cmdSlice = append(cmdSlice, "--")
	cmdSlice = append(cmdSlice, command...)
	cmdStr := strings.Join(cmdSlice, " ")

	cmd := exec.Command(cmdSlice[0], cmdSlice[1:]...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return crash.Info{}, fmt.Errorf("Error creating pipe: %s", err)
	}
	if err := cmd.Start(); err != nil {
		return crash.Info{}, fmt.Errorf("Error launching tool: %s", err)
	}

	out, _ := ioutil.ReadAll(stdout)
	err = cmd.Wait()

	// handles clean exit, exit with errorcode ( no crash ) or timeout
	if len(out) == 0 ||
		bytes.Contains(out, []byte("exited with status")) ||
		bytes.Contains(out, []byte("killing the process...")) {
		// No crash.
		return crash.Info{}, fmt.Errorf("No lldb output for %s", cmdStr)
	}

	ci := parse(out, cmdStr)
	// If, for some reason, we managed to parse the instrumentation output
	// without crashing, but there's an error from running the command we
	// return both. The caller can abort if they're being strict or just
	// ignore it otherwise.
	return ci, err

}
