package analytics

func (m *Analytics) RegisterCallback(fn CallbackFunc) {
	m.cm.Lock()
	defer m.cm.Unlock()
	m.callbacks = append(m.callbacks, fn)
}

func (m *Analytics) forEachCallback(ldstatus bool) {
	m.cm.RLock()
	defer m.cm.RUnlock()

	if len(m.callbacks) == 0 {
		return
	}

	for _, cb := range m.callbacks {
		cb(ldstatus)
	}
}
