package utils

import (
	"encoding/xml"
	"github.com/darkliquid/leader1/config"
)

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
		logger.Printf("Couldn't decode shoutcast server stats XML: %s", err.Error())
		return ShoutcastServerStats{}, err
	}
	return
}

// Does the full process for returning a populated shoutcast stat object
func GetShoutcastStats(cfg *config.Settings) (ShoutcastServerStats, error) {
	var res string
	var err error

	if cfg.Stream.StatsUser != "" && cfg.Stream.StatsPass != "" {
		res, err = GetPageWithAuth(cfg.Stream.StatsURL, cfg.Stream.StatsUser, cfg.Stream.StatsPass)
	} else {
		res, err = GetPage(cfg.Stream.StatsURL)
	}
	if err != nil {
		logger.Printf("Couldn't load page - %s", err.Error())
		return ShoutcastServerStats{}, err
	}

	stats, err := shoutcastStats(FixInvalidUTF8(res))
	if err != nil {
		logger.Printf("Couldn't parse stats - %s", err.Error())
		return ShoutcastServerStats{}, err
	}

	return stats, err
}