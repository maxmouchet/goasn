package peeringdb

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type IX struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type LAN struct {
	ID   int `json:"id"`
	IXID int `json:"ix_id"`
}

type Prefix struct {
	ID      int    `json:"id"`
	IXLanID int    `json:"ixlan_id"`
	Prefix  string `json:"prefix"`
}

type DB struct {
	IXs      []IX
	LANs     []LAN
	Prefixes []Prefix
}

func (db *DB) FromAPI() error {
	var ixResponse struct {
		Data []IX `json:"data"`
	}

	var lanResponse struct {
		Data []LAN `json:"data"`
	}

	var pfxResponse struct {
		Data []Prefix `json:"data"`
	}

	err := getJSON("https://peeringdb.com/api/ix.json", &ixResponse)
	if err != nil {
		return err
	}

	err = getJSON("https://peeringdb.com/api/ixlan.json", &lanResponse)
	if err != nil {
		return err
	}

	err = getJSON("https://peeringdb.com/api/ixpfx.json", &pfxResponse)
	if err != nil {
		return err
	}

	db.IXs = ixResponse.Data
	db.LANs = lanResponse.Data
	db.Prefixes = pfxResponse.Data

	return nil
}

func getJSON(url_ string, v interface{}) error {
	client := http.Client{
		Timeout: time.Second * 30,
	}

	req, err := http.NewRequest("GET", url_, nil)
	if err != nil {
		return err
	}

	log.Printf("GET %s", req.URL)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	return json.NewDecoder(resp.Body).Decode(&v)
}
