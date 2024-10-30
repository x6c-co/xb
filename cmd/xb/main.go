package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	birdsocket "github.com/czerwonk/bird_socket"
	"github.com/olekukonko/tablewriter"
)

var (
	socketPath  = "/var/run/bird/bird.ctl"
	debug       = false
	header      = []string{"Session", "State", "Neighbor", "AS", "Import", "Export"}
	borders     = tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false}
	version     = "alpha"
	commit      = ""
	date        = ""
	showVersion = false
)

const (
	bufferSize          = 10_000
	showProtocolsAllCMD = "show protocols all"
	newEntryKey         = "1002-"
	neighborASKey       = "Neighbor AS:"
	neighborAddressKey  = "Neighbor address:"
	importUpdatesKey    = "Import updates:"
	exportUpdatesKey    = "Export updates:"
	newline             = "\n"
	emptyString         = ""
	bgp                 = "BGP"
	birdSocketKey       = "BIRD_SOCKET"
)

func main() {
	if os.Getenv(birdSocketKey) != emptyString {
		socketPath = os.Getenv(birdSocketKey)
	}

	flag.StringVar(&socketPath, "socket", socketPath, "path to socket")
	flag.BoolVar(&debug, "debug", debug, "output more data")
	flag.BoolVar(&showVersion, "version", showVersion, "shows the version")
	flag.Parse()

	if showVersion {
		fmt.Printf("%s (%s) [%s]\n", version, commit, date)
	}

	socket := birdsocket.NewSocket(socketPath, birdsocket.WithBufferSize(bufferSize))

	_, err := socket.Connect()
	if err != nil {
		fmt.Printf("failed to connect to socket '%s'\n", socketPath)
		return
	}

	reply, err := socket.Query(showProtocolsAllCMD)
	if err != nil {
		fmt.Printf("failed to query BIRD on socket '%s'\n", socketPath)
		return
	}

	protocols := string(reply)

	if debug {
		fmt.Println(protocols)
	}

	lines := strings.Split(protocols, newline)

	isOpen := false
	isBGP := false

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorders(borders)
	table.SetHeader(header)

	session := emptyString
	state := emptyString
	neighbor := emptyString
	as := emptyString
	im := emptyString
	ex := emptyString

	for _, line := range lines {
		if strings.HasPrefix(line, newEntryKey) {
			isOpen = true

			line = strings.ReplaceAll(line, newEntryKey, emptyString)
			fields := strings.Fields(line)
			if fields[1] != bgp {
				continue
			}
			isBGP = true
			session = fields[0]

			state = getElementAtIndex(fields, 3)
		}

		if strings.Contains(line, neighborASKey) && isBGP && isOpen {
			neighborASFields := strings.Fields(line)
			as = getElementAtIndex(neighborASFields, 2)
		}

		if strings.Contains(line, neighborAddressKey) && isBGP && isOpen {
			neighborAddressFields := strings.Fields(line)
			neighbor = getElementAtIndex(neighborAddressFields, 2)
		}

		if strings.Contains(line, importUpdatesKey) && isBGP && isOpen {
			importUpdateFields := strings.Fields(line)
			im = getElementAtIndex(importUpdateFields, 6)
		}

		if strings.Contains(line, exportUpdatesKey) && isBGP && isOpen {
			exportUpdateFields := strings.Fields(line)
			ex = getElementAtIndex(exportUpdateFields, 6)
		}

		if len(line) <= 1 {
			isOpen = false
			isBGP = false

			if session != emptyString {
				table.Append([]string{session, state, neighbor, as, im, ex})
			}

			session = emptyString
			state = emptyString
			neighbor = emptyString
			as = emptyString
			im = emptyString
			ex = emptyString
		}
	}

	table.Render()
}

func getElementAtIndex(slice []string, index int) string {
	if index < 0 || index >= len(slice) {
		return ""
	}
	return slice[index]
}
