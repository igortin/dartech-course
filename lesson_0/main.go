package main

import (
	"dartech-course"
	"fmt"
)

func main() {
	p := dartech_course.Instance()
	A := dartech_course.CreateSong("Paradise a", "Paradise group a", nil, nil)
	B := dartech_course.CreateSong("Paradise b", "Paradise group b", nil, nil)
	C := dartech_course.CreateSong("Paradise c", "Paradise group c", nil, nil)
	W := dartech_course.CreateSong("Paradise w", "Paradise group w", nil, nil)

	dartech_course.Point(&p, &A)

	dartech_course.App(&p, &B)
	dartech_course.App(&p, &C)
	dartech_course.App(&p, &W)

	fmt.Printf("%v\n", dartech_course.GetN(&p))
	fmt.Printf("%v\n", dartech_course.GetN(&p))
	fmt.Printf("%v\n", dartech_course.GetN(&p))
	fmt.Printf("%v\n", dartech_course.GetP(&p))

	//a += 1
	// b := &a      // b = "0X0500999"
	//fmt.Println(&a, a) // 0xc00009e000 2
	// fmt.Println(&b, b, *b) // 0xc0000a0000 0xc00009e000 2
	//*b = *b + 1 // a = 3 and b = 3
	//fmt.Println(a, *b)
}