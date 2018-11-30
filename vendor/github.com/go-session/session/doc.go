/*
Package session implements a efficient, safely and easy-to-use session library for Go.


Example:

	package main

	import (
		"context"
		"fmt"
		"net/http"

		"github.com/go-session/session"
	)

	func main() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			store, err := session.Start(context.Background(), w, r)
			if err != nil {
				fmt.Fprint(w, err)
				return
			}

			store.Set("foo", "bar")
			err = store.Save()
			if err != nil {
				fmt.Fprint(w, err)
				return
			}

			http.Redirect(w, r, "/foo", 302)
		})

		http.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
			store, err := session.Start(context.Background(), w, r)
			if err != nil {
				fmt.Fprint(w, err)
				return
			}

			foo, ok := store.Get("foo")
			if ok {
				fmt.Fprintf(w, "foo:%s", foo)
				return
			}
			fmt.Fprint(w, "does not exist")
		})

		http.ListenAndServe(":8080", nil)
	}

Open in your web browser at http://localhost:8080

Output:
	foo:bar

Learn more at https://github.com/go-session/session
*/
package session
