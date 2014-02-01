package commands

import (
	"encoding/xml"
	"fmt"
	"github.com/fluffle/golog/logging"
	"io/ioutil"
	"net/http"
)

func getPage(url string) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logging.Warn("Couldn't build http request")
		return "", err
	}

	req.Header.Set("User-Agent", "Leader-1/Mighty, Mighty GoBot")

	resp, err := client.Do(req)
	if err != nil {
		logging.Warn("Couldn't perform http request")
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logging.Warn("Couldn't read http response body")
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
		logging.Warn("Couldn't decode shoutcast server stats XML")
		return ShoutcastServerStats{}, err
	}
	return
}

// Does the full process for returning a populated shoutcast stat object
func getShoutcastStats() (ShoutcastServerStats, error) {
	res, err := getPage("http://streamhere.g3-radio.net:4122/stats?sid=1")
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
