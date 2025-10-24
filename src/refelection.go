package main

import (
	"fmt"
	"reflect"
)

type Person struct {
	Name string
	Age  int
	id   int
}

type Greeter struct{}

func (g Greeter) Greet(name string) string {
	return "Hello, " + name + "!"
}

func main() {
	x := 42
	v := reflect.ValueOf(x)
	fmt.Println("v: ", v)
	fmt.Println("v.Type(): ", v.Type())
	fmt.Println("v.Kind(): ", v.Kind())
	fmt.Println("Is int type: ", v.Kind() == reflect.Int)
	fmt.Println("v.Interface(): ", v.Interface())
	fmt.Println("v.IsZero(): ", v.IsZero())
	fmt.Println("v.IsValid(): ", v.IsValid())

	y := 10
	v2 := reflect.ValueOf(&y).Elem()
	fmt.Println("Original value of y: ", v2.Int())
	v2.SetInt(20)
	fmt.Println("New value of y: ", v2.Int())

	var itf interface{} = "Hello, World!"
	v3 := reflect.ValueOf(itf)
	fmt.Println("v3: ", v3)
	fmt.Println("v3.Type(): ", v3.Type())
	fmt.Println("v3.Kind(): ", v3.Kind())
	fmt.Println("v3.Interface(): ", v3.Interface())
	fmt.Println("v3.IsZero(): ", v3.IsZero())
	fmt.Println("v3.IsValid(): ", v3.IsValid())

	p := Person{Name: "John", Age: 30, id: 123}
	v4 := reflect.ValueOf(p)
	fmt.Println("v4: ", v4)
	for i := 0; i < v4.NumField(); i++ {
		if !v4.Field(i).CanSet() {
			continue
		}
		fmt.Println("Field: ", v4.Field(i))
		fmt.Println("Field type: ", v4.Field(i).Type())
		fmt.Println("Field kind: ", v4.Field(i).Kind())
		fmt.Println("Field interface: ", v4.Field(i).Interface())
		fmt.Println("Field is zero: ", v4.Field(i).IsZero())
		fmt.Println("Field is valid: ", v4.Field(i).IsValid())
	}
	v5 := reflect.ValueOf(&p).Elem()
	nameField := v5.FieldByName("Name")
	if nameField.IsValid() && nameField.CanSet() {
		nameField.SetString("Jane")
	}
	fmt.Println("New name: ", p.Name)
	idField := v5.FieldByName("id")
	if idField.IsValid() && idField.CanSet() {
		idField.SetInt(123)
	} else {
		fmt.Println("id field is not valid or cannot be set")
	}
	fmt.Println("New id: ", p.id)

	g := Greeter{}
	v6 := reflect.TypeOf(g)
	v7 := reflect.ValueOf(g)
	for i := range v6.NumMethod() {
		method := v6.Method(i)
		fmt.Println("Method: ", method.Name)
	}
	method := v7.MethodByName("Greet")
	if method.IsValid() {
		result := method.Call([]reflect.Value{reflect.ValueOf("Bob")})
		fmt.Println("Result: ", result[0].String())
	}
}
