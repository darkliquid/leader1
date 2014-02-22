package utils

import (
	"errors"
	"github.com/darkliquid/leader1/database"
	"strings"
	"time"
)

func LikeTrack(nick, track string) (bool, error) {
	db, err := database.DB()
	if err != nil {
		return false, err
	}

	if strings.TrimSpace(nick) == "" {
		return false, errors.New("Empty nick")
	}

	if strings.TrimSpace(track) == "" {
		return false, errors.New("Empty track")
	}

	_, err = db.Exec("INSERT INTO likelogs (type, user, song, date) VALUES ('like', ?, ?, ?)", nick, track, time.Now().Unix())

	if err != nil {
		logger.Printf("DB ERROR: %s", err.Error())
		return false, err
	}

	return true, nil
}

func HateTrack(nick, track string) (bool, error) {
	db, err := database.DB()
	if err != nil {
		return false, err
	}

	if strings.TrimSpace(nick) == "" {
		return false, errors.New("Empty nick")
	}

	if strings.TrimSpace(track) == "" {
		return false, errors.New("Empty track")
	}

	_, err = db.Exec("INSERT INTO likelogs (type, user, song, date) VALUES ('dislike', ?, ?, ?)", nick, track, time.Now().Unix())

	if err != nil {
		logger.Printf("DB ERROR: %s", err.Error())
		return false, err
	}

	return true, nil
}

func Request(nick, request string) (bool, error) {
	db, err := database.DB()
	if err != nil {
		return false, err
	}

	if strings.TrimSpace(nick) == "" {
		return false, errors.New("Empty nick")
	}

	if strings.TrimSpace(request) == "" {
		return false, errors.New("Empty request")
	}

	_, err = db.Exec("INSERT INTO requests (user, song, date) VALUES (?, ?, ?)", nick, request, time.Now().Unix())

	if err != nil {
		logger.Printf("DB ERROR: %s", err.Error())
		return false, err
	}

	return true, nil
}
