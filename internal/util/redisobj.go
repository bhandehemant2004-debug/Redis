package util

import "time"


type Redisobj struct{

	Value string 
	Expiryat time.Time
	Expiry bool
}