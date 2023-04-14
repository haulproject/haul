package history

import (
	"fmt"
	"time"
)

func Log(id, message string) error {
	fmt.Printf("[%s] %s: %s\n", time.Now(), id, message)
	//TODO insert message into database as a document in the 'histories' collection
	return nil
}

func LogCreate(id, name string) error {
	err := Log(id, fmt.Sprintf("'%s' created", name))
	if err != nil {
		return err
	}

	return nil
}

func LogDelete(id string) error {
	err := Log(id, "deleted")
	if err != nil {
		return err
	}

	return nil
}
