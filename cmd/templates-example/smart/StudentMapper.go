package templates

//DO NOT EDIT. This code is auto-generated using go:generate template-generator

type StudentSlice []Student

func (ps StudentSlice) ListToMap() map[string]Student {
	entityMap := make(map[string]Student)
	slice := []Student(ps)
	for i := range slice {
		entityMap[slice[i].ID] = slice[i]
	}
	return entityMap
}
