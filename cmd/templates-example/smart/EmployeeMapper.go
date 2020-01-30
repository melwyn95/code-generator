package templates

//DO NOT EDIT. This code is auto-generated using go:generate template-generator

type EmployeeSlice []Employee

func (ps EmployeeSlice) ListToMap() map[string]Employee {
	entityMap := make(map[string]Employee)
	slice := []Employee(ps)
	for i := range slice {
		entityMap[slice[i].ID] = slice[i]
	}
	return entityMap
}
