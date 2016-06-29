package main

type Public interface {
	Public() interface{}
}

func (u *UserRole) Public() interface{} {
	return map[string]interface{}{
		"id":      u.Id.Encode(),
		"user_id": u.UserId,
	}
}
