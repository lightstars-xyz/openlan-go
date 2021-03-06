package storage

import (
	"github.com/danieldin95/openlan-go/src/libol"
	"github.com/danieldin95/openlan-go/src/models"
)

type user struct {
	Users *libol.SafeStrMap
}

var User = user{
	Users: libol.NewSafeStrMap(1024),
}

func (w *user) Init(size int) {
	w.Users = libol.NewSafeStrMap(size)
}

func (w *user) Add(user *models.User) {
	libol.Debug("user.Add %v", *user)
	name := user.Name
	if name == "" {
		name = user.Token
	}
	w.Users.Del(name)
	_ = w.Users.Set(name, user)
}

func (w *user) Del(name string) {
	libol.Debug("user.Add %s", name)
	w.Users.Del(name)
}

func (w *user) Get(name string) *models.User {
	if v := w.Users.Get(name); v != nil {
		return v.(*models.User)
	}
	return nil
}

func (w *user) List() <-chan *models.User {
	c := make(chan *models.User, 128)

	go func() {
		w.Users.Iter(func(k string, v interface{}) {
			c <- v.(*models.User)
		})
		c <- nil //Finish channel by nil.
	}()

	return c
}
