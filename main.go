package main

import (
    "flag"
    "strings"
    "time"

    "io/ioutil"
    "math/rand"
    "os/exec"
    "path/filepath"
)

func init() {
    rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
    pos := flag.Int("pos", 0, "position (head index)")
    dir := flag.String("dir", "", "directory to parse")
    heads := flag.Int("heads", 1, "number of heads")

    flag.Parse()

    fehbg, err := ioutil.ReadFile("~/.fehbg")
    if err != nil {
        panic(err)
    }

    bgs := strings.Split(string(fehbg), " ")
    bgs = bgs[len(bgs) - *heads - 1:]

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
