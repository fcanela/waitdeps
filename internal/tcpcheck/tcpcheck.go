package tcpcheck

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"
)

var hostPortRegex = regexp.MustCompile(`^([a-zA-Z0-9\-\_\.])+:(\d+)`)

func parseAddress(uri_string string) (string, error) {
	isHostPortFormat := hostPortRegex.MatchString(uri_string)
	if isHostPortFormat {
		uri_string = "//" + uri_string
	}

	uri, err := url.Parse(uri_string)
	if err != nil {
		return "", err
	}

	host := uri.Hostname()
	port := uri.Port()

	if port == "" {
		resolvedPort, err := net.LookupPort("tcp", uri.Scheme)
		if err != nil {
			return "", err
		}
		port = strconv.Itoa(resolvedPort)
	}

	if host == "" || port == "" {
		return "", fmt.Errorf("ensure it is HOST:PORT or URI format")
	}

	return host + ":" + port, nil
}

func checkAddress(address string) bool {
	conn, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func checkEntry(ctx context.Context, wg *sync.WaitGroup, uri string) {
	defer wg.Done()

	address, err := parseAddress(uri)
	if err != nil {
		fmt.Printf("unable to parse address %s. %v\n", uri, err)
		os.Exit(1)
	}

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("%s is unavailable\n", address)
			return
		default:
			isAvailable := checkAddress(address)
			if isAvailable {
				fmt.Printf("successful tcp connection check for %s\n", address)
				return
			} else {
				fmt.Printf("tcp connection check failed for %s. Retrying...\n", address)
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func Check(ctx context.Context, wg *sync.WaitGroup, uris []string) {
	for _, uri := range uris {
		wg.Add(1)
		go checkEntry(ctx, wg, uri)
	}
}
