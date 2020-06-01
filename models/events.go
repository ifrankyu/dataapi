package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/astaxie/beego"
	"github.com/gomodule/redigo/redis"
)

// EventList.
var (
	EventList map[string]*Event
)

const (
	// DYMFEVENTSREPORT represents fatal errors
	DYMFEVENTSREPORT = "dymf_events_report"
)

// addEventToRedis func.
func addEventToRedis(e *Event) bool {
	client, err := redis.DialURL("redis://" + beego.AppConfig.String("redishost") + ":" + beego.AppConfig.String("redisport"))

	if err != nil {
		log.Fatalln(err)
	}

	defer client.Close()

	rel, _ := json.Marshal(&e)
	fmt.Println(string(rel))

	v, err := client.Do("lpush", DYMFEVENTSREPORT, string(rel))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(v)
	return true
}

// getEventFromRedis func.
func getEventFromRedis(id string) (e *Event, err error) {
	client, err := redis.DialURL("redis://" + beego.AppConfig.String("redishost") + ":" + beego.AppConfig.String("redisport"))

	if err != nil {
		log.Fatalln(err)
	}

	defer client.Close()
	values, err := redis.Values(client.Do("lrange", DYMFEVENTSREPORT, "0", "1000"))
	if err != nil {
		return e, err
	}

	var str string
	for _, v := range values {
		str = string(v.([]byte))
		if id == str {
			json.Unmarshal(v.([]byte), &e)
			break
		}
		fmt.Println(string(v.([]byte)))
	}

	return e, err
}

// GetAllEventsFromRedis func.
func GetAllEventsFromRedis() map[string]*Event {
	client, err := redis.DialURL("redis://" + beego.AppConfig.String("redishost") + ":" + beego.AppConfig.String("redisport"))
	if err != nil {
		log.Fatalln(err)
	}

	defer client.Close()

	values, err := redis.Values(client.Do("lrange", DYMFEVENTSREPORT, "0", "1000"))
	if err != nil {
		return EventList
	}

	var e Event
	for _, v := range values {
		fmt.Println(string(v.([]byte)))
		err = json.Unmarshal(v.([]byte), &e)
		if err != nil {
			fmt.Println("json Unmarshal error:", err)
		}
		EventList[e.PlatformID] = &e
	}

	return EventList
}

func init() {
	EventList = make(map[string]*Event)
	// event := Event{"dymf", "iostest", 1587109493067, "1.2.23", "768", "dbc218edc", "data list"}

	// EventList["dbc218edc"] = &event

	// addEventToRedis(&event)
}

// Event struct.
type Event struct {
	Game       string
	Platform   string
	Time       int64
	Version    string
	UID        string
	PlatformID string
	Data       string
}

// AddEvent func.
func AddEvent(e Event, id string) string {
	e.PlatformID = id
	EventList[id] = &e
	addEventToRedis(&e)
	return id
}

// GetEvent func.
func GetEvent(id string) (e *Event, err error) {
	if e, ok := EventList[id]; ok {
		return e, nil
	}

	e, err = getEventFromRedis(id)
	if err != nil {
		return e, nil
	}

	return nil, errors.New("Event not exists")
}

// GetAllEvents func.
func GetAllEvents() map[string]*Event {
	if len(EventList) > 0 {
		fmt.Println("EventList data from memory")
		return EventList
	}

	return GetAllEventsFromRedis()
}

// UpdateEvent func.
func UpdateEvent(id string, ee *Event) (e *Event, err error) {
	if e, ok := EventList[id]; ok {
		if ee.Game != "" {
			e.Game = ee.Game
		}
		if ee.Platform != "" {
			e.Platform = ee.Platform
		}
		if ee.Time != 0 {
			e.Time = ee.Time
		}
		if ee.Version != "" {
			e.Version = ee.Version
		}
		if ee.UID != "" {
			e.UID = ee.UID
		}
		if ee.PlatformID != "" {
			e.PlatformID = ee.PlatformID
		}
		if ee.Data != "" {
			e.Data = ee.Data
		}
		return e, nil
	}
	return nil, errors.New("Event Not Exist")
}

// DeleteEvent func.
func DeleteEvent(id string) {
	delete(EventList, id)
}
