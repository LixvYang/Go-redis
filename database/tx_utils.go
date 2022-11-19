package database

func readAllKeys(args [][]byte) ([]string, []string) {
	keys := make([]string, len(args))
	for i, v := range args {
		keys[i] = string(v)
	}
	return nil, keys
}
