package goconf

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"unicode"
)

type Config struct {
	items map[string]interface{}
}

// New function initials a Config and load a JSON format configure file in it.
func New(filename string) (c *Config, err error) {
	c = NewConfig()
	err = c.Load(filename)
	return
}

// NewConfig() generates an empty Config struct
func NewConfig() (c *Config) {
	c = new(Config)
	c.items = make(map[string]interface{})
	return
}

// Load loads JSON format configure file.
func (c *Config) Load(filename string) (err error) {
	if filename == "" {
		return errors.New("Need configure file")
	}
	fi, err := os.Stat(filename)
	if err != nil {
		return
	}
	if fi.IsDir() {
		return errors.New(filename + "is a directory")
	}

	// Open configure file
	fileHandler, err := os.Open(filename)
	if err != nil {
		return
	}
	defer fileHandler.Close()
	// Get file size
	fsize := fi.Size()
	// Make file buffer
	buf := make([]byte, fsize+1)
	// Read file to buffer
	_, err = fileHandler.Read(buf)
	if err != nil {
		return
	}

	// Remove comment
	// The char before single commment line requires start with \s*// to different from URLs
	re, err := regexp.Compile(`((?m)^\s*//.*\n)|((?ms)/\*.*?\*/)`)
	buf = re.ReplaceAllLiteral(buf[0:fsize], []byte(""))

	// Unmarshal JSON
	c.items = make(map[string]interface{}) // configure map key is string, value could be any type
	err = json.Unmarshal(buf, &c.items)

	return
}

// Get function read configure item value to the pointer v pointed variate by path, etc. "/ModuleA/SubModule1/itemIntArray"
// The type of variate that pointer v pointed to MUST be the same type with item value.
// JSON object type should correspond to struct or map[string]interface{} type in Golang.
// Get 函数按照配置项的在配置文件中的层次路径,如 "/ModuleA/SubModule1/itemIntArray", 读取配置项的值到指针变量v指向的变量中.
// v指向的变量必须与配置值数据类型相符，JSON对象类型应对应Golang中Struct类型或Map[string]interface{}类型.
func (c *Config) Get(nodePath string, v interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			if str, ok := r.(string); ok {
				err = errors.New(str)
			} else {
				err = r.(error)
			}
		}
	}()

	nodes := strings.Split(nodePath, "/")
	// Here MUST has type assert first before use conf["nodename"] and its child node value
	var nodeValue interface{}
	nodeValue = c.items
	for i := 1; i < len(nodes); i++ { // i := 1 to skip first root '/'
		nv, ok := nodeValue.(map[string]interface{})
		if ok {
			if nodes[i] != "" {
				nodeValue = nv[nodes[i]]
			} else {
				nodeValue = nv
			}
		} else {
			return errors.New("Configuration node " + nodes[i] + " is not a map[string]interface{} type.")
		}
	}

	if nodeValue == nil {
		return errors.New("Configuration node value is nil, node name may be mistake.")
	}

	if v == nil {
		return errors.New("Value receiver is nil.")
	}

	if reflect.TypeOf(v).Kind() == reflect.Ptr {
		if reflect.TypeOf(v).Elem().Kind() == reflect.Struct {
			err = c.setStructValue(nodeValue, v)
		} else {
			err = c.getJsonValue(nodeValue, v)
		}
	} else {
		err = errors.New("Value receiver type MUST be a pointer.")
	}
	return
}

func (c *Config) setStructValue(nodeValue interface{}, v interface{}) (err error) {
	if nv, ok := nodeValue.(map[string]interface{}); ok {
		if v != nil && reflect.ValueOf(v).MethodByName("UnmarshalJSON").IsValid() { // Suport JSON Unmarshal
			bytes, err := json.Marshal(nv)
			if err != nil {
				return err
			} else {
				r := reflect.ValueOf(v).MethodByName("UnmarshalJSON").Call([]reflect.Value{reflect.ValueOf(bytes)})
				if len(r) == 1 {
					if r[0].Interface() == nil {
						return nil
					} else if err, ok := r[0].Interface().(error); ok {
						return err
					}
				}
				return errors.New("Invalid UnmarshalJSON function")
			}
		} else {
			for name, value := range nv {
				fv := reflect.ValueOf(v).Elem().FieldByName(c.firstCharToUpper(name))
				if fv.IsValid() {
					if fv.Type().Kind() == reflect.Struct {
						err = c.setStructValue(value, fv.Addr().Interface())
					} else {
						err = c.getJsonValue(value, fv.Addr().Interface())
					}
				}
			}
		}
	}
	return
}

func (c *Config) getJsonValue(nodeValue interface{}, v interface{}) (err error) {
	switch nodeValue.(type) {
	case float64:
		switch v.(type) {
		case *float64:
			v, ok := v.(*float64)
			if ok {
				*v = nodeValue.(float64)
			} else {
				err = errors.New("Value receiver type is not *float64.")
			}
		case *int:
			v, ok := v.(*int)
			if ok {
				*v = int(nodeValue.(float64))
			} else {
				err = errors.New("Value receiver type is not *int.")
			}
		case *int8:
			v, ok := v.(*int8)
			if ok {
				*v = int8(nodeValue.(float64))
			} else {
				err = errors.New("Value receiver type is not *int8.")
			}
		case *int16:
			v, ok := v.(*int16)
			if ok {
				*v = int16(nodeValue.(float64))
			} else {
				err = errors.New("Value receiver type is not *int16.")
			}
		case *int32:
			v, ok := v.(*int32)
			if ok {
				*v = int32(nodeValue.(float64))
			} else {
				err = errors.New("Value receiver type is not *int32.")
			}
		case *int64:
			v, ok := v.(*int64)
			if ok {
				*v = int64(nodeValue.(float64))
			} else {
				err = errors.New("Value receiver type is not *int64.")
			}
		case *uint:
			v, ok := v.(*uint)
			if ok {
				*v = uint(nodeValue.(float64))
			} else {
				err = errors.New("Value receiver type is not *uint.")
			}
		case *uint8:
			v, ok := v.(*uint8)
			if ok {
				tmpV := nodeValue.(float64)
				*v = uint8(tmpV)
			} else {
				err = errors.New("Value receiver type is not *uint16.")
			}
		case *uint16:
			v, ok := v.(*uint16)
			if ok {
				tmpV := nodeValue.(float64)
				*v = uint16(tmpV)
			} else {
				err = errors.New("Value receiver type is not *uint16.")
			}
		case *uint32:
			v, ok := v.(*uint32)
			if ok {
				*v = uint32(nodeValue.(float64))
			} else {
				err = errors.New("Value receiver type is not *uint32.")
			}
		case *uint64:
			v, ok := v.(*uint64)
			if ok {
				*v = uint64(nodeValue.(float64))
			} else {
				err = errors.New("Value receiver type is not *uint64.")
			}
		default:
			err = errors.New("Value receiver type is not *float64.")
		}
	case string:
		v, ok := v.(*string)
		if ok {
			*v = nodeValue.(string)
		} else {
			err = errors.New("Value receiver type is not *string.")
		}
	case map[string]interface{}:
		v, ok := v.(*map[string]interface{})
		if ok {
			*v = nodeValue.(map[string]interface{})
		} else {
			err = errors.New("Value receiver type is not *map[string]interface{}.")
		}
	case []interface{}:
		nodeArrayValue := nodeValue.([]interface{})
		if reflect.TypeOf(v).Kind() == reflect.Ptr && reflect.TypeOf(v).Elem().Kind() == reflect.Slice && reflect.ValueOf(v).Elem().CanSet() {
			sv := reflect.ValueOf(v).Elem()
			if reflect.TypeOf(v).Elem().Elem().Kind() == reflect.Struct {
				for _, nv := range nodeArrayValue {
					if nsv, ok := nv.(map[string]interface{}); ok {
						tmpV := reflect.New(reflect.TypeOf(v).Elem().Elem())
						err = c.setStructValue(nsv, tmpV.Interface())
						sv.Set(reflect.Append(sv, tmpV.Elem()))
					} else {
						err = errors.New("Array node element value is not a struct.")
					}
				}
			} else {
				for _, nv := range nodeArrayValue {
					tmpV := reflect.New(reflect.TypeOf(v).Elem().Elem())
					err = c.getJsonValue(nv, tmpV.Interface())
					sv.Set(reflect.Append(sv, tmpV.Elem()))
				}
			}
		} else {
			err = errors.New("Value receiver type is not pointer of slice.")
		}
	default:
		err = errors.New("Node value type is not correct")
	}

	return err
}

// firstCharToUpper upper first character of string to upper case,
// It is used to convert name of configure item key to EXPORTED struct field name, which first character MUST be upper case.
func (c *Config) firstCharToUpper(str string) string {
	if len(str) > 0 {
		runes := []rune(str)
		firstRune := unicode.ToUpper(runes[0])
		str = string(firstRune) + string(runes[1:])
	}
	return str
}
