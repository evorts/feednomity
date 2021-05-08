package kv

type manager struct {
	kv map[string]interface{}
	hashed map[string]map[string]interface{}
}

type IManager interface {
	Add(key string, value interface{})
	AddHash(key, hash string, value interface{})
	Get(key string) interface{}
	GetHash(key, hash string) interface{}
	GetAllHash(key string) map[string]interface{}
	Remove(key string)
	RemoveHash(key, hash string)
	RemoveAllHash(key string)
}

func NewKV() IManager {
	return &manager{
		kv:     make(map[string]interface{}),
		hashed: make(map[string]map[string]interface{}),
	}
}

func (m *manager) Add(key string, value interface{}) {
	m.kv[key] = value
}

func (m *manager) AddHash(key, hash string, value interface{}) {
	if _, ok := m.hashed[key]; !ok {
		m.hashed[key] = make(map[string]interface{})
	}
	m.hashed[key][hash] = value
}

func (m *manager) Get(key string) interface{} {
	if v, ok := m.kv[key]; ok {
		return v
	}
	return nil
}

func (m *manager) GetHash(key, hash string) interface{} {
	if _, ok := m.hashed[key]; !ok {
		return nil
	}
	if v, ok := m.hashed[key][hash]; ok {
		return v
	}
	return nil
}

func (m *manager) GetAllHash(key string) map[string]interface{} {
	if v, ok := m.hashed[key]; ok {
		return v
	}
	return nil
}

func (m *manager) Remove(key string) {
	if _, ok := m.kv[key]; ok {
		delete(m.kv, key)
	}
}

func (m *manager) RemoveHash(key, hash string) {
	if _, ok := m.hashed[key]; !ok {
		return
	}
	if _, ok := m.hashed[key][hash]; !ok {
		return
	}
	delete(m.hashed[key], hash)
}

func (m *manager) RemoveAllHash(key string) {
	if _, ok := m.hashed[key]; ok {
		delete(m.hashed, key)
	}
}
