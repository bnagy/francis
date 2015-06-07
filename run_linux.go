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
	cmd.Wait()

	return getCrashInfo(out, cmdStr)

}
