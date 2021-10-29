package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Process struct {
	Name string
	PID int
	SessionName string
	SessionId int
	MemUsage float64
}


func main() {
	start := time.Now()
	output, err := exec.Command("cmd", "/c", "tasklist /fo list").Output()

	if err != nil {
		fmt.Printf("Failed to get task list because: %v\n", err)
		os.Exit(-1)
	}

	strOutput := string(output)

	var processes []Process
	var prc Process

	for _, line := range strings.Split(strOutput, "\n") {
		l := strings.ReplaceAll(line, "\r", "")
		if l == "" {
			if prc.IsValid() {
				processes = append(processes, prc)
			} else {
				prc = Process{}
			}

			continue
		}

		name, found := GetImageName(l)
		if found {
			prc.Name = name
			continue
		}

		var pid int
		pid, found = GetPid(l)
		if found {
			prc.PID = pid
			continue
		}

		var sessionName string
		sessionName, found = GetSessionName(l)
		if found {
			prc.SessionName = sessionName
			continue
		}

		var sessionId int
		sessionId, found = GetSessionId(l)
		if found {
			prc.SessionId = sessionId
			continue
		}

		var memUsage float64
		memUsage, found = GetMemUsage(l)
		if found {
			prc.MemUsage = memUsage
			continue
		}
	}

	for _, prc := range processes {
		fmt.Println("------------------------------------------------------------------------------------")
		fmt.Printf("Name (PID): %s (%d)\n", prc.Name, prc.PID)
		fmt.Printf("Session (id): %s (%d)\n", prc.SessionName, prc.SessionId)
		fmt.Printf("Mem Usage: %f\n", prc.MemUsage)
		fmt.Println()
	}

	log.Printf("All done! (elapsed time: %s)\n", time.Since(start))
}

func (p *Process) IsValid() bool {
	return p.Name != "" && p.SessionName != "" && p.SessionId >= 0 && p.PID >= 0 && p.MemUsage >= 0
}

func GetImageName(s string) (string, bool) {
	return extractValue(s, "Image Name:")
}

func GetPid(s string) (int, bool) {
	return extractIntValue(s, "PID:")
}

func GetSessionName(s string) (string, bool) {
	return extractValue(s, "Session Name:")
}

func GetSessionId(s string) (int, bool) {
	return extractIntValue(s, "Session#:")
}

func GetMemUsage(s string) (float64, bool) {
	return extractFloatValue(s, "Mem Usage:")
}

func extractFloatValue(s string, name string) (float64, bool) {
	value, found := extractValue(s, name)

	if !found {
		return 0, false
	}

	result, err := strconv.ParseFloat(strings.Trim(strings.ReplaceAll(value, "K", ""), " "), 64)
	if err != nil {
		fmt.Println(err)
		return -1, false
	}
	return result, true
}

func extractIntValue(s string, name string) (int, bool) {
	value, found := extractValue(s, name)

	if !found {
		return 0, false
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		fmt.Println(err)
		return -1, false
	}
	return result, true
}

func extractValue(s string, name string) (string, bool) {
	if !strings.Contains(s, name) {
		return "", false
	}

	return strings.Trim(strings.ReplaceAll(s, name, ""), " "), true
}
