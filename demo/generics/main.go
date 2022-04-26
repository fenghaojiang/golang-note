package main

import (
	"fmt"
	"golang.org/x/exp/constraints"
)

func GMin [T constraints.Ordered] (x, y T) T {
	if x < y {
		return x
	}
	return y
}

type Point []int32



// Scale returns a copy of s with each element multiplied by c.
func Scale[S ~[]E, E constraints.Integer](s S, c E) S {
    r := make(S, len(s))
    for i, v := range s {
        r[i] = v * c
    }
    return r
}

func main() {
	x := GMin[int](2, 3)
	fmt.Println(x)
	fmin := GMin[float64]
	m := fmin(1.1, 2.2)
	var point Point = Point{1,3,4,5,56}
	fmt.Println(m)

	fmt.Println(Scale[Point](point, 2))
	engineer := Engineer{Salary: 10}
	manager := Manager{Salary: 100}
	getSalary(engineer)
	getSalary(manager)
}



type Employee interface {
	PrintSalary()
}

func getSalary[E Employee] (e E) {
	e.PrintSalary()
}


type Engineer struct {
	Salary int32
}

func (e Engineer) PrintSalary() {
	fmt.Println(e.Salary)
}

type Manager struct {
	Salary int64
}

func (m Manager) PrintSalary() {
	fmt.Println(m.Salary)
}

