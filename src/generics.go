package main

import "fmt"

func swap[T any](a, b T) (T, T) {
	return b, a
}

type Stack[T any] struct {
	elements []T
}

func (s *Stack[T]) Push(element T) {
	s.elements = append(s.elements, element)
}

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.elements) == 0 {
		var zero T
		return zero, false
	}
	element := s.elements[len(s.elements)-1]
	s.elements = s.elements[:len(s.elements)-1]
	return element, true
}

func (s *Stack[T]) PrettyPrint() {
	prettyString := "["
	for _, element := range s.elements {
		prettyString += fmt.Sprintf("%v ", element)
	}
	prettyString += "]"
	fmt.Println(prettyString)
}

func main() {
	fmt.Println(swap(1, 2))
	fmt.Println(swap("Hello", "World"))
	fmt.Println(swap(true, false))
	fmt.Println(swap(1.0, 2.0))

	stack := Stack[int]{}
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)
	stack.PrettyPrint()
	fmt.Println(stack.Pop())
	stack.PrettyPrint()
	fmt.Println(stack.Pop())
	fmt.Println(stack.Pop())
	fmt.Println(stack.Pop())

	stringStack := Stack[string]{}
	stringStack.Push("Hello")
	stringStack.Push("World")
	stringStack.PrettyPrint()
	fmt.Println(stringStack.Pop())
	stringStack.PrettyPrint()
}
