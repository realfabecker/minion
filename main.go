package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var dryRun bool

func runSingleCommand(command string, flags map[string]string) error {
	act := strings.Split(command, " ")
	for k, v := range flags {
		act = append(act, fmt.Sprintf("--%s=%s", k, v))
	}
	if dryRun {
		fmt.Println(strings.Join(act, " "))
		return nil
	}
	cmd := exec.Command(act[:1][0], act[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func runParallelCommands(command string, mFlags []map[string]string, pll int) {

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

func main() {
	var argFile = flag.String("f", "", "path to csv file with flags")
	var workers = flag.Int("n", 1, "number of concurrent workers")
	var command = flag.String("c", "", "command to be executed")
	flag.BoolVar(&dryRun, "d", false, "dry run mode")

	flag.Parse()
	if (*argFile) != "" && (*command) != "" {
		mFlags, err := readFlagsFromCsv(*argFile)
		if err != nil {
			log.Fatalln("error reading csv file:", err)
		}
		runParallelCommands(*command, mFlags, *workers)
		fmt.Println("Done")
	} else {
		flag.PrintDefaults()
	}
}
