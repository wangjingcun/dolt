package earl

import (
	"github.com/liquidata-inc/ld/dolt/go/libraries/utils/osutil"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var validHostRegex = regexp.MustCompile("^[-.a-zA-z0-9]*$")
var validHostWithPortRegex = regexp.MustCompile("^[-.a-zA-z0-9]*:[0-9]*$")

func isValidHost(hostAndPortStr string) bool {
	hostStr := hostAndPortStr
	portStr := ""

	if idx := strings.IndexRune(hostAndPortStr, ':'); idx != -1 {
		hostStr = hostAndPortStr[:idx]
		portStr = strings.TrimSpace(hostAndPortStr[idx+1:])
	}

	if len(portStr) > 0 {
		if _, err := strconv.ParseUint(portStr, 10, 16); err != nil {
			return false
		}
	}

	if hostStr == "" {
		return false
	} else if hostStr == "localhost" {
		return true
	} else if strings.Index(hostStr, ".") == -1 {
		return false
	}

	return validHostRegex.MatchString(hostStr) || validHostWithPortRegex.MatchString(hostStr)
}

func Parse(urlStr string) (*url.URL, error) {
	u, err := parse(urlStr)

	if err != nil {
		return nil, err
	}

	// if Path is e.g. "/C$/" for a network location, it should instead be "C:/"
	if len(u.Path) >= 3 && u.Path[0] == '/' && u.Path[1] >= 'A' && u.Path[1] <= 'Z' && u.Path[2] == '$' {
		u.Path = u.Path[1:2] + ":" + u.Path[3:]
	} else if !osutil.StartsWithWindowsVolume(u.Path) { // normalize some
		if len(u.Path) == 0 || u.Path[0] != '/' {
			u.Path = "/" + u.Path
		}
	}
	u.Path = strings.ReplaceAll(u.Path, `\`, "/")

	return u, nil
}

func parse(urlStr string) (*url.URL, error) {
	if strIdx := strings.Index(urlStr, ":///"); strIdx != -1 && osutil.StartsWithWindowsVolume(urlStr[strIdx+4:]) {
		return &url.URL{
			Scheme: urlStr[:strIdx],
			Path: urlStr[strIdx+4:],
		}, nil
	}
	if strings.Index(urlStr, "://") == -1 {
		u, err := url.Parse("http://" + urlStr)

		if err == nil && isValidHost(u.Host) {
			u.Scheme = ""
			return u, nil
		} else if err != nil {
			return nil, err
		}
	}

	return url.Parse(urlStr)
}
