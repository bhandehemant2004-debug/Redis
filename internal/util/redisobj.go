package util

import "time"


type Redisobj struct{

	value string 
	expiryat time.Time
	expirty bool
}