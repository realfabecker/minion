package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

func runSingleCommand(command string, flags map[string]string) error {
	act := strings.Split(command, " ")
	if len(flags) > 0 {
		for k, v := range flags {
			act = append(act, fmt.Sprintf("--%s=%s", k, v))
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

func runParallelCommandsWithFlags(command string, pll int, mFlags []map[string]string) {
	var wg sync.WaitGroup
	wg.Add(len(mFlags))

	start := time.Now()

	fmt.Println("Running", pll, "commands in parallel")
	var ch = make(chan struct{}, pll)
	defer close(ch)

	for _, v := range mFlags {
		go func(flags map[string]string) {
			defer wg.Done()
			ch <- struct{}{}
			if err := runSingleCommand(command, flags); err != nil {
				fmt.Println(fmt.Errorf("runPll: %w", err))
			}
			<-ch
		}(v)
	}
	wg.Wait()

	end := time.Since(start)
	fmt.Println("Done in", end.String(), "(", end.Seconds(), "s)")
}

func readFlagsFromCsv(filePath string) ([]map[string]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	csvReader := csv.NewReader(f)
	header, err := csvReader.Read()

	var args = make([]map[string]string, 0)
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		var m = make(map[string]string)
		for i, v := range rec {
			m[header[i]] = v
		}
		args = append(args, m)
	}
	return args, nil
}

func runParallelCommands(command string, pll int) {
	var wg sync.WaitGroup
	wg.Add(pll)
	start := time.Now()

	fmt.Println("Running", pll, "commands in parallel")
	var ch = make(chan struct{}, pll)
	defer close(ch)
	for i := 0; i < pll; i++ {
		go func() {
			defer wg.Done()
			ch <- struct{}{}
			if err := runSingleCommand(command, nil); err != nil {
				fmt.Println(fmt.Errorf("runPll: %w", err))
			}
			<-ch
		}()
	}
	wg.Wait()
	end := time.Since(start)
	fmt.Println("Done in", end.String(), "(", end.Seconds(), "s)")
}

func main() {
	var workers = flag.Int("w", 1, "number of concurrent workers")
	var command = flag.String("c", "", "command to be executed")
	var argFile = flag.String("f", "", "path to csv file with flags")

	flag.Parse()
	if (*argFile) != "" && (*command) != "" {
		mFlags, err := readFlagsFromCsv(*argFile)
		if err != nil {
			log.Fatalln("error reading csv file:", err)
		}
		runParallelCommandsWithFlags(*command, *workers, mFlags)
		fmt.Println("Done")
	} else if *command != "" {
		runParallelCommands(*command, *workers)
		fmt.Println("Done")
	} else {
		flag.PrintDefaults()
	}
}
