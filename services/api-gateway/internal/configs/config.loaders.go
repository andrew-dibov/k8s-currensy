package configs

import (
	"os"
	"strings"
)

func loadEnv(k string, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

func loadKeys(k string, d map[string]bool) map[string]bool {
	s := os.Getenv(k)
	ks := make(map[string]bool)

	if s != "" {
		for _, key := range strings.Split(s, ",") {
			ks[strings.TrimSpace(key)] = true
		}
	} else {
		ks = d
	}

	return ks
}
