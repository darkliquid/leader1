package commands

import (
	"fmt"
	irc "github.com/darkliquid/goirc/client"
	"github.com/fluffle/golog/logging"
	"html"
	"net"
	"strings"
)

func URLTitler(conn *irc.Conn, line *irc.Line, target string, url string) {
	// Don't expand our own URLs
	if line.Nick == conn.Me.Nick {
		return
	}

	// use custom client from utils.go because otherwise risk infinite-waits/hangs
	client := newHttpTimeoutClient()
	resp, err := client.Get(url)

	if err != nil {
		if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
			conn.Privmsg(target, fmt.Sprintf("[Link] request timed out"))
		}
		logging.Warn(fmt.Sprintf("Failed to GET %s due to %s", url, err.Error()))
		return
	}

	// Make sure we close our response reader like a good citizen
	defer resp.Body.Close()

	var title string

	// No point in parsing response if it wasn't a success
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// Read 10kb of the body
		body := make([]byte, 0)
		var count int

		for {
			more_body := make([]byte, 10240) // Read in up to 10kb increments
			count, err = resp.Body.Read(more_body)
			body = append(body, more_body[:count]...)
			if len(body) >= 10240 { // Got 10kb, so break out
				break
			}

			// Couldn't read the page for some reason, maybe EOF
			if err != nil {
				if len(body) == 0 {
					fmt.Printf("Failed to read page response for %s due to %s", url, err.Error())
					return
				}
				break
			}
		}

		pageData := string(body)

		start := strings.Index(strings.ToLower(pageData), "<title")
		end := strings.Index(strings.ToLower(pageData), "</title>")

		switch {
		case start > -1 && end > -1: // Found a title tag, get it's contents
			title = pageData[start:end]
			title = pageData[start+strings.Index(title, ">")+1 : end]
		case start > -1: // If for some reason within 10kb we get a chopped off <title>, use the remainder
			title = pageData[start:]
			title = pageData[start+strings.Index(title, ">")+1:]
		default: // No title tags to speak off, use the mime type instead to be helpful
			title = fmt.Sprintf("unknown [%s]", resp.Header.Get("Content-Type"))
		}
	} else {
		title = resp.Status // Lets just return the status text for non-successful responses
	}

	// Post back to channel
	conn.Privmsg(target, fmt.Sprintf("[Link] %s", html.UnescapeString(title)))
}
