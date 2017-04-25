package main

import (
    "flag"
    "strings"
    "time"
    "os"

    "io/ioutil"
    "math/rand"
    "os/exec"
    "os/user"
    "path/filepath"
)

func init() {
    rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
    pos := flag.Int("pos", 0, "position (xinerama index)")
    dir := flag.String("dir", "", "directory to parse")
    heads := flag.Int("heads", 1, "number of heads")

    flag.Parse()

    if *dir == "" {
        flag.PrintDefaults()
        os.Exit(1)
    }

    usr, err := user.Current()
    if err != nil {
        panic(err)
    }

    // TODO handle non-existent .fehbg
    fehpath := filepath.Join(usr.HomeDir, ".fehbg")
    fehbg, err := ioutil.ReadFile(fehpath)
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
