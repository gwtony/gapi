package utils

import (
	"io"
	"fmt"
	"crypto/rand"
	"errors"
	"encoding/binary"
	"net"
	"sync"
	"time"
)
const (
	// Intervals bewteen 1/1/1970 and 15/10/1582 (Julain days of 1 Jan 1970 - Julain days of 15 Oct 1582) * 100-Nanoseconds Per Day
	intervals = (2440587 - 2299160) * 86400 * 10000000
)

var (
	lastGenerated time.Time  // last generated time
	clockSequence uint16     // clock sequence for same tick
	nodeID        []byte     // node id (MAC Address)
	locker        sync.Mutex // global lock
)

// NewUUID generates a time-based UUID according to RFC 4122
func NewUUID() (string, error) {
	// Get and release a global lock
	locker.Lock()
	defer locker.Unlock()

	uuid := make([]byte, 16)

	// get timestamp
	now := time.Now().UTC()
	timestamp := uint64(now.UnixNano()/100) + intervals // get timestamp
	if !now.After(lastGenerated) {
		clockSequence++ // last generated time known, then just increment clock sequence
	} else {
		b := make([]byte, 2)
		_, err := rand.Read(b)
		if err != nil {
			return "", errors.New("Could not generate clock sequence")
		}
		clockSequence = uint16(int(b[0])<<8 | int(b[1])) // set to a random value (network byte order)
	}

	lastGenerated = now // remember the last generated time

	timeLow := uint32(timestamp & 0xffffffff)
	timeMiddle := uint16((timestamp >> 32) & 0xffff)
	timeHigh := uint16((timestamp >> 48) & 0xfff)

	// network byte order(BigEndian)
	binary.BigEndian.PutUint32(uuid[0:], timeLow)
	binary.BigEndian.PutUint16(uuid[4:], timeMiddle)
	binary.BigEndian.PutUint16(uuid[6:], timeHigh)
	binary.BigEndian.PutUint16(uuid[8:], clockSequence)

	// get node id(mac address)
	if nodeID == nil {
		interfaces, err := net.Interfaces()
		if err != nil {
			return "", errors.New("Could not get network interfaces")
		}

		for _, i := range interfaces {
			if len(i.HardwareAddr) >= 6 {
				nodeID = make([]byte, 6)
				copy(nodeID, i.HardwareAddr)
				break
			}
		}

		if nodeID == nil {
			nodeID = make([]byte, 6)
			_, err := rand.Read(nodeID)
			if err != nil {
				return "", errors.New("Could not generate node id")
			}
		}
	}

	copy(uuid[10:], nodeID)

	// set version(v1)
	uuid[6] = (uuid[6] | 0x10) & 0x1f

	// set layout(RFC4122)
	uuid[8] = (uuid[8] | 0x80) & 0x8f // Msb0=1, Msb1=0

	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

// NewUUIDv4 generates a random UUID according to RFC 4122
func NewUUIDv4() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	/* variant bits */
	uuid[8] = uuid[8]&^0xc0 | 0x80
	/* pseudo-random */
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
