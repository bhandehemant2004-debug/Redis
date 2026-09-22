package obj


type entry struct{
	entries map[string]string
}

func (e *entry) AddEntry(en []string) {
    for i := 0; i+1 < len(en); i += 2 {
        key := en[i]
        value := en[i+1]
        e.entries[key]=value
    }
}
func NewEntry()*entry{
	return &entry{
		entries: make(map[string]string),
	}
}