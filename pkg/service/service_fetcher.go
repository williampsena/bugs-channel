package service

import (
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
	plugin "github.com/williampsena/bugs-channel-plugins/pkg/service"
	"github.com/williampsena/bugs-channel/pkg/settings"
)

// Represents an error when service is not found
var ErrServiceNotFound = errors.New("an error occurred when attempting to fetch the service")

type YAMLServiceFetcher struct {
	services []settings.ConfigFileService
}

func (s *YAMLServiceFetcher) GetServiceByAuthKey(authKey string) (plugin.Service, error) {
	if authKey == "" {
		return plugin.Service{}, ErrServiceNotFound
	}

	for _, s := range s.services {
		for _, a := range s.AuthKeys {
			if a.Key == authKey && !a.Disabled && !isAuthKeyExpired(a.ExpiredAt) {
				return plugin.Service{Id: s.Id, Name: s.Name}, nil
			}
		}
	}

	log.Debugf("AuthKey: %v", authKey)

	return plugin.Service{}, ErrServiceNotFound
}

func isAuthKeyExpired(expiredAt int64) bool {
	if expiredAt == 0 {
		return false
	}

	return expiredAt < time.Now().Unix()
}

// Build a new yaml service fetcher instance
func NewYAMLServiceFetcher(services []settings.ConfigFileService) plugin.ServiceFetcher {
	return &YAMLServiceFetcher{services}
}
