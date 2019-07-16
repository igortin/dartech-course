package main

import (
	"fmt"
	"reflect"
)

func main(){
	b :=Box{}
	b.Put(1,2,3,5,"Alex",9,999)
	fmt.Printf("slice : %v\n", b.interfaces)
	b.Drop("Alex")
	fmt.Printf("slice : %v\n", b.interfaces)


}



type iBox interface{
	Put(...interface{}) // ... variate argument BOx.Put(1,2,5,7,"Peter")
	Drop(interface{})
}
// generic array save int string and we need create type with realize interface Box
type Box struct {
	interfaces []interface{}
}


func (b *Box) Put(a ...interface{}) {
	for _,val := range a {
		b.interfaces = append(b.interfaces, val)
	}
}

func (b *Box) Drop(unknown interface{}){
	for index,i := range b.interfaces {
		if unknown == i && reflect.TypeOf(unknown) == reflect.TypeOf(i) {
			b.interfaces = append(b.interfaces[0:index], b.interfaces[index+1:]...)
		}
	}
}