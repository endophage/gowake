package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"strings"
)

const (
	port = 80 // ports 0 or 7 would be typical but 0 fails and 7 doesn't seem to work
)

var (
	broadcastAddr = net.UDPAddr{
		IP:   net.IPv4bcast,
		Port: port,
	}
	basePacket = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF} // sync header
)

func parseMAC(m string) ([]byte, error) {
	log.Println("Waking", m)
	var err error
	res := basePacket
	macHex := strings.Split(m, ":")
	if len(macHex) != 6 {
		return nil, fmt.Errorf("invalid MAC address. Must be in format 00:00:00:00:00:00, received '%s'", m)
	}
	macBytes := make([][]byte, len(macHex))
	for i, seg := range macHex {
		macBytes[i], err = hex.DecodeString(seg)
		if err != nil {
			return nil, err
		}
	}
	mac := bytes.Join(macBytes, nil)
	// append 16 iterations of MAC address
	for i := 0; i < 16; i++ {
		res = append(res, mac...)
	}
	return res, nil
}

func magicBroadcast(packet []byte) error {
	sock, err := net.DialUDP("udp", nil, &broadcastAddr)
	if err != nil {
		return err
	}
	n, err := sock.Write(packet)
	if err != nil {
		return err
	}
	if n < len(packet) {
		return fmt.Errorf("failed to send complete magic packet")
	}
	return nil
}
