/*
* @Author: Yajun
* @Date:   2021/12/3 17:06
 */

package abtest

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Domain 流量域
type Domain struct {
	Layers      []*Layer             `json:"layers" validate:"required,unique=Name,gte=1,dive" label:"[流量层]"`
	Experiments map[string]Parameter `json:"experiments" validate:"required" label:"[实验配置]"`
	name        string
}

func validate(data []byte) (*Domain, error) {
	d := Domain{}
	if err := json.Unmarshal(data, &d); err != nil {
		return nil, err
	}
	if err := Validator.Struct(d); err != nil {
		sb := strings.Builder{}
		for k, v := range err.(validator.ValidationErrors).Translate(trans) {
			sb.WriteString(v + " <" + k + ">\n")
		}
		err = errors.New(sb.String())
		return nil, err
	}
	return &d, nil
}

func Parse(name string, data []byte) (*Domain, error) {
	d, err := validate(data)
	if err != nil {
		return nil, err
	}
	d.initLayers()
	d.name = name
	return d, nil
}

func (d *Domain) initLayers() {
	for l := range d.Layers {
		d.Layers[l].init(d)
	}
}

type Result struct {
	Exp    *Traffic `json:"Exp,omitempty"`
	Reason Reason   `json:"Reason,omitempty"`
}

func (d *Domain) Execute(id string) []Result {
	res := make([]Result, 0)
	for _, l := range d.Layers {
		exp, r := l.Allocate(ID(id))
		res = append(res, Result{
			Exp:    exp,
			Reason: r,
		})
	}
	return res
}
