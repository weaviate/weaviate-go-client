package db

import "sync"

type GetVersionFn = func() string

type VersionProvider struct {
	mutex        sync.RWMutex
	getVersionFn GetVersionFn
	version      string
}

func NewVersionProvider(getVersionFn GetVersionFn) *VersionProvider {
	return &VersionProvider{getVersionFn: getVersionFn, version: getVersionFn()}
}

func (v *VersionProvider) Refresh() {
	v.refreshIfEmpty(false)
}

func (v *VersionProvider) ForceRefresh() {
	v.refreshIfEmpty(true)
}

func (v *VersionProvider) Version() string {
	v.refreshIfEmpty(false)
	return v.getVersion()
}

func (v *VersionProvider) getVersion() string {
	return v.version
}

func (v *VersionProvider) updateVersion() {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	v.version = v.getVersionFn()
}

func (v *VersionProvider) isEmptyVersion() bool {
	return len(v.getVersion()) == 0
}

func (v *VersionProvider) refreshIfEmpty(force bool) {
	if force || v.isEmptyVersion() {
		v.updateVersion()
	}
}
