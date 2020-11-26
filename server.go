package main

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
)

type Grade struct {
	Student string
	Course  string
	Grade   float32
}

type Server struct {
	Courses  map[string]map[string]float32
	Students map[string]map[string]float32
}

func (this *Server) Init(a bool, b *bool) error {
	this.Courses = make(map[string]map[string]float32)
	this.Students = make(map[string]map[string]float32)
	return nil
}

func Print(students, courses map[string]map[string]float32) {
	fmt.Println("Students:")
	for s := range students {
		fmt.Println(s + " {")
		for c := range students[s] {
			fmt.Println("  " + c + ": " + fmt.Sprintf("%.2f", students[s][c]))
		}
		fmt.Println("}")
	}

	fmt.Println("\nCourses:")
	for c := range courses {
		fmt.Println(c + " {")
		for s := range courses[c] {
			fmt.Println("  " + s + ": " + fmt.Sprintf("%.2f", courses[c][s]))
		}
		fmt.Println("}")
	}
	fmt.Println("---------------")
}

func (this *Server) RegisterGrades(g Grade, response *string) error {
	if _, ok := this.Courses[g.Course]; ok {
		if _, ok := this.Courses[g.Course][g.Student]; ok {
			*response = "fail"
			return errors.New(g.Student + " already has a grade for " + g.Course)
		}
		this.Courses[g.Course][g.Student] = g.Grade
	} else {
		this.Courses[g.Course] = map[string]float32{g.Student: g.Grade}
	}

	if _, ok := this.Students[g.Student]; ok {
		this.Students[g.Student][g.Course] = g.Grade
	} else {
		this.Students[g.Student] = map[string]float32{g.Course: g.Grade}
	}

	Print(this.Students, this.Courses)

	*response = "success"
	return nil
}

func (this *Server) StudentAverage(student string, average *float32) error {
	if _, ok := this.Students[student]; !ok {
		return errors.New(student + " does not exist")
	}
	var temp float32
	for _, v := range this.Students[student] {
		temp += v
	}
	*average = temp / float32(len(this.Students[student]))
	return nil
}

func (this *Server) GeneralAverage(_ string, average *float32) error {
	if len(this.Students) == 0 {
		return errors.New("No students rigistered")
	}
	var sg, gg float32
	for student := range this.Students {
		this.StudentAverage(student, &sg)
		gg += sg
	}
	*average = gg / float32(len(this.Students))
	return nil
}

func (this *Server) CourseAverage(course string, average *float32) error {
	if _, ok := this.Courses[course]; !ok {
		return errors.New(course + " does not exist")
	}
	var temp float32
	for _, v := range this.Courses[course] {
		temp += v
	}
	*average = temp / float32(len(this.Courses[course]))
	return nil
}

func server() {
	rpc.Register(new(Server))
	ln, err := net.Listen("tcp", ":5400")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Listening on port 5400")
	for {
		c, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go rpc.ServeConn(c)
	}
}

func main() {
	go server()

	var input string
	fmt.Scanln(&input)
}
