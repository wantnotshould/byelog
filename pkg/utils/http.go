// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package utils

import (
	"net"
	"net/http"
)

func IP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
