package dumb

// Example of dumb code with repetition for all the types

type Profile_ struct {
	ID   string
	Name string
	Age  int
}

type ProfileSlice []Profile_

type Employee_ struct {
	ID     string
	Name   string
	Salary float64
}

type EmployeeSlice []Employee_

type User_ struct {
	ID    string
	Name  string
	Email string
}

type UserSlice []User_

type Student_ struct {
	ID    string
	Name  string
	Marks int
}

type StudentSlice []Student_

func (ps ProfileSlice) ListToMap() map[string]Profile_ {
	entityMap := make(map[string]Profile_)
	slice := []Profile_(ps)
	for i := range slice {
		entityMap[slice[i].ID] = slice[i]
	}
	return entityMap
}

func (ps EmployeeSlice) ListToMap() map[string]Employee_ {
	entityMap := make(map[string]Employee_)
	slice := []Employee_(ps)
	for i := range slice {
		entityMap[slice[i].ID] = slice[i]
	}
	return entityMap
}

func (ps UserSlice) ListToMap() map[string]User_ {
	entityMap := make(map[string]User_)
	slice := []User_(ps)
	for i := range slice {
		entityMap[slice[i].ID] = slice[i]
	}
	return entityMap
}

func (ps StudentSlice) ListToMap() map[string]Student_ {
	entityMap := make(map[string]Student_)
	slice := []Student_(ps)
	for i := range slice {
		entityMap[slice[i].ID] = slice[i]
	}
	return entityMap
}
