package main

import (
	"fmt"
	"net/rpc"
)

type Grade struct {
	Student string
	Course  string
	Grade   float32
}

func client() {
	c, err := rpc.Dial("tcp", "127.0.0.1:5400")
	if err != nil {
		fmt.Println(err)
		return
	}
	n := true
	c.Call("Server.Init", n, &n)
	var op int
	for {
		fmt.Println("1) Register")
		fmt.Println("2) Student average")
		fmt.Println("3) General average")
		fmt.Println("4) Course average")
		fmt.Println("0) Exit")
		fmt.Scanln(&op)

		switch op {
		case 1:
			var name string
			fmt.Print("Name: ")
			fmt.Scanln(&name)

			var course string
			fmt.Print("Course: ")
			fmt.Scanln(&course)

			var grade float32
			fmt.Print("Grade: ")
			fmt.Scanln(&grade)

			g := Grade{name, course, grade}
			var response string
			err = c.Call("Server.RegisterGrades", g, &response)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(response)
		case 2:
			var name string
			fmt.Print("Student: ")
			fmt.Scanln(&name)

			var response float32
			err = c.Call("Server.StudentAverage", name, &response)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(name+":", response)
			}
		case 3:
			var response float32
			err = c.Call("Server.GeneralAverage", "", &response)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("General average:", response)
			}
		case 4:
			var course string
			fmt.Print("Course: ")
			fmt.Scanln(&course)

			var response float32
			err = c.Call("Server.CourseAverage", course, &response)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(course+":", response)
			}
		case 0:
			return
		}
	}
}

func main() {
	client()
}
