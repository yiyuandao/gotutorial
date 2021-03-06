package main

import "github.com/cmceniry/gotutorial/mapping"
import "bufio"
import "fmt"
import "os/exec"
import "regexp"
import "strings"

func substitute(line string) string {
	for orig, replace := range mapping.Subs {
		re := regexp.MustCompile(orig + "(\\s+)")
		line = re.ReplaceAllString(line, replace+strings.Repeat(" ", 16-len(orig)))
	}
	return line
}

func cmdExec() {
	args := []string{"-xt", "/dev/loop1", "/dev/loop2", "/dev/loop3", "/dev/loop4", "5"}
	cmd := exec.Command("iostat", args...)
	stdout, err := cmd.StdoutPipe()
	out := bufio.NewReader(stdout)
	if err != nil {
		panic(err)
	}
	if err := cmd.Start(); err != nil {
		panic(err)
	}

	for {
		if buf, err := out.ReadString('\n'); err != nil {
			panic(err)
		} else {
			fmt.Print(substitute(buf))
		}
	}
}

func main() {
	firstupdate := make(chan bool) // (1)
	go mapping.SignalUpdater(firstupdate) // (2)
	if <-firstupdate { // (3)
		cmdExec() // (4)
	} else {
		panic("Unable to find correct ASM/device mappings")
	}
}
