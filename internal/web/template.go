package web

import (
	"bytes"

	"github.com/MindHunter86/eyesonly/internal/utils"
)

// This code uses a lot of byte allocations,
// but it will be called only if new Host header detected
type Template struct {
	BasePrefix, BasePostfix []byte
}
type Templates map[string]*Template

func (m Templates) GetPrefix(host string) []byte {
	dynamicMu.RLock()
	defer dynamicMu.RUnlock()

	if t, ok := m[host]; ok {
		return t.BasePrefix
	}
	return nil
}

func (m Templates) GetPostfix(host string) []byte {
	dynamicMu.RLock()
	defer dynamicMu.RUnlock()

	if t, ok := m[host]; ok {
		return t.BasePostfix
	}
	return nil
}

func (m Templates) NewTemplate(host string, tpl *Template) {
	dynamicMu.Lock()
	defer dynamicMu.Unlock()

	k := utils.CopyString(host)
	m[k] = tpl
}

func (m Templates) GenerateIndex(baseHost string, cfg *Config) ([]byte, []byte) {
	var e error
	var idxbuf []byte
	if idxbuf, e = openFile("dist/pages/index.html"); e != nil {
		panic(utils.ExtraErrorWrapper(e, "try to open index.html file"))
	}

	var ok bool
	var befhead, afthead []byte
	if befhead, afthead, ok = bytes.Cut(idxbuf, []byte("<qtpl-head></qtpl-head>")); !ok {
		panic("index.html <qtpl-head> template delimiter not found")
	}

	var dynhead []byte
	// dynhead = []byte(RenderHead(baseHost, cfg))

	tplbuf := make([]byte, 0, len(idxbuf)+len(dynhead))
	tplbuf = append(tplbuf, befhead...)
	tplbuf = append(tplbuf, dynhead...)
	tplbuf = append(tplbuf, afthead...)

	tpl := &Template{}
	if tpl.BasePrefix, tpl.BasePostfix, ok = bytes.Cut(tplbuf, []byte("<qtpl-data></qtpl-data>")); !ok {
		panic("index.html <qtpl-data> template delimiter not found")
	}

	m.NewTemplate(baseHost, tpl)
	return tpl.BasePrefix, tpl.BasePostfix
}

func (m Templates) GenerateSettings() ([]byte, []byte) {
	var e error
	var setbuf []byte
	if setbuf, e = openFile("dist/pages/settings.html"); e != nil {
		panic(utils.ExtraErrorWrapper(e, "try to open settings.html file"))
	}

	var ok bool
	tpl := &Template{}
	if tpl.BasePrefix, tpl.BasePostfix, ok = bytes.Cut(setbuf, []byte("<qtpl-data></qtpl-data>")); !ok {
		panic("settings.html <qtpl-data> template delimiter not found")
	}

	// m.NewTemplate(baseHost, tpl)
	return tpl.BasePrefix, tpl.BasePostfix
}
