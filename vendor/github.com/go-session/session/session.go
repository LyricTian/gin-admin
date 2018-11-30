package session

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	// ErrInvalidSessionID invalid session id
	ErrInvalidSessionID = errors.New("invalid session id")
)

// Define default options
var defaultOptions = options{
	cookieName:     "go_session_id",
	cookieLifeTime: 3600 * 24 * 7,
	expired:        7200,
	secure:         true,
	sessionID: func() string {
		return newUUID()
	},
	enableSetCookie:     true,
	enableSIDInURLQuery: true,
}

type options struct {
	sign                    []byte
	cookieName              string
	cookieLifeTime          int
	secure                  bool
	domain                  string
	expired                 int64
	sessionID               func() string
	enableSetCookie         bool
	enableSIDInURLQuery     bool
	enableSIDInHTTPHeader   bool
	sessionNameInHTTPHeader string
	store                   ManagerStore
}

// Option A session parameter options
type Option func(*options)

// SetSign Set the session id signature value
func SetSign(sign []byte) Option {
	return func(o *options) {
		o.sign = sign
	}
}

// SetCookieName Set the cookie name
func SetCookieName(cookieName string) Option {
	return func(o *options) {
		o.cookieName = cookieName
	}
}

// SetCookieLifeTime Set the cookie expiration time (in seconds)
func SetCookieLifeTime(cookieLifeTime int) Option {
	return func(o *options) {
		o.cookieLifeTime = cookieLifeTime
	}
}

// SetDomain Set the domain name of the cookie
func SetDomain(domain string) Option {
	return func(o *options) {
		o.domain = domain
	}
}

// SetSecure Set cookie security
func SetSecure(secure bool) Option {
	return func(o *options) {
		o.secure = secure
	}
}

// SetExpired Set session expiration time (in seconds)
func SetExpired(expired int64) Option {
	return func(o *options) {
		o.expired = expired
	}
}

// SetSessionID Set callback function to generate session id
func SetSessionID(sessionID func() string) Option {
	return func(o *options) {
		o.sessionID = sessionID
	}
}

// SetEnableSetCookie Enable writing session id to cookie
// (enabled by default, can be turned off if no cookie is written)
func SetEnableSetCookie(enableSetCookie bool) Option {
	return func(o *options) {
		o.enableSetCookie = enableSetCookie
	}
}

// SetEnableSIDInURLQuery Allow session id from URL query parameters (enabled by default)
func SetEnableSIDInURLQuery(enableSIDInURLQuery bool) Option {
	return func(o *options) {
		o.enableSIDInURLQuery = enableSIDInURLQuery
	}
}

// SetEnableSIDInHTTPHeader Allow session id to be obtained from the request header
func SetEnableSIDInHTTPHeader(enableSIDInHTTPHeader bool) Option {
	return func(o *options) {
		o.enableSIDInHTTPHeader = enableSIDInHTTPHeader
	}
}

// SetSessionNameInHTTPHeader The key name in the request header where the session ID is stored
// (if it is empty, the default is the cookie name)
func SetSessionNameInHTTPHeader(sessionNameInHTTPHeader string) Option {
	return func(o *options) {
		o.sessionNameInHTTPHeader = sessionNameInHTTPHeader
	}
}

// SetStore Set session management storage
func SetStore(store ManagerStore) Option {
	return func(o *options) {
		o.store = store
	}
}

// NewManager Create a session management instance
func NewManager(opt ...Option) *Manager {
	opts := defaultOptions
	for _, o := range opt {
		o(&opts)
	}

	if opts.enableSIDInHTTPHeader && opts.sessionNameInHTTPHeader == "" {
		opts.sessionNameInHTTPHeader = opts.cookieName
	}

	if opts.store == nil {
		opts.store = NewMemoryStore()
	}
	return &Manager{opts: &opts}
}

// Manager A session management instance, including start and destroy operations
type Manager struct {
	opts *options
}

func (m *Manager) getContext(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx = newReqContext(ctx, r)
	ctx = newResContext(ctx, w)
	return ctx
}

func (m *Manager) signature(sid string) string {
	h := hmac.New(sha1.New, m.opts.sign)
	h.Write([]byte(sid))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (m *Manager) decodeSessionID(value string) (string, error) {
	value, err := url.QueryUnescape(value)
	if err != nil {
		return "", err
	}

	vals := strings.Split(value, ".")
	if len(vals) != 2 {
		return "", ErrInvalidSessionID
	}

	bsid, err := base64.StdEncoding.DecodeString(vals[0])
	if err != nil {
		return "", err
	}
	sid := string(bsid)

	sign := m.signature(sid)
	if sign != vals[1] {
		return "", ErrInvalidSessionID
	}
	return sid, nil
}

func (m *Manager) sessionID(r *http.Request) (string, error) {
	var cookieValue string

	cookie, err := r.Cookie(m.opts.cookieName)
	if err == nil && cookie.Value != "" {
		cookieValue = cookie.Value
	} else {
		if m.opts.enableSIDInURLQuery {
			err := r.ParseForm()
			if err != nil {
				return "", err
			}
			cookieValue = r.FormValue(m.opts.cookieName)
		}

		if m.opts.enableSIDInHTTPHeader && cookieValue == "" {
			cookieValue = r.Header.Get(m.opts.sessionNameInHTTPHeader)
		}
	}

	if cookieValue != "" {
		return m.decodeSessionID(cookieValue)
	}

	return "", nil
}

func (m *Manager) encodeSessionID(sid string) string {
	b := base64.StdEncoding.EncodeToString([]byte(sid))
	s := fmt.Sprintf("%s.%s", b, m.signature(sid))
	return url.QueryEscape(s)
}

func (m *Manager) isSecure(r *http.Request) bool {
	if !m.opts.secure {
		return false
	}
	if r.URL.Scheme != "" {
		return r.URL.Scheme == "https"
	}
	if r.TLS == nil {
		return false
	}
	return true
}

func (m *Manager) setCookie(sessionID string, w http.ResponseWriter, r *http.Request) {
	cookieValue := m.encodeSessionID(sessionID)
	cookie := &http.Cookie{
		Name:     m.opts.cookieName,
		Value:    cookieValue,
		Path:     "/",
		HttpOnly: true,
		Secure:   m.isSecure(r),
		Domain:   m.opts.domain,
	}

	if v := m.opts.cookieLifeTime; v > 0 {
		cookie.MaxAge = v
		cookie.Expires = time.Now().Add(time.Duration(v) * time.Second)
	}

	if m.opts.enableSetCookie {
		http.SetCookie(w, cookie)
	}

	r.AddCookie(cookie)

	if m.opts.enableSIDInHTTPHeader {
		key := m.opts.sessionNameInHTTPHeader
		r.Header.Set(key, cookieValue)
		w.Header().Set(key, cookieValue)
	}
}

// Start a session and return to session storage
func (m *Manager) Start(ctx context.Context, w http.ResponseWriter, r *http.Request) (Store, error) {
	ctx = m.getContext(ctx, w, r)

	sid, err := m.sessionID(r)
	if err != nil {
		return nil, err
	}

	if sid != "" {
		if exists, verr := m.opts.store.Check(ctx, sid); verr != nil {
			return nil, verr
		} else if exists {
			return m.opts.store.Update(ctx, sid, m.opts.expired)
		}
	}

	store, err := m.opts.store.Create(ctx, m.opts.sessionID(), m.opts.expired)
	if err != nil {
		return nil, err
	}

	m.setCookie(store.SessionID(), w, r)
	return store, nil
}

// Refresh a session and return to session storage
func (m *Manager) Refresh(ctx context.Context, w http.ResponseWriter, r *http.Request) (Store, error) {
	ctx = m.getContext(ctx, w, r)

	sid, err := m.sessionID(r)
	if err != nil {
		return nil, err
	} else if sid == "" {
		sid = m.opts.sessionID()
	}

	store, err := m.opts.store.Refresh(ctx, sid, m.opts.sessionID(), m.opts.expired)
	if err != nil {
		return nil, err
	}

	m.setCookie(store.SessionID(), w, r)
	return store, nil
}

// Destroy a session
func (m *Manager) Destroy(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx = m.getContext(ctx, w, r)

	sid, err := m.sessionID(r)
	if err != nil {
		return err
	} else if sid == "" {
		return nil
	}

	exists, err := m.opts.store.Check(ctx, sid)
	if err != nil {
		return err
	} else if !exists {
		return nil
	}

	err = m.opts.store.Delete(ctx, sid)
	if err != nil {
		return err
	}

	if m.opts.enableSIDInHTTPHeader {
		key := m.opts.sessionNameInHTTPHeader
		r.Header.Del(key)
		w.Header().Del(key)
	}

	if m.opts.enableSetCookie {
		cookie := &http.Cookie{
			Name:     m.opts.cookieName,
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Now(),
			MaxAge:   -1,
		}

		http.SetCookie(w, cookie)
	}

	return nil
}
