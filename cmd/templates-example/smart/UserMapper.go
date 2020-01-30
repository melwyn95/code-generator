package templates

//DO NOT EDIT. This code is auto-generated using go:generate template-generator

type UserSlice []User

func (ps UserSlice) ListToMap() map[string]User {
	entityMap := make(map[string]User)
	slice := []User(ps)
	for i := range slice {
		entityMap[slice[i].ID] = slice[i]
	}
	return entityMap
}
