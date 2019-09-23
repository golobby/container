package main

import (
	"fmt"
	"reflect"
)

type IoC map[interface{}]interface{}

func (i IoC) Bind(f interface{}) {
	i[reflect.TypeOf(f).Out(0).String()] = reflect.ValueOf(f).Call([]reflect.Value{})[0].Elem().Interface()
}

func (i IoC) Make(f interface{}) interface{} {
	if t, ok := f.(string); ok {
		return i[t]
	}

	t := reflect.TypeOf(f).In(0)
	reflect.ValueOf(f).Call([]reflect.Value{reflect.ValueOf(i[t.String()])})

	return nil
}

type Repository interface {
	Find() string
}

type UserRepository struct {
	Name string
}

func (u UserRepository) Find() string {
	return u.Name
}

func main() {
	i := IoC{}
	i.Bind(func() Repository {
		return &UserRepository{}
	})

	x := i.Make("main.Repository")
	fmt.Println(x.(Repository).Find())

	i.Make(func(r Repository) {
		fmt.Println(r.Find())
	})
}
