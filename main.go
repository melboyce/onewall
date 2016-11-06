package main

import "flag"
import "fmt"
import "io/ioutil"
import "math/rand"
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
	//stanza := strings.Join(bgs[:*numHeads-1], " ")
	stanza := "#!/bin/sh\nfeh --bg-fill"
	bgs = bgs[len(bgs)-*numHeads-1:]

	if *pos > len(bgs) - 1 {
    	panic("pos > heads")
	}

	patt := filepath.Join(*dir, "*")

	walls, err := filepath.Glob(patt)
	if err != nil {
    	panic(err)
	}

	wallList := " "
	for i, f := range bgs {
    	if i == *pos {
        	wallList += "'" + walls[rand.Intn(len(walls))] + "'"
    	} else {
        	wallList += f
    	}
    	wallList += " "
	}

	fmt.Println(stanza + wallList)
}
