package snmp

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/metalpoch/olt-blueprint/update/model"
)

const (
	sysname_oid   = "1.3.6.1.2.1.1.5.0"
	ifname_oid    = "1.3.6.1.2.1.31.1.1.1.1"
	bytes_in_oid  = "1.3.6.1.4.1.2011.6.128.1.1.4.21.1.15"
	bytes_out_oid = "1.3.6.1.4.1.2011.6.128.1.1.4.21.1.30"
	bandwidth_oid = "1.3.6.1.2.1.31.1.1.1.15"
)

func runCmd(command string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stdout
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func Sysname(ip, community string) (string, error) {
	var sysname string
	cmd := fmt.Sprintf("snmpwalk -v 2c -c %s %s %s", community, ip, sysname_oid)

	stdout, stderr, err := runCmd(cmd)
	if err != nil {
		log.Printf("snmp error on device: %s - %s | stderr: %s | err: %s\n", ip, community, stderr, err.Error())
		return sysname, err
	}

	rows := strings.Split(string(stdout), "\n")
	for _, row := range rows {
		if len(row) < 1 {
			break
		}
		sysname = strings.Split(row, "STRING: ")[1]
	}
	return sysname, nil
}

func Measurements(device *model.Device) model.Snmp {
	oids := [4]string{ifname_oid, bytes_in_oid, bytes_out_oid, bandwidth_oid}
	ifnames := make(map[int]string)
	byteInOcts := make(map[int]int64)
	byteOutOcts := make(map[int]int64)
	bandwidths := make(map[int]int16)

	var wg sync.WaitGroup
	wg.Add(4)

	for _, oid := range oids {
		defer wg.Done()
		cmd := fmt.Sprintf("snmpwalk -v 2c -c %s %s %s", device.Community, device.IP, oid)
		stdout, stderr, err := runCmd(cmd)
		if err != nil {
			log.Fatalf("snmp error on device: %s - %s | stderr: %s | err: %s", device.IP, device.Community, stderr, err)
		}

		rows := strings.Split(string(stdout), "\n")
		for _, r := range rows {
			if len(r) < 1 {
				break
			}

			var val string
			splited := strings.Split(r, " = ")
			idxStr := strings.Split(splited[0], ".")
			val = strings.Split(splited[1], ": ")[1]

			idx, err := strconv.Atoi(idxStr[len(idxStr)-1])
			if err != nil {
				log.Println(device.Sysname, "-", err, string(stdout))
			}

			if oid == ifname_oid {
				ifnames[idx] = val
			} else {
				v, err := strconv.Atoi(val)
				if err != nil {
					log.Println(device.Sysname, "-", err, string(stdout))
				}

				switch oid {
				case bytes_in_oid:
					byteInOcts[idx] = int64(v)
				case bytes_out_oid:
					byteOutOcts[idx] = int64(v)
				case bandwidth_oid:
					bandwidths[idx] = int16(v)
				}
			}
		}
	}
	wg.Wait()

	return model.Snmp{
		IfName:    ifnames,
		ByteIn:    byteInOcts,
		ByteOut:   byteOutOcts,
		Bandwidth: bandwidths,
	}
}
