package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
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
	landscape := flag.Bool("l", false, "only landscape images")
	portrait := flag.Bool("p", false, "only portrait images")

	flag.Usage = usage
	flag.Parse()
	if flag.NArg() < 1 {
		usage()
		os.Exit(64)
	}
	dirs := flag.Args()

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

	nwall, err := getWall(dirs, *landscape, *portrait)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cant get wall from %s: %v\n", dirs, err)
		os.Exit(4)
	}

	if *pos > len(walls)-1 {
		walls = append(walls, nwall)
	} else {
		walls[*pos] = nwall
	}

	walls = append([]string{"--bg-fill"}, walls...)

	cmd := exec.Command("feh", walls...)
	err = cmd.Run()
	if err != nil {
		fmt.Printf("%+v\n", walls)
		fmt.Fprintf(os.Stderr, "cant run feh: %v\n", err)
		os.Exit(5)
	}
	os.Exit(0)
}

func getWall(dd []string, l, p bool) (string, error) {
	walls := []string{}
	for _, d := range dd {
		ww, err := filepath.Glob(filepath.Join(d, "*"))
		if err != nil {
			return "", fmt.Errorf("cant get walls: %w", err)
		}
		walls = append(walls, ww...)
	}
	rand.Shuffle(len(walls), func(i, j int) { walls[i], walls[j] = walls[j], walls[i] })
	for _, wall := range walls {
		ext := filepath.Ext(wall)
		switch ext {
		case ".jpg", ".jpeg", ".png", ".gif":
		default:
			continue
		}
		f, err := os.Open(wall)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cant open file %s: %v\n", wall, err)
			continue
		}
		defer f.Close()
		img, _, err := image.DecodeConfig(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cant decode image %s: %v\n", wall, err)
			continue
		}
		if !l && !p {
			return wall, nil
		}
		if l && img.Width > img.Height {
			return wall, nil
		}
		if p && img.Width < img.Height {
			return wall, nil
		}
	}
	return "", errors.New("bottomed out")
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
	fmt.Printf("usage: %s [-pos N] [-l|-p] <dir> [<dir> <dir> ...]\n", os.Args[0])
	flag.PrintDefaults()
}
