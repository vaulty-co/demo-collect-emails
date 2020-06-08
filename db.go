package main

import "sync"

type DB struct {
	emails []email
	mux    sync.Mutex
}

func (d *DB) Save(e email) {
	d.mux.Lock()
	d.emails = append(d.emails, e)
	d.mux.Unlock()
}

func (d *DB) List() []email {
	return d.emails
}
