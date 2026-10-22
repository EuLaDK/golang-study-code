package main

import (
	"fmt"
	"strings"
)

// 变长参数和JavaScript里面的剩余参数差不多
func concat(sep string, args ...interface{}) string {
	var result string
	for i, v := range args {
		if i != 0 {
			result += sep
		}
		switch v.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			result += fmt.Sprintf("%d", v)
		case string:
			result += fmt.Sprintf("%s", v)
		case []int:
			ints := v.([]int)
			for i, v := range ints {
				if i != 0 {
					result += sep
				}
				result += fmt.Sprintf("%d", v)
			}
		case []string:
			strs := v.([]string)
			result += strings.Join(strs, sep)
		default:
			fmt.Printf("the argument type[%T] is not supported", v)
			return ""
		}
	}
	return result
}

// 如果参数在传入时有隐式要求的固定顺序（这点由调用者保证），我们还可以利用变长参数函数模拟实现函数的可选参数和默认参数。
type record struct {
	name    string
	gender  string
	age     uint16
	city    string
	country string
}

func enroll(args ...interface{}) (*record, error) {
	if len(args) > 5 || len(args) < 3 {
		return nil, fmt.Errorf("错误")
	}

	r := &record{
		city:    "北京",
		country: "中国",
	}

	fmt.Println(args, "-------------------------")

	for i, v := range args {
		switch i {
		case 0: // name
			name, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("name is not passed as string")
			}
			r.name = name
		case 1: // gender
			gender, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("gender is not passed as string")
			}
			r.gender = gender
		case 2: // age
			age, ok := v.(int)
			if !ok {
				return nil, fmt.Errorf("age is not passed as int")
			}
			r.age = uint16(age)
		case 3: // city
			city, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("city is not passed as string")
			}
			r.city = city
		case 4: // country
			country, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("country is not passed as string")
			}
			r.country = country
		default:
			return nil, fmt.Errorf("unknown argument passed")
		}
	}
	return r, nil
}

// 功能选项模式
type FinishedHouse struct {
	style                  int    // 0: Chinese; 1:American;2:European
	centralAirConditioning bool   // true or false
	floorMaterial          string //"ground-tile" or "wood"
	wallMaterial           string // "latex" or "paper" or "diatom-mud"
}

type Option func(*FinishedHouse)

func NewFinishedHouse(options ...Option) *FinishedHouse {
	h := &FinishedHouse{
		style:                  0,
		centralAirConditioning: true,
		floorMaterial:          "wood",
		wallMaterial:           "paper",
	}
	for _, option := range options {
		option(h)
	}
	return h
}

func WithStyle(style int) Option {
	return func(h *FinishedHouse) {
		h.style = style
	}
}

func WithFloorMaterial(material string) Option {
	return func(h *FinishedHouse) {
		h.floorMaterial = material
	}
}

func WithWallMaterial(material string) Option {
	return func(h *FinishedHouse) {
		h.wallMaterial = material
	}
}

func WithCentralAirConditioning(centralAirConditioning bool) Option {
	return func(h *FinishedHouse) {
		h.centralAirConditioning = centralAirConditioning
	}
}

func main() {
	println(concat("-", 1, 2))
	println(concat("-", "hello", "gopher"))
	println(concat("-", "hello", 1, uint32(2), []int{11, 12, 13}, 17, []string{"robot", "ai", "ml"}, "hacker", 33))

	r, _ := enroll("小明", "male", 23)
	fmt.Printf("%+v\n", *r)
	r, _ = enroll("小红", "female", 13, "Hangzhou")
	fmt.Printf("%+v\n", *r)
	r, _ = enroll("Leo Messi", "male", 33, "Barcelona", "Spain")
	fmt.Printf("%+v\n", *r)
	r, err := enroll("小吴", 21, "Suzhou")
	if err != nil {
		fmt.Println(err)
		// return
	}
	fmt.Println("--------------------------")
	fmt.Printf("%+v\n", NewFinishedHouse()) // 使用默认选项
	fmt.Printf("%+v\n", NewFinishedHouse(WithStyle(1), WithFloorMaterial("ground-tile"), WithCentralAirConditioning(false)))
}
