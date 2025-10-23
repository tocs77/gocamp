package main

import (
	"fmt"
	"sort"
)

type Person struct {
	Name string
	Age  int
}

type By func(p1, p2 *Person) bool

type personSorter struct {
	people []Person
	by     By
}

func (ps *personSorter) Len() int {
	return len(ps.people)
}

func (ps *personSorter) Swap(i, j int) {
	ps.people[i], ps.people[j] = ps.people[j], ps.people[i]
}

func (ps *personSorter) Less(i, j int) bool {

	return ps.by(&ps.people[i], &ps.people[j])
}

func (by By) Sort(people []Person) {
	ps := &personSorter{people: people, by: by}
	sort.Sort(ps)
}

func main() {
	people := []Person{
		{Name: "John", Age: 30},
		{Name: "Jane", Age: 25},
		{Name: "Jim", Age: 35},
	}

	ageSort := By(func(p1, p2 *Person) bool {
		return p1.Age < p2.Age
	})
	ageSort.Sort(people)

	fmt.Println("ageSort: ", people)

	ageSortDesc := By(func(p1, p2 *Person) bool {
		return p1.Age > p2.Age
	})
	ageSortDesc.Sort(people)

	fmt.Println("ageSortDesc: ", people)

	sort.Slice(people, func(i, j int) bool {
		return people[i].Name < people[j].Name
	})

	fmt.Println("sort.Slice: ", people)
}
