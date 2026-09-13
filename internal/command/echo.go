package command 


func HandleechoCommand(args []string) string {
	str := ""
	for i := 1; i < len(args); i++ {
		str += args[i]
		if i != len(args)-1 {
			str += " "
		}
	}
	return str;
}