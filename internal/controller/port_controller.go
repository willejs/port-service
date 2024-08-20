package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/willejs/ports-service/internal/app"
	"github.com/willejs/ports-service/internal/domain/entity"
)

// PortController defines the controller for handling port-related operations.
type PortController struct {
	service *app.PortService
}

// NewPortController creates a new PortController.
func NewPortController(service *app.PortService) *PortController {
	return &PortController{service: service}
}

// ListAllPorts retrieves all ports and returns them.
func (c *PortController) ListAllPorts() ([]*entity.Port, error) {
	return c.service.ListAllPorts()
}

// I should probably do this in a service but for now this is fine
func (c *PortController) UpsertPortsFromFile() error {
	// Load ports data from JSON file
	// I should open the file and get the handle here instead of reading it all into memory
	data, err := ioutil.ReadFile("../../data/ports.json")
	if err != nil {
		log.Fatalf("error reading ports file: %v", err)
		return err
	}

	var ports map[string]*entity.Port
	if err := json.Unmarshal(data, &ports); err != nil {
		return err
	}

	// I think you can call decode next or similar in a for loop to avoid loading all into memory
	// however, then i will not have a key needed for a primary key for the db
	// i gave up on this and in future i would change how this works
	for key, port := range ports {
		// the data sometimes missess the code attribute, instead of any logic to santize or normalise data lets just use the key of the map.
		if port.Code == "" {
			// if port does not have a arrtibute code use key as code
			port.Code = key
		}
		if err := c.service.UpsertPort(port); err != nil {
			return err
		}
	}
	return nil
}
