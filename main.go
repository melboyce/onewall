// update a .fehbg file and refresh the xroot.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func bail(s string, e int) {
	fmt.Fprintf(os.Stderr, s+"\n\n")
	flag.Usage()
	os.Exit(e)
}

func usage() {
	fmt.Printf("usage: %s [-pos] <directory>\n", os.Args[0])
	flag.PrintDefaults()
}

var pos = flag.Int("pos", 0, "position (xinerama screen index)")

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() < 1 {
		bail("no directory specified", 64)
	}

	dotfeh := fehBG()
	if _, err := os.Stat(dotfeh); err != nil {
		bail(err.Error(), 1)
	}
	b, err := ioutil.ReadFile(dotfeh)
	if err != nil {
		bail(err.Error(), 1)
	}
	curr := strings.Split(strings.TrimSpace(string(b)), " ")
	curr = curr[3:] // TODO process each line to drop the shebang properly
	if *pos > len(curr) {
		bail("position greater than number of entries in .fehbg", 1)
	}
	for i, fp := range curr {
		curr[i] = strings.Replace(fp, "'", "", -1)
	}

	dir := flag.Arg(0)
	walls, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		bail(err.Error(), 1)
	}
	curr[*pos] = walls[rand.Intn(len(walls))]

	args := []string{"--bg-fill"}
	for _, fp := range curr {
		args = append(args, fp)
	}
	cmd := exec.Command("feh", args...)
	err = cmd.Run()
	if err != nil {
		bail(err.Error(), 2)
	}
}

func fehBG() string {
	u, err := user.Current()
	if err != nil {
		bail(err.Error(), 1)
	}
	return filepath.Join(u.HomeDir, ".fehbg")
}
