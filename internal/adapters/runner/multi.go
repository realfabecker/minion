package runner

import (
	"fmt"
	"github.com/realfabecker/kevin/internal/core/ports"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

type multi struct{}

func NewMulti() ports.ParallelRunner {
	return &multi{}
}

func (m *multi) Run(command string, pll int, mFlags []map[string]string) {
	start := time.Now()

	fmt.Println("Running", pll, "commands in parallel")
	if len(mFlags) > 0 {
		m.runParallelWithFlags(command, pll, mFlags)
	} else {
		m.runParallel(command, pll)
	}

	end := time.Since(start)
	fmt.Println("Done in", end.String(), "(", end.Seconds(), "s)")
}

func (m *multi) runParallel(command string, pll int) {
	var wg sync.WaitGroup
	wg.Add(pll)

	var ch = make(chan struct{}, pll)
	defer close(ch)
	for i := 0; i < pll; i++ {
		go func() {
			defer wg.Done()
			ch <- struct{}{}
			if err := m.runCmd(command, nil); err != nil {
				fmt.Println(fmt.Errorf("runPll: %w", err))
			}
			<-ch
		}()
	}
	wg.Wait()
}

func (m *multi) runParallelWithFlags(command string, pll int, mFlags []map[string]string) {
	var wg sync.WaitGroup
	wg.Add(len(mFlags))

	var ch = make(chan struct{}, pll)
	defer close(ch)

	for _, v := range mFlags {
		go func(flags map[string]string) {
			defer wg.Done()
			ch <- struct{}{}
			if err := m.runCmd(command, flags); err != nil {
				fmt.Println(fmt.Errorf("runPll: %w", err))
			}
			<-ch
		}(v)
	}
	wg.Wait()
}

func (m *multi) runCmd(command string, flags map[string]string) error {
	act := strings.Split(command, " ")
	if len(flags) > 0 {
		for k, v := range flags {
			act = append(act, fmt.Sprintf(`--%s='%s'`, k, v))
		}
	}
	var cmd *exec.Cmd
	if runtime.GOOS == "linux" {
		cmd = exec.Command("bash", "-c", strings.Join(act, " "))
	} else {
		cmd = exec.Command(act[:1][0], act[1:]...)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
