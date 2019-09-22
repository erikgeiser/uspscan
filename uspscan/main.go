package uspscan

import (
	"fmt"
)

// Run performs the unquoted service path check
func Run() {
	serviceCfgs, err := listServices()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, cfg := range serviceCfgs {
		// ignore weird services without path name
		if cfg.PathName == "" {
			continue
		}

		paths, err := searchPaths(cfg.PathName)
		if err != nil {
			fmt.Printf("couldn't build search paths: %v\n", err)
			continue
		}

		for _, p := range paths {
			if fileExists(p) {
				break
			}
			if fileCreationPossible(p) {
				fmt.Printf("[%s|%s] Vulnerable Service: %s (%s)", cfg.StartName, cfg.StartMode, p, cfg.PathName)
				break
			}
		}
	}
}
