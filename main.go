package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
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

func main() {
	pos := flag.Int("pos", 0, "position (xinerama index)")
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() < 1 {
		usage()
		os.Exit(64)
	}
	dir := flag.Arg(0)

	user, err := user.Current()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cant get user: %v\n", err)
		os.Exit(1)
	}

	fehbg := filepath.Join(user.HomeDir, ".fehbg")
	f, err := os.Open(fehbg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cant open %s: %v\n", fehbg, err)
		os.Exit(2)
	}
	defer f.Close()

	feh, err := getCmd(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing %s: %v\n", fehbg, err)
		os.Exit(3)
	}

	walls := []string{}
	els := strings.Fields(feh)
	for _, el := range els {
		if el == "feh" || strings.HasPrefix(el, "-") {
			continue
		}
		el = strings.Trim(el, "'")
		walls = append(walls, el)
	}

	nwall, err := getWall(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cant get wall from %s: %v\n", dir, err)
		os.Exit(4)
	}

	if *pos > len(walls)-1 {
		walls = append(walls, nwall)
	} else {
		walls[*pos] = nwall
	}

	args := []string{"--bg-fill"}
	for _, wall := range walls {
		args = append(args, wall)
	}
	cmd := exec.Command("feh", args...)
	err = cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cant run feh: %v\n", err)
		os.Exit(5)
	}
	os.Exit(0)
}

func getWall(d string) (string, error) {
	walls, err := filepath.Glob(filepath.Join(d, "*"))
	if err != nil {
		return "", fmt.Errorf("cant get walls: %w", err)
	}
	return walls[rand.Intn(len(walls))], nil
}

func getCmd(f *os.File) (string, error) {
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "feh ") {
			return line, nil
		}
	}
	if err := sc.Err(); err != nil {
		return "", fmt.Errorf("error scanning file: %w", err)
	}
	return "", errors.New("cant find feh command in file")
}

func usage() {
	fmt.Printf("usage: %s [-pos N] <directory>\n", os.Args[0])
	flag.PrintDefaults()
}
