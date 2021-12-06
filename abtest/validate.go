/*
* @Author: Yajun
* @Date:   2021/12/3 16:25
 */

package abtest

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	trans, _  = ut.New(zh.New()).GetTranslator("zh")
	Validator = NewValidator()
)

func NewValidator() *validator.Validate {
	validate := validator.New()
	// 注册自定义valid函数
	validate.RegisterStructValidation(ValidDomain, Domain{})
	// 注册自定义valid函数错误翻译
	_ = validate.RegisterTranslation("validShare", trans,
		registerTranslator("validShare", "{0}流量份额之和不为100，重新配置{1}"), translate)
	_ = validate.RegisterTranslation("unique", trans,
		registerTranslator("unique", "{0}中{1}重复"), translate)
	_ = validate.RegisterTranslation("validLayer", trans,
		registerTranslator("validLayer", "{0}中的流量无法对应实验: {1}"), translate)
	// 注册一个函数，获取struct tag里自定义的label作为字段名
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string { return strings.ToLower(fld.Name) })
	//注册翻译器
	_ = zhTranslations.RegisterDefaultTranslations(validate, trans)
	return validate
}

func registerTranslator(tag string, msg string) validator.RegisterTranslationsFunc {
	return func(trans ut.Translator) (err error) {
		if err = trans.Add(tag, msg, false); err != nil {
			return
		}
		return
	}
}

func translate(ut ut.Translator, fe validator.FieldError) string {
	t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
	if err != nil {
		log.Printf("警告: 翻译字段错误: %#v", fe)
		return fe.(error).Error()
	}
	return t
}

func ValidDomain(sl validator.StructLevel) {
	su := sl.Current().Interface().(Domain)
	exps := su.Experiments
	var sum uint
	for i, l := range su.Layers {
		layer := fmt.Sprintf("layers[%d]", i)
		sum = 0
		for _, t := range l.Traffics {
			if _, ok := exps[t.Name]; !ok {
				sl.ReportError(t, layer, "", "validLayer", t.Name)
			}
			sum += t.Share
		}
		if sum != 100 {
			sl.ReportError(l, layer, "", "validShare", l.Name)
		}
	}
}
