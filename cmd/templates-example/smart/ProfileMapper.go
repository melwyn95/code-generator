package templates

//DO NOT EDIT. This code is auto-generated using go:generate template-generator

type ProfileSlice []Profile

func (ps ProfileSlice) ListToMap() map[string]Profile {
	entityMap := make(map[string]Profile)
	slice := []Profile(ps)
	for i := range slice {
		entityMap[slice[i].ID] = slice[i]
	}
	return entityMap
}
