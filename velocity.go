package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const version = "0.0.1"

func findVelocity(path string, verbose *bool, rs *int) bool {
	rs1 := []byte{0x75, 0x0C, 0xBC, 0xA3, 0x3A, 0x07, 0x8A, 0x41}
	rs2 := []byte{0x75, 0x7C, 0xBC, 0xA3, 0x3A, 0x07, 0x8A, 0x41} // and up

	f, err := ioutil.ReadFile(path)
	if err != nil && *verbose {
		fmt.Println("Could not process file:", err)
	}

	pattern := rs1
	if *rs == 2 {
		pattern = rs2
	}

	if strings.Contains(string(f), string(pattern)) {
		return true
	}
	return false
}

func main() {
	dir := flag.String("d", ".", "Directory to search in. Uses current directory by default.")
	rs := flag.Int("rs", 1, "Redstone variant to run against. Valid values: 1/2 (used in >=RS2)")
	verbose := flag.Bool("v", false, "Show verbose output.")
	ver := flag.Bool("ver", false, "Show version and exit.")
	flag.Parse()

	if *ver {
		fmt.Printf("velocitysearch v%s by Daniel Gurney", version)
		return
	}

	if *rs > 2 {
		*rs = 2
	}

	filepath.Walk(*dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil && *verbose {
				// We never return an error since we want to keep going.
				fmt.Println(err)
			}
			if strings.Contains(path, "WinSxS") || strings.Contains(path, "Microsoft.NET") || strings.Contains(path, "SysWOW64") {
				return nil
			}
			name := info.Name()
			if filepath.Ext(name) == ".exe" || filepath.Ext(name) == ".dll" {
				if findVelocity(path, verbose, rs) {
					fmt.Printf("Velocity signature found in %s :D\n", path)
				}
			}
			return nil
		})
}
