package commands

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/fluffle/golog/logging"
	"github.com/darkliquid/leader1/config"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

func getPage(url string) (string, error) {
	client := newHttpTimeoutClient()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logging.Warn(fmt.Sprintf("Couldn't build http request: %s", err.Error()))
		return "", err
	}

	req.Header.Set("User-Agent", "Leader-1/Mighty, Mighty GoBot")

	resp, err := client.Do(req)
	if err != nil {
		logging.Warn(fmt.Sprintf("Couldn't perform http request: %s", err.Error()))
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logging.Warn(fmt.Sprintf("Couldn't read http response body: %s", err.Error()))
		return "", err
	}

	return string(body), err
}

type ShoutcastServerStats struct {
	CurrentListeners int    `xml:"CURRENTLISTENERS"`
	PeakListeners    int    `xml:"PEAKLISTENERS"`
	MaxListeners     int    `xml:"MAXLISTENERS"`
	UniqueListeners  int    `xml:"UNIQUELISTENERS"`
	AverageTime      int    `xml:"AVERAGETIME"`
	ServerGenre      string `xml:"SERVERGENRE"`
	ServerUrl        string `xml:"SERVERURL"`
	ServerTitle      string `xml:"SERVERTITLE"`
	SongTitle        string `xml:"SONGTITLE"`
	StreamHits       int    `xml:"STREAMHITS"`
	StreamStatus     int    `xml:"STREAMSTATUS"`
	BackupStatus     int    `xml:"BACKUPSTATUS"`
	StreamPath       string `xml:"STREAMPATH"`
	StreamUptime     int    `xml:"STREAMUPTIME"`
	BitRate          int    `xml:"BITRATE"`
	Content          string `xml:"CONTENT"`
	Version          string `xml:"VERSION"`
}

// Returns a shoutcast stats object
func shoutcastStats(stats string) (xml_stats ShoutcastServerStats, err error) {
	if err = xml.Unmarshal([]byte(stats), &xml_stats); err != nil {
		logging.Warn(fmt.Sprintf("Couldn't decode shoutcast server stats XML: %s", err.Error()))
		return ShoutcastServerStats{}, err
	}
	return
}

// Does the full process for returning a populated shoutcast stat object
func getShoutcastStats() (ShoutcastServerStats, error) {
	res, err := getPage(config.Config.Stream.StatsURL)
	if err != nil {
		logging.Error(fmt.Sprintf("Couldn't load page - %s", err.Error()))
		return ShoutcastServerStats{}, err
	}

	stats, err := shoutcastStats(res)
	if err != nil {
		logging.Error(fmt.Sprintf("Couldn't parse stats - %s", err.Error()))
		return ShoutcastServerStats{}, err
	}

	return stats, err
}

func timeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, cTimeout)
		if err != nil {
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(rwTimeout))
		return conn, nil
	}
}

// Sets up a 2 second timeout http client because waiting longer for a page request is too costly
func newHttpTimeoutClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Dial: timeoutDialer(time.Second*2, time.Second*2),
		},
	}
}

// Extracts an URL out of a string
func ExtractURL(str string) (url string, err error) {
	// Assume errors by default
	err = errors.New("No URL found")

	// Capture links and post their titles, etc
	start := strings.Index(str, "http://")

	// Try https if no http match
	if start == -1 {
		start = strings.Index(str, "https://")
	}

	// Found a link... maybe
	if start > -1 {
		url = strings.SplitN(str[start:], " ", 2)[0]
		// String isn't just a protocol
		if len(url) > 9 && !strings.HasSuffix(url, "://") {
			err = nil
		} else {
			url = ""
		}
	}
	return
}
