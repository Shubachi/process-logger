package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/shirou/gopsutil/process"
)

const LOG_DIR = "logs"
const PERIOD_SEC = 5

func main() {
	trackedProcesses := os.Args[1:]
	log.Println(fmt.Sprintf("Tracking the following Processes: %v", trackedProcesses))

	if _, err := os.Stat(LOG_DIR); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(LOG_DIR, os.ModePerm)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}

	procs, err := process.Processes()
	if err != nil {
		log.Fatal(err)
	}

	for {
		processFound := false

		for _, proc := range procs {
			procName, err := proc.Name()
			if err != nil {
				continue
			}

			if !isTrackedProcess(procName, trackedProcesses) {
				continue
			}

			processFound = true
			file, err := os.OpenFile(fmt.Sprintf("%s/%s.csv", LOG_DIR, procName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				log.Fatal(fmt.Sprintf("Open file failed: Error: %s", err))
				os.Exit(1)
			}

			cpu, err := proc.CPUPercent()
			if err != nil {
				log.Fatal(fmt.Sprintf("Get CPU Failed: Error: %s", err))
				os.Exit(1)
			}

			memRss := "n/a"
			memPriv := "n/a"
			mem, err := proc.MemoryInfo()
			if err == nil {
				memRss = fmt.Sprintf("%d", mem.RSS)
				memPriv = fmt.Sprintf("%d", mem.VMS)
			}

			diskRead := "n/a"
			diskWrite := "n/a"
			disk, err := proc.IOCounters()
			if err == nil {
				diskRead = fmt.Sprintf("%d", disk.ReadBytes)
				diskWrite = fmt.Sprintf("%d", disk.WriteBytes)
			}

			netIn := "n/a"
			netOut := "n/a"

			network, err := proc.NetIOCounters(false)
			if err == nil {
				netIn = fmt.Sprintf("%d", network[0].BytesRecv)
				netOut = fmt.Sprintf("%d", network[0].BytesSent)
			}

			//time, cpu, mem_rss, mem_priv, disk_read, disk_write, network_in, network_out
			now := time.Now()
			outString := fmt.Sprintf("%s,%f,%s,%s,%s,%s,%s,%s\n", now.Format(time.RFC3339), cpu, memRss, memPriv, diskRead, diskWrite, netIn, netOut)
			log.Println(fmt.Sprintf("%s,%s", procName, outString))
			file.Write([]byte(outString))
			file.Close()
		}

		if !processFound {
			log.Println("None of the provided processes were found. Exiting")
			os.Exit(1)
		}

		time.Sleep(PERIOD_SEC * time.Second)
	}

}

func isTrackedProcess(procName string, trackedProcesses []string) bool {
	for _, trackedProc := range trackedProcesses {
		if procName == trackedProc {
			return true
		}
	}

	return false
}
