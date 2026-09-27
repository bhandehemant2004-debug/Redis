package command

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/obj"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
	"github.com/codecrafters-io/redis-starter-go/internal/stream"
	"github.com/codecrafters-io/redis-starter-go/internal/util"
)


func HandleXaddCommand(commands []string , storage store.Storage)(string,bool){
	

	key := commands[1]


	object , exists := storage.Get(key)

	if !exists{
		object = util.NewRedisobj("stream",stream.New())
	}

	tree := object.Value.(*stream.Tree)

	entry_key := commands[2]

	last_entry_key,_,ok :=tree.Last() 

	
	if entry_key == "*"{
		return fmt.Sprintf("%d-%d", time.Now().UnixMilli(), 0),false
	}else if entry_key == "0-0" {
    	return "ERR The ID specified in XADD must be greater than 0-0", true
	} else if strings.HasSuffix(entry_key, "*") {
    	if ok {
        	entry_key, _ = validateID(entry_key, string(last_entry_key))
    	} else {
        	entry_key = generateFirstID(entry_key)
    }
	}else if(ok && strings.Compare(entry_key,string(last_entry_key))<=0){
		return "ERR The ID specified in XADD is equal or smaller than the target stream top item",true
	}
	Entry := obj.NewEntry([]byte(entry_key))

	Entry.AddEntry(commands[3:])

	tree.Insert([]byte(entry_key),Entry);


	storage.Set(key,object)

	waiters, exists := storage.GetStreamBlocking(key)

	if exists {
		for _, waiter := range waiters {

			if strings.Compare(entry_key, waiter.StartID) > 0 {
				close(waiter.Ch)
			}
		}
	}

	return entry_key,false
}

func validateID(id, lastID string) (string, bool) {
	toValidate := strings.Split(id, "-")
	validator := strings.Split(lastID, "-")
	ms1, _ := strconv.Atoi(toValidate[0])
	ms2, _ := strconv.Atoi(validator[0])

	seq2, _ := strconv.Atoi(validator[1])
	if toValidate[1] == "*" {
		if ms1 == ms2 {
			toValidate[1] = strconv.Itoa(seq2 + 1)
		} else if ms1 == 0 {
			toValidate[1] = "1"
		} else {
			toValidate[1] = "0"
		}
	}
	seq1, _ := strconv.Atoi(toValidate[1])

	id = toValidate[0] + "-" + toValidate[1]
	if ms1 < ms2 {
		return id, false
	}

	if ms1 == ms2 && seq1 <= seq2 {
		return id, false
	}

	return id, true
}
func generateFirstID(key string)string{
	val := strings.Split(key, "-");
	if val[0]=="0"{
		return val[0]+"-1"
	}
	return val[0]+"-0"
}

