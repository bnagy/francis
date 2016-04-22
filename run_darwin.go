package francis

import (
	"fmt"
	"github.com/bnagy/crashwalk/crash"
	"go/build"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// Run runs a command under francis. This API matches the interface for
// Debugger for some of my other tools, but the `filename` and `memlimit`
// params are not presently used. The caller should set them to "" and 0
func (e *Engine) Run(command []string, filename string, memlimit, timeout int) (crash.Info, error) {

	pkg, err := build.Import("github.com/bnagy/francis", ".", build.FindOnly)
	if err != nil {
		return crash.Info{}, fmt.Errorf("Couldn't find import path: %s", err)
	}
	tool := filepath.Join(pkg.Dir, "exploitaben/exploitaben.py")

	// Construct the command array
	// XXX: This doesn't do anything, I think I forgot to export an env var or something
	cmdSlice := []string{tool, "-e", "MallocScribble=1", "-e", "MallocGuardEdges=1"}
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
	cmd.Wait()

	return getCrashInfo(out, cmdStr)

}
