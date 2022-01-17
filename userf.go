package main

// Wrapper save function
func (user User) Save() (err error) {
	err = saveUserData(user.id, *user.data)
	return
}

func IDs(vs []Item) []int {
	vsm := make([]int, len(vs))
	for i, v := range vs {
		vsm[i] = v.ID
	}
	return vsm
}

func (user User) Has(i Item) bool {
	ids := IDs(user.data.inv.Data)
	m := make(map[int]struct{})
	for _, data := range ids {
		m[data] = struct{}{}
	}
	_, e := m[i.ID]
	return e
}
