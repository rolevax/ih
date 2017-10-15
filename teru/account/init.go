package account

import "os"

var jwtIssuer = ""

func init() {
	hostname, err := os.Hostname()
	if err == nil {
		jwtIssuer = "teru@" + hostname
	}
}
