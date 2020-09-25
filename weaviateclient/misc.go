package weaviateclient

import (
	"net/http"
)

type Misc struct {
	config *Config
}

func (misc *Misc) ReadyChecker() *readyChecker {
	return &readyChecker{config: misc.config}
}

func (misc *Misc) LifeChecker() *liveChecker {
	return &liveChecker{config: misc.config}
}

func (misc *Misc) OpenIDConfigurationGetter() *openIDConfigGetter {
	return &openIDConfigGetter{config: misc.config}
}


type readyChecker struct {
	config *Config
}

func (rc *readyChecker) Do() (bool, error) {
	return checker(rc.config, "ready")
}


type liveChecker struct {
	config *Config
}

func (lc *liveChecker) Do() (bool, error) {
	return checker(lc.config, "live")
}


func checker(config *Config, endpoint string) (bool, error) {
	path := config.basePath()+"/.well-known/"+endpoint
	response, getErr := http.Get(path)
	if getErr != nil {
		return false, getErr
	}
	defer response.Body.Close()
	return response.StatusCode == 200, nil
}


type openIDConfigGetter struct {
	config *Config
}

func (oidcg *openIDConfigGetter) Do() (map[string]string, error) {
	path := oidcg.config.basePath()+"/.well-known/openid-configuration"
	response, getErr := http.Get(path)
	if getErr != nil {
		return nil, getErr
	}
	defer response.Body.Close()
	if response.StatusCode == 200 {

		return nil, nil
		// TODO set response body to struct
	}
	// TODO return error
	return nil, nil
}