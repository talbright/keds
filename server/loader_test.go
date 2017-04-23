package server_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/talbright/keds/server"
	us "github.com/talbright/keds/utils/system"

	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

var pluginContent = []byte("#!/bin/bash\n echo true > $0.out")

func createMockPlugin(rootDir string, name string, data []byte, perm os.FileMode) (dir string) {
	dir = filepath.Join(rootDir, name)
	if err := os.Mkdir(dir, 0755); err != nil {
		log.Fatal(err)
	}
	tmpfn := filepath.Join(dir, name)
	if data != nil {
		if err := ioutil.WriteFile(tmpfn, data, perm); err != nil {
			log.Fatal(err)
		}
	}
	return tmpfn
}

func createMockPlugins() (rootDir string, dirs []string) {
	dirs = make([]string, 0)
	if dir, err := ioutil.TempDir("", "plugins"); err != nil {
		log.Fatal(err)
	} else {
		rootDir = dir
	}
	for i := 0; i < 3; i++ {
		dirs = append(dirs, createMockPlugin(rootDir, fmt.Sprintf("plugin%d", i), pluginContent, 0755))
	}
	dirs = append(dirs, createMockPlugin(rootDir, "plugin3", pluginContent, 0644))
	dirs = append(dirs, createMockPlugin(rootDir, "plugin4", nil, 0644))
	return
}

var _ = Describe("Loader", func() {
	Describe("FindCmdInPath", func() {
		It("should load and run the executables", func() {
			rootPluginDir, _ := createMockPlugins()
			l := NewLoader([]string{rootPluginDir})
			Expect(l.Load()).Should(Succeed())
			for i := 0; i < 3; i++ {
				plugin := fmt.Sprintf("plugin%d", i)
				out := path.Join(rootPluginDir, plugin, plugin+".out")
				Eventually(func() bool {
					return us.Exists(out)
				}).Should(BeTrue())
			}
		})
	})
	Describe("FindCmdInPath", func() {
		It("should find executables in subdirs of path", func() {
			rootPluginDir, _ := createMockPlugins()
			l := NewLoader([]string{rootPluginDir})
			cmds := l.FindCmdsInPath(rootPluginDir)
			Expect(cmds).Should(HaveLen(3))
			Expect(cmds).Should(ContainElement(filepath.Join(rootPluginDir, "plugin0", "plugin0")))
			Expect(cmds).Should(ContainElement(filepath.Join(rootPluginDir, "plugin1", "plugin1")))
			Expect(cmds).Should(ContainElement(filepath.Join(rootPluginDir, "plugin2", "plugin2")))
		})
	})
	Describe("BuildCmd", func() {
		It("should create the command to execute", func() {
			rootPluginDir, _ := createMockPlugins()
			l := NewLoader([]string{rootPluginDir})
			cmds := l.FindCmdsInPath(rootPluginDir)
			cmd := l.BuildCmd(cmds[0])
			Expect(cmd.Path).Should(Equal(cmds[0]))
			Expect(cmd.Dir).Should(Equal(path.Dir(cmds[0])))
		})
	})
	Describe("StartCmd", func() {
		It("should start the command", func() {
			rootPluginDir, _ := createMockPlugins()
			l := NewLoader([]string{rootPluginDir})
			cmds := l.FindCmdsInPath(rootPluginDir)
			cmd := l.BuildCmd(cmds[0])
			Expect(l.StartCmd(cmd)).Should(Succeed())
			out := cmds[0] + ".out"
			Eventually(func() bool {
				return us.Exists(out)
			}).Should(BeTrue())
		})
	})
})
