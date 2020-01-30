package templates

//go:generate template-generator

type Profile struct {
	ID   string
	Name string
	Age  int
}

type Employee struct {
	ID     string
	Name   string
	Salary float64
}

type User struct {
	ID    string
	Name  string
	Email string
}

type Student struct {
	ID    string
	Name  string
	Marks int
}
