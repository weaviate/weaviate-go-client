package util

import "sync"

type GetVersionFn = func() string

type DBVersionProvider struct {
	mutex        sync.RWMutex
	getVersionFn GetVersionFn
	version      string
}

func NewDBVersionProvider(getVersionFn GetVersionFn) *DBVersionProvider {
	return &DBVersionProvider{getVersionFn: getVersionFn, version: getVersionFn()}
}

func (s *DBVersionProvider) Refresh() {
	s.refreshIfEmpty(false)
}

func (s *DBVersionProvider) ForceRefresh() {
	s.refreshIfEmpty(true)
}

func (s *DBVersionProvider) Version() string {
	s.refreshIfEmpty(false)
	return s.getVersion()
}

func (s *DBVersionProvider) getVersion() string {
	return s.version
}

func (s *DBVersionProvider) updateVersion() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.version = s.getVersionFn()
}

func (s *DBVersionProvider) isEmptyVersion() bool {
	return len(s.getVersion()) == 0
}

func (s *DBVersionProvider) refreshIfEmpty(force bool) {
	if force || s.isEmptyVersion() {
		s.updateVersion()
	}
}
