package util

type GetVersionFn = func() string

type DBVersionProvider struct {
	getVersionFn GetVersionFn
	version      string
}

func NewDBVersionProvider(getVersionFn GetVersionFn) *DBVersionProvider {
	return &DBVersionProvider{getVersionFn: getVersionFn, version: getVersionFn()}
}

func (s *DBVersionProvider) Refresh() {
	if len(s.version) == 0 {
		s.version = s.getVersionFn()
	}
}

func (s *DBVersionProvider) Version() string {
	return s.version
}
