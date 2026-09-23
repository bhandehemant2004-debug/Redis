package obj


type Entry struct{
    Id []byte
	Entries map[string]string
}

func (e *Entry) AddEntry(en []string) {
    for i := 0; i+1 < len(en); i += 2 {
        key := en[i]
        value := en[i+1]
        e.Entries[key]=value
    }
}
func NewEntry(ids[]byte)*Entry{

	return &Entry{
        Id: ids,
		Entries: make(map[string]string),
	}
}

func (e *Entry)Entryid()[]byte{
    return e.Id
}

