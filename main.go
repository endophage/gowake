package main

// This utility implements WakeOnLan Magic Packet sending.
// Packets should be broadcast over UDP with the following format:
// Byte Range    - Value
//   0 - 6       - 0xFF, this is a static header value
//   7 - 102     - 16 repetitions of MAC address of system to wake
// 103 - 106/108 - Optional.
//                 If 4 byte range, IPv5 address.
//                 If 6 byte range, Ethernet address

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app  = kingpin.New("wake", "Sends Wake On LAN magic packets")
	now  = app.Command("now", "Send packet now")
	name = now.Flag("name", "Name to store mac address as for aliased use").Short('n').String()

	mac  = now.Arg("MAC", "MAC address to send magic packet to").String()
	list = app.Command("list", "list saved MAC addresses and their aliases")
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case list.FullCommand():
		listLocations()
	case now.FullCommand():
		wake()
	default:
		log.Fatal("command not recognized")
	}
}

func listLocations() {
	locations, err := listNames()
	if err != nil {
		log.Fatal(err.Error())
	}
	for k, v := range locations {
		fmt.Println(k, ": ", v)
	}
}

func wake() {
	var (
		err    error
		packet []byte
	)
	if *mac != "" {
		packet, err = parseMAC(*mac)
	} else if *name != "" {
		packet, err = lookupName(*name)
	}
	if err != nil {
		log.Fatal(err.Error())
	}

	err = magicBroadcast(packet)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("System should be awake")

	// save the name and mac if both were provided. Possibly overwrites an
	// existing value
	if *name != "" && *mac != "" {
		err = saveName(*name, *mac)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}
