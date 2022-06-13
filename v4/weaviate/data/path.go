package data

import "path"

func buildObjectsPath(id, className string) string {
	p := "/objects"
	if className != "" {
		p = path.Join(p, className)
	}
	return path.Join(p, id)
}
