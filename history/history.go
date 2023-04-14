package history

import "log"

func Log(message string) error {
	log.Printf("history: %s\n", message)
	return nil
}
