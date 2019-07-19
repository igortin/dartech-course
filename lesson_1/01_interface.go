package lesson_1

import (
	"fmt"
	"reflect"
)

func main() {
	b := Box{}
	b.Put(1, 2, 3, 5, "Alex", 9, 999)
	b.Drop(5)
	fmt.Println(b.interfaces)
}

type iBox interface {
	Put(...interface{}) // ... variate argument BOx.Put(1,2,5,7,"Peter")
	Drop(interface{})
}

type Box struct {
	interfaces []interface{}
}

func (b *Box) Put(a ...interface{}) {
	for _, val := range a {
		b.interfaces = append(b.interfaces, val)
	}
}

func (b *Box) Drop(unknown interface{}) {
	for index, i := range b.interfaces {
		if unknown == i && reflect.TypeOf(unknown) == reflect.TypeOf(i) {
			b.interfaces = append(b.interfaces[0:index], b.interfaces[index+1:]...)
		}
	}
}



