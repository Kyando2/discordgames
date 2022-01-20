package main

// Wrapper save function
func (user User) Save() (err error) {
	err = saveUserData(user.id, *user.data)
	return
}

func (user User) Items() Inventory {
	return user.data.inv
}

// Shortcut function for contains(user.Inv().getIDs(), i.ID)
// returns whether the receiver user's inventory contains the item.
func (user User) Has(i Item) bool {
	return contains(user.Items().getIDs(), i.ID)
}
