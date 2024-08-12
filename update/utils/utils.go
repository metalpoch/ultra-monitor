package utils

import (
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/metalpoch/olt-blueprint/update/model"
)

func BytesToKbps(prevBytes, currBytes, diffDate int64) int32 {
	bps := math.Abs(float64(8*currBytes)-float64(8*prevBytes)) / float64(diffDate)
	return int32(math.Round(bps / 1000))
}

func Mkdir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func GetGPONInterface(olt, ifname string) model.ElementOLT {
	s := strings.Replace(ifname, "GPON ", "", 1)
	parts := strings.Split(s, "/")
	shell, err := strconv.Atoi(parts[0])

	if err != nil {
		log.Printf("error parsing %s: %s\n", s, err.Error())
	}

	card, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Printf("error parsing %s: %s\n", s, err.Error())
	}

	port, err := strconv.Atoi(parts[2])
	if err != nil {
		log.Printf("error parsing %s: %s\n", s, err.Error())
	}

	return model.ElementOLT{
		OLT:       olt,
		Interface: ifname,
		Slot:      int8(shell),
		Card:      int8(card),
		Port:      int8(port),
	}
}
