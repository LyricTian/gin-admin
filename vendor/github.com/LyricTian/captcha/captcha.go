// Copyright 2011 Dmitry Chestnykh. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package captcha implements generation and verification of image and audio
// CAPTCHAs.
//
// A captcha solution is the sequence of digits 0-9 with the defined length.
// There are two captcha representations: image and audio.
//
// An image representation is a PNG-encoded image with the solution printed on
// it in such a way that makes it hard for computers to solve it using OCR.
//
// An audio representation is a WAVE-encoded (8 kHz unsigned 8-bit) sound with
// the spoken solution (currently in English, Russian, Chinese, and Japanese).
// To make it hard for computers to solve audio captcha, the voice that
// pronounces numbers has random speed and pitch, and there is a randomly
// generated background noise mixed into the sound.
//
// This package doesn't require external files or libraries to generate captcha
// representations; it is self-contained.
//
// To make captchas one-time, the package includes a memory storage that stores
// captcha ids, their solutions, and expiration time. Used captchas are removed
// from the store immediately after calling Verify or VerifyString, while
// unused captchas (user loaded a page with captcha, but didn't submit the
// form) are collected automatically after the predefined expiration time.
// Developers can also provide custom store (for example, which saves captcha
// ids and solutions in database) by implementing Store interface and
// registering the object with SetCustomStore.
//
// Captchas are created by calling New, which returns the captcha id.  Their
// representations, though, are created on-the-fly by calling WriteImage or
// WriteAudio functions. Created representations are not stored anywhere, but
// subsequent calls to these functions with the same id will write the same
// captcha solution. Reload function will create a new different solution for
// the provided captcha, allowing users to "reload" captcha if they can't solve
// the displayed one without reloading the whole page.  Verify and VerifyString
// are used to verify that the given solution is the right one for the given
// captcha id.
//
// Server provides an http.Handler which can serve image and audio
// representations of captchas automatically from the URL. It can also be used
// to reload captchas.  Refer to Server function documentation for details, or
// take a look at the example in "capexample" subdirectory.
package captcha

import (
	"bytes"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/LyricTian/captcha/store"
)

const (
	// DefaultLen Default number of digits in captcha solution.
	DefaultLen = 6
	// Expiration time of captchas used by default store.
	Expiration = 10 * time.Minute
)

var (
	// ErrNotFound id not found
	ErrNotFound = errors.New("captcha: id not found")

	globalStore store.Store
	once        sync.Once
)

func getStore() store.Store {
	once.Do(func() {
		if globalStore == nil {
			globalStore = store.NewMemoryStore(time.Second, Expiration)
		}
	})
	return globalStore
}

// SetCustomStore sets custom storage for captchas, replacing the default
// memory store. This function must be called before generating any captchas.
func SetCustomStore(s store.Store) {
	globalStore = s
}

// New creates a new captcha with the standard length, saves it in the internal
// storage and returns its id.
func New() string {
	return NewLen(DefaultLen)
}

// NewLen is just like New, but accepts length of a captcha solution as the
// argument.
func NewLen(length int) (id string) {
	id = randomID()
	getStore().Set(id, RandomDigits(length))
	return
}

// Reload generates and remembers new digits for the given captcha id.  This
// function returns false if there is no captcha with the given id.
//
// After calling this function, the image or audio presented to a user must be
// refreshed to show the new captcha representation (WriteImage and WriteAudio
// will write the new one).
func Reload(id string) bool {
	old := getStore().Get(id, false)
	if old == nil {
		return false
	}
	getStore().Set(id, RandomDigits(len(old)))
	return true
}

// WriteImage writes PNG-encoded image representation of the captcha with the
// given id. The image will have the given width and height.
func WriteImage(w io.Writer, id string, width, height int) error {
	d := getStore().Get(id, false)
	if d == nil {
		return ErrNotFound
	}
	_, err := NewImage(id, d, width, height).WriteTo(w)
	return err
}

// WriteAudio writes WAV-encoded audio representation of the captcha with the
// given id and the given language. If there are no sounds for the given
// language, English is used.
func WriteAudio(w io.Writer, id string, lang string) error {
	d := getStore().Get(id, false)
	if d == nil {
		return ErrNotFound
	}
	_, err := NewAudio(id, d, lang).WriteTo(w)
	return err
}

// Verify returns true if the given digits are the ones that were used to
// create the given captcha id.
//
// The function deletes the captcha with the given id from the internal
// storage, so that the same captcha can't be verified anymore.
func Verify(id string, digits []byte) bool {
	if digits == nil || len(digits) == 0 {
		return false
	}
	reald := getStore().Get(id, true)
	if reald == nil {
		return false
	}
	return bytes.Equal(digits, reald)
}

// VerifyString is like Verify, but accepts a string of digits.  It removes
// spaces and commas from the string, but any other characters, apart from
// digits and listed above, will cause the function to return false.
func VerifyString(id string, digits string) bool {
	if digits == "" {
		return false
	}
	ns := make([]byte, len(digits))
	for i := range ns {
		d := digits[i]
		switch {
		case '0' <= d && d <= '9':
			ns[i] = d - '0'
		case d == ' ' || d == ',':
			// ignore
		default:
			return false
		}
	}
	return Verify(id, ns)
}
