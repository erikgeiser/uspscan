package uspscan

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type serviceCfg struct {
	DisplayName string
	PathName    string
	StartMode   string
	StartName   string
}

func runWMICQuery() (string, error) {
	wmic, err := exec.LookPath("wmic")
	if err != nil {
		return "", fmt.Errorf("can't find wmic: %v", err)
	}
	cmd := exec.Command(wmic, "service", "get", "DisplayName,PathName,StartMode,StartName", "/format:csv")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("wmic call failed: %v", err)
	}
	serviceCSV := strings.TrimSpace(string(output))
	serviceCSV = regexp.MustCompile(`[\t\r\n]+`).ReplaceAllString(serviceCSV, "\n")
	return serviceCSV, nil
}

func listServices() ([]serviceCfg, error) {
	serviceCSV, err := runWMICQuery()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(strings.NewReader(serviceCSV))
	ok := scanner.Scan()
	if !ok {
		return nil, fmt.Errorf("can't read wmic csv output header: %v", err)
	}
	header, _ := csvSplit(scanner.Text(), -1)

	configs := []serviceCfg{}

	for scanner.Scan() {
		record, err := csvSplit(scanner.Text(), len(header))
		if err != nil {
			return nil, fmt.Errorf("can't read record: %v", err)
		}
		cfg, err := parseConfig(header, record)
		if err != nil {
			return nil, fmt.Errorf("can't parse record: %v", err)
		}
		configs = append(configs, cfg)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	return configs, nil
}

func parseConfig(csvHeader, csvRecord []string) (serviceCfg, error) {
	mustGetValue := func(name string) string {
		idx := indexOf(csvHeader, name)
		if idx < 0 {
			panic(fmt.Errorf("%s not in csv header %v", name, csvHeader))
		}
		if idx >= len(csvRecord) {
			panic(fmt.Errorf("%s header index larger than record %v", name, csvRecord))
		}
		return csvRecord[idx]
	}
	return serviceCfg{
		DisplayName: mustGetValue("DisplayName"),
		PathName:    mustGetValue("PathName"),
		StartMode:   mustGetValue("StartMode"),
		StartName:   mustGetValue("StartName"),
	}, nil
}

func indexOf(data []string, element string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1
}

func csvSplit(line string, wantedLength int) ([]string, error) {
	parts := strings.Split(line, ",")
	if wantedLength >= 0 && len(parts) != wantedLength {
		return nil, fmt.Errorf("read %d values instead of %d", len(parts), wantedLength)
	}
	return parts, nil
}
