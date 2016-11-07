package main

import "flag"
//import "fmt"
import "io/ioutil"
import "math/rand"
import "os/exec"
import "path/filepath"
import "strings"
import "time"

func init() {
    rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
    pos := flag.Int("pos", 0, "position (head index)")
    dir := flag.String("dir", "/home/mel/i/walls/default", "directory to parse")
    numHeads := flag.Int("numHeads", 4, "number of heads")

    flag.Parse()

	fehbg, err := ioutil.ReadFile("/home/mel/.fehbg")
	if err != nil {
    	panic(err)
	}

	bgs := strings.Split(string(fehbg), " ")
	bgs = bgs[len(bgs)-*numHeads-1:]

	if *pos > len(bgs) - 1 {
    	panic("pos > heads")
	}

	patt := filepath.Join(*dir, "*")

	walls, err := filepath.Glob(patt)
	if err != nil {
    	panic(err)
	}

	fehArgs := []string{"--bg-fill"}
	for i, w := range bgs {
    	w = strings.TrimSpace(w)
    	if len(w) > 0 {
        	if i == *pos {
            	fehArgs = append(fehArgs, walls[rand.Intn(len(walls))])
        	} else {
            	fehArgs = append(fehArgs, strings.Replace(w, "'", "", -1))
        	}
    	}
	}

	cmd := exec.Command("feh", fehArgs...)

	err = cmd.Run()
	if err != nil {
    	panic(err)
	}
}
