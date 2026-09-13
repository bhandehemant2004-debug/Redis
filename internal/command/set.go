package command 


func HandleSetCommand(args []string) string {
	if len(args) != 2 {
		return "ERR wrong number of arguments for 'set' command"
	}
	return "OK"
}