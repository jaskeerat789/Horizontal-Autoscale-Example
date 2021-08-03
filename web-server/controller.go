package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/rs/xid"
)

type Controller struct {
	l  hclog.Logger
	rc *RabbitMQClient
}

func NewController() *Controller {
	rc := NewClient()
	log := hclog.New(&hclog.LoggerOptions{
		Name: "Handler",
	})

	return &Controller{l: log, rc: rc}
}

func (c *Controller) GenerateOrder(rw http.ResponseWriter, r *http.Request) {
	c.l.Info("generete order")
	id := xid.New().String()
	c.rc.SendMessage([]byte(id))
	fmt.Fprintln(rw, "Welcome!", id)
}

func (c *Controller) GetStatus(rw http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	c.l.Info("get status", "ID", id)

	fmt.Fprintln(rw, "Welcome! ", id)

}
