package controllers

import (
	"beeapi/models"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/astaxie/beego"
)

// EventsController Operations about Events.
type EventsController struct {
	beego.Controller
}

// GetRandomString func
func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// Post func.
// @Title CreateEvent
// @Description create events
// @Param	body		body 	models.Event	true		"body for event content"
// @Success 200 {int} models.Event.Id
// @Failure 403 body is empty
// @router / [post]
func (e *EventsController) Post() {
	var event models.Event
	err := json.Unmarshal(e.Ctx.Input.RequestBody, &event)
	if err != nil {
		fmt.Println("err", err)
	}
	var id string
	e.Ctx.Input.Bind(&id, "id")
	if id == "" {
		id = GetRandomString(10)
	}

	id = models.AddEvent(event, id)
	fmt.Println("id:", id)
	fmt.Println("event", event)

	e.Data["json"] = map[string]string{"id": id}
	e.ServeJSON()
}

// GetAll func.
// @Title GetAll
// @Description get all Events
// @Success 200 {object} models.Event
// @router / [get]
func (e *EventsController) GetAll() {
	events := models.GetAllEvents()
	e.Data["json"] = events
	e.ServeJSON()
}

// Get func.
// @Title Get
// @Description get event by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Event
// @Failure 403 :id is empty
// @router /:id [get]
func (e *EventsController) Get() {
	id := e.GetString(":id")

	if id != "" {
		event, err := models.GetEvent(id)
		if err != nil {
			e.Data["json"] = err.Error()
		} else {
			e.Data["json"] = event
		}
	}
	e.ServeJSON()
}

// Put func.
// @Title Update
// @Description update the event
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Event	true		"body for event content"
// @Success 200 {object} models.Event
// @Failure 403 :id is not int
// @router /:id [put]
func (e *EventsController) Put() {
	uid := e.GetString(":id")
	if uid != "" {
		var event models.Event
		json.Unmarshal(e.Ctx.Input.RequestBody, &event)
		uu, err := models.UpdateEvent(uid, &event)
		if err != nil {
			e.Data["json"] = err.Error()
		} else {
			e.Data["json"] = uu
		}
	}
	e.ServeJSON()
}

// Delete func.
// @Title Delete
// @Description delete the event
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (e *EventsController) Delete() {
	id := e.GetString(":id")
	models.DeleteEvent(id)
	e.Data["json"] = "delete success!"
	e.ServeJSON()
}
