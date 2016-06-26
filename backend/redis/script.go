package redis

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/garyburd/redigo/redis"
)

var scripts = map[string]*redis.Script{}

func (c *Conn) LoadScripts() error {
	names, err := AssetDir("")
	if err != nil {
		return err
	}
	for _, name := range names {
		scriptName, script, err := c.loadScript(name)
		if err != nil {
			return err
		}
		scripts[scriptName] = script
	}

	return nil
}

func (c *Conn) loadScript(assetName string) (scriptName string, script *redis.Script, err error) {
	log.Printf("Loading script: %s", assetName)

	nameWithArity := strings.TrimSuffix(assetName, ".lua")
	i := strings.LastIndex(nameWithArity, "_")
	scriptName, arityStr := nameWithArity[:i], nameWithArity[i+1:]

	content, err := Asset(assetName)
	if err != nil {
		return
	}

	arity, err := strconv.Atoi(arityStr)
	if err != nil {
		return
	}

	script = redis.NewScript(arity, string(content))
	err = script.Load(c)
	return
}

func (c *Conn) DoScript(name string, keysAndArgs ...interface{}) (interface{}, error) {
	script, err := scriptByName(name)
	if err != nil {
		return nil, err
	}
	return script.Do(c, keysAndArgs...)
}

func (c *Conn) SendScript(name string, keysAndArgs ...interface{}) error {
	script, err := scriptByName(name)
	if err != nil {
		return err
	}
	return script.SendHash(c, keysAndArgs...)
}

func scriptByName(name string) (*redis.Script, error) {
	script, ok := scripts[name]
	if !ok {
		return nil, fmt.Errorf("Cannot find redis script: %s", name)
	}
	return script, nil
}
