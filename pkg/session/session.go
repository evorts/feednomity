package session

import (
	"context"
	"encoding/json"
	"github.com/alexedwards/scs/v2"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

type manager struct {
	sm *scs.SessionManager
}

type Cookie scs.SessionCookie

type IManager interface {
	LoadAndSave(next http.Handler) http.Handler
	Get(ctx context.Context, key string) interface{}
	GetString(ctx context.Context, key string) string
	GetBool(ctx context.Context, key string) bool
	GetInt(ctx context.Context, key string) int
	GetFloat(ctx context.Context, key string) float64
	GetBytes(ctx context.Context, key string) []byte
	GetJson(ctx context.Context, key string, dst interface{}) error
	GetTime(ctx context.Context, key string) time.Time
	Pop(ctx context.Context, key string) interface{}
	PopString(ctx context.Context, key string) string
	PopBool(ctx context.Context, key string) bool
	PopInt(ctx context.Context, key string) int
	PopFloat(ctx context.Context, key string) float64
	PopBytes(ctx context.Context, key string) []byte
	PopTime(ctx context.Context, key string) time.Time
	Put(ctx context.Context, key string, val interface{})
	PutJson(ctx context.Context, key string, val interface{}) error
	Remove(ctx context.Context, key string)
	Clear(ctx context.Context) error
	Exists(ctx context.Context, key string) bool
	Keys(ctx context.Context) []string
	RenewToken(ctx context.Context) error
	Destroy(ctx context.Context) error
	Commit(ctx context.Context) (string, time.Time, error)
	SetSessionLifetime(lifetime time.Duration)
}

func NewSession(expire, idle time.Duration, cookie Cookie) IManager {
	sm := scs.New()
	sm.Lifetime = time.Second * expire
	sm.IdleTimeout = time.Minute * idle
	sm.Cookie.Name = cookie.Name
	sm.Cookie.Domain = cookie.Domain
	sm.Cookie.HttpOnly = cookie.HttpOnly
	sm.Cookie.Persist = cookie.Persist
	sm.Cookie.SameSite = cookie.SameSite
	sm.Cookie.Secure = cookie.Secure
	sm.Cookie.Path = cookie.Path
	return &manager{sm: sm}
}

func (m *manager) LoadAndSave(next http.Handler) http.Handler {
	return m.sm.LoadAndSave(next)
}

func (m *manager) SetSessionLifetime(lifetime time.Duration) {
	m.sm.Lifetime = lifetime
}

func (m *manager) Get(ctx context.Context, key string) interface{} {
	return m.sm.Get(ctx, key)
}

func (m *manager) GetString(ctx context.Context, key string) string {
	return m.sm.GetString(ctx, key)
}

func (m *manager) GetBool(ctx context.Context, key string) bool {
	return m.sm.GetBool(ctx, key)
}

func (m *manager) GetInt(ctx context.Context, key string) int {
	return m.sm.GetInt(ctx, key)
}

func (m *manager) GetFloat(ctx context.Context, key string) float64 {
	return m.sm.GetFloat(ctx, key)
}

func (m *manager) GetBytes(ctx context.Context, key string) []byte {
	return m.sm.GetBytes(ctx, key)
}

func (m *manager) GetJson(ctx context.Context, key string, dst interface{}) error {
	s := m.GetBytes(ctx, key)
	if len(s) < 1 {
		return errors.New("no data found")
	}
	return json.Unmarshal(s, dst)
}

func (m *manager) PutJson(ctx context.Context, key string, val interface{}) error {
	v, err := json.Marshal(val)
	if err == nil {
		m.sm.Put(ctx, key, v)
	}
	return err
}

func (m *manager) GetTime(ctx context.Context, key string) time.Time {
	return m.sm.GetTime(ctx, key)
}

func (m *manager) Pop(ctx context.Context, key string) interface{} {
	return m.sm.Pop(ctx, key)
}

func (m *manager) PopString(ctx context.Context, key string) string {
	return m.sm.PopString(ctx, key)
}

func (m *manager) PopBool(ctx context.Context, key string) bool {
	return m.sm.PopBool(ctx, key)
}

func (m *manager) PopInt(ctx context.Context, key string) int {
	return m.sm.PopInt(ctx, key)
}

func (m *manager) PopFloat(ctx context.Context, key string) float64 {
	return m.sm.PopFloat(ctx, key)
}

func (m *manager) PopBytes(ctx context.Context, key string) []byte {
	return m.sm.PopBytes(ctx, key)
}

func (m *manager) PopTime(ctx context.Context, key string) time.Time {
	return m.sm.PopTime(ctx, key)
}

func (m *manager) Put(ctx context.Context, key string, val interface{}) {
	m.sm.Put(ctx, key, val)
}

func (m *manager) Remove(ctx context.Context, key string) {
	m.sm.Remove(ctx, key)
}

func (m *manager) Clear(ctx context.Context) error {
	return m.sm.Clear(ctx)
}

func (m *manager) Exists(ctx context.Context, key string) bool {
	return m.sm.Exists(ctx, key)
}

func (m *manager) Keys(ctx context.Context) []string {
	return m.sm.Keys(ctx)
}

func (m *manager) RenewToken(ctx context.Context) error {
	return m.sm.RenewToken(ctx)
}

func (m *manager) Destroy(ctx context.Context) error {
	return m.sm.Destroy(ctx)
}

func (m *manager) Commit(ctx context.Context) (string, time.Time, error) {
	return m.sm.Commit(ctx)
}
