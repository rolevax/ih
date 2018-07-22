package mako

import "fmt"

func keyAuth(username string) string {
	return fmt.Sprintf("mako.auth:%v", username)
}
