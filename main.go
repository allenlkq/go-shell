package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"flag"
	"strings"
	"os/exec"
	"time"
	"sync"
)

var numberOfProcesses = flag.Int("n", 5, "number of processes")
var sleep = flag.Int("s", 10, "sleep beween tasks per process, in seconds")
var loop = flag.Bool("l", false, "loop the commands")

func Map(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func Filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func main() {
	flag.Parse()

	bytes, _ := ioutil.ReadAll(os.Stdin)
	commands := strings.Split(string(bytes),"\n")
	commands = Map(commands, func(s string) string{
		return strings.Trim(s, " ")
	})
	commands = Filter(commands, func(s string) bool{
		return s != ""
	})
	//commands = []string{"ls", "pwd"}
	if len(commands) == 0 {
		fmt.Println("No commands to run\n")
		os.Exit(1);
	}

	fmt.Println("Commands to run: \n" + strings.Join(commands, "\n"))

	cmdChan := make(chan string)
	var wg sync.WaitGroup

	for i:=1; i<=*numberOfProcesses; i++ {
		go func(cmdChan <- chan string) {
			for cmd := range cmdChan{
				fmt.Printf("running command %s \n", cmd)
				out, err := exec.Command("sh","-c",cmd).Output()
				if err != nil {
					fmt.Printf("ERR: %q \n", err)
				} else {
					fmt.Printf("OUT: %s \n", string(out))
				}
				wg.Done()
				time.Sleep(time.Duration(*sleep) * time.Second)
			}
		}(cmdChan)
	}

	for _, cmd := range commands {
		cmdChan <- cmd
		wg.Add(1)
	}
	for *loop {
		for _, cmd := range commands {
			cmdChan <- cmd
			wg.Add(1)
		}
	}

	wg.Wait()
}