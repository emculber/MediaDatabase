package main

type Public interface {
	Public() interface{}
}

func (user *UserRole) Public() interface{} {
	return map[string]interface{}{
		"id":      user.Id,
		"user_id": user.UserId,
	}
}
