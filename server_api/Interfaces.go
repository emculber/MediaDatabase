package main

type Public interface {
	Public() interface{}
}

func (userKeys *UserKeys) Public() interface{} {
	return map[string]interface{}{
		"id":      userKeys.Id,
		"user_id": userKeys.User.Id,
	}
}
