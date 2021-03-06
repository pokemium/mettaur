package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pokemium/magia/pkg/emulator"
	"github.com/pokemium/magia/pkg/emulator/audio"
	"github.com/pokemium/magia/pkg/emulator/joypad"
	"github.com/pokemium/magia/pkg/gba"

	"github.com/hajimehoshi/ebiten/v2"
)

var version string

const (
	title   = "Magia"
	exeName = "magia"
)

// ExitCode represents program's status code
type ExitCode int

// exit code
const (
	ExitCodeOK ExitCode = iota
	ExitCodeError
)

func init() {
	if version == "" {
		version = "Develop"
	}

	flag.Usage = func() {
		usage := fmt.Sprintf(`Usage:
    %s [arg] [input]
input: a filepath
Arguments: 
`, exeName)

		fmt.Fprintf(os.Stderr, usage)
		flag.PrintDefaults()
	}
}

func main() {
	os.Exit(int(Run()))
}

// Run program
func Run() ExitCode {
	var (
		showVersion   = flag.Bool("v", false, "show version")
		debug         = flag.Bool("d", false, "exec in debug mode")
		showBIOSIntro = flag.Bool("b", false, "show BIOS intro")
		showCartInfo  = flag.Bool("c", false, "show cartridge info")
		mute          = flag.Bool("m", false, "mute sound")
	)

	flag.Parse()
	if *showVersion {
		fmt.Println(title+":", version)
		return ExitCodeOK
	}

	path := flag.Arg(0)
	data, err := readROM(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read ROM data: %s\n", err)
		return ExitCodeError
	}

	emu := emulator.New(gba.New(data, &audio.Stream, *debug, *mute), path)
	if *showCartInfo {
		fmt.Println(emu.GBA.CartInfo())
		return ExitCodeOK
	}

	emu.GBA.SetJoypadHandler(joypad.Handler)
	emu.LoadSav()
	if *showBIOSIntro {
		emu.GBA.Reset()
	} else {
		emu.GBA.SoftReset()
	}

	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle(emu.GBA.CartHeader.Title)
	ebiten.SetWindowSize(240*2, 160*2)
	if err := ebiten.RunGame(emu); err != nil {
		fmt.Fprintf(os.Stderr, "crash in emulation: %s\n", err)
	}
	return ExitCodeOK
}

func readROM(path string) ([]byte, error) {
	if path == "" {
		return []byte{}, errors.New("please select gba file path")
	}
	if filepath.Ext(path) != ".gba" {
		return []byte{}, errors.New("please select .gba file")
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return []byte{}, errors.New("fail to read file")
	}
	return bytes, nil
}
