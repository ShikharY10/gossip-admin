package handler

import (
	"gbADMIN/config"
	"gbADMIN/schema"
)

type Handler struct {
	Database *DataBaseHandler
	Cache    *CacheHandler
	Env      *config.ENV
	Services []schema.Service
}

func (handle *Handler) RemoveService(serviceName string) {
	var index int = -1
	for i, service := range handle.Services {
		if service.Name == serviceName {
			index = i
			break
		}
	}
	if index != -1 {
		handle.Services = append(handle.Services[:index], handle.Services[index+1:]...)
	}
}
