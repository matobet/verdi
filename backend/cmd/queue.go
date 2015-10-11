package cmd

import "fmt"

const GlobalQueue = "cmd:global"

func queueByClass(class Class, params map[string]interface{}) (queue string, err error) {
	switch class {
	case Global:
		return GlobalQueue, nil
	default:
		return queueByClassAndParams(class, params)
	}
}

func queueByClassAndParams(class Class, params map[string]interface{}) (queue string, err error) {
	className := classes[class]
	if id, ok := params[className]; ok {
		return QueueByClassAndID(class, id.(string)), nil
	}
	return "", fmt.Errorf(`Commands of class "%s" require '%s' id as parameter`, className, className)
}

func QueueByClassAndID(class Class, id string) string {
	className := classes[class]
	return fmt.Sprintf("cmd:%s:%s", className, id)
}

func replyQueue(id string) string {
	return "cmd:result:" + id
}
