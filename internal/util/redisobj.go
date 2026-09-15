package util

import "time"


type Redisobj struct{
	Type string
	Value any
	Expiryat time.Time
	Expiry bool
}

func NewRedisobj(Type string, value any) Redisobj {
	return Redisobj{
		Type: Type,
		Value: value,
	}
}
func SetExpiry(obj Redisobj, expiry time.Time) Redisobj {
	obj.Expiry = true
	obj.Expiryat = expiry
	return obj
}
