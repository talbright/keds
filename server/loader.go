package server

import (
	us "github.com/talbright/keds/utils/system"
	"golang.org/x/net/trace"

	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
)

type ILoader interface {
	Load() (err error)
}

type Loader struct {
	loadPath []string
	cmds     []*exec.Cmd
	events   trace.EventLog
}

func NewLoader(loadPath []string) *Loader {
	_, file, line, _ := runtime.Caller(0)
	return &Loader{
		loadPath: loadPath,
		cmds:     make([]*exec.Cmd, 0),
		events:   trace.NewEventLog("server.Loader", fmt.Sprintf("%s:%d", file, line)),
	}
}

func (l *Loader) Load() (err error) {
	executables := make([]string, 0)
	for _, p := range l.loadPath {
		l.events.Printf("searching for plugins in %s", p)
		executables = append(executables, l.FindPluginsInPath(p)...)
	}
	for _, e := range executables {
		if cmd := l.BuildCmd(e); cmd != nil {
			l.cmds = append(l.cmds, cmd)
		}
	}
	for _, c := range l.cmds {
		l.events.Printf("running plugin %s", c.Path)
		if err := l.StartCmd(c); err != nil {
			log.Fatal(err)
		}
	}
	return
}

func (l *Loader) FindPluginsInPath(path string) (cmds []string) {
	cmds = make([]string, 0)
	if files, err := ioutil.ReadDir(path); err == nil {
		for _, file := range files {
			ex := filepath.Join(path, file.Name(), file.Name())
			if us.IsExecutable(ex) {
				cmds = append(cmds, ex)
			}
		}
	}
	return
}

func (l *Loader) BuildCmd(exPath string) (cmd *exec.Cmd) {
	cmd = exec.Command(exPath)
	cmd.Dir = path.Dir(exPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return
}

func (l *Loader) StartCmd(cmd *exec.Cmd) (err error) {
	return cmd.Start()
}
