package def

import (
	"encoding/json"
	"fmt"
	"github.com/hootuu/gelato/errors"
	"strconv"
	"strings"
)

type Payload struct {
	data map[string]string
}

func NewPayload() *Payload {
	return &Payload{data: make(map[string]string)}
}

func PayloadOf(bData []byte) (*Payload, *errors.Error) {
	p := NewPayload()
	err := json.Unmarshal(bData, &p.data)
	if err != nil {
		return nil, errors.System("parse payload failed:" + err.Error())
	}
	return p, nil
}

func (p *Payload) ToBytes() ([]byte, *errors.Error) {
	jsonStr, err := json.Marshal(p.data)
	if err != nil {
		return []byte("{}"), errors.System("Marshal payload failed:" + err.Error())
	}
	return jsonStr, nil
}

func (p *Payload) Set(k string, v string) *Payload {
	p.data[k] = v
	return p
}

func (p *Payload) GetString(k string) (string, *errors.Error) {
	v, ok := p.data[k]
	if !ok {
		return "", errors.System(fmt.Sprintf("no such attribute: %s", k))
	}
	return v, nil
}

func (p *Payload) GetBoolean(k string) (bool, *errors.Error) {
	str, err := p.GetString(k)
	if err != nil {
		return false, errors.System("parse to string failed:" + err.Error())
	}
	if strings.ToLower(str) == "true" {
		return true, nil
	}
	return false, nil
}

func (p *Payload) GetInt64(k string) (int64, *errors.Error) {
	str, err := p.GetString(k)
	if err != nil {
		return 0, err
	}
	iv, nErr := strconv.ParseInt(str, 10, 64)
	if nErr != nil {
		return 0, errors.System(fmt.Sprintf("invalid int value: %s [%s]", str, nErr.Error()))
	}
	return iv, nil
}
