package goconf_test

import (
	"fmt"

	"github.com/pantsing/goconf"
)

func ExampleConfigLoadAndGet() {
	c, err := goconf.New("configure.json")
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Configure into struct
	type ModuleAConf struct {
		SubModule1 struct {
			ItemInt      int64
			ItemIntArray []int64

			ItemFloat      float64
			ItemFloatArray []float64

			ItemString      string
			ItemStringArray []string

			ItemBool bool
		}
		SubModule2 struct {
			ItemObject struct {
				Email []string
				Sms   []string
			}
			ItemObjectArray []struct {
				Listen string
				port   uint16 // Unexported field can not receive configure item value
				Port   uint16
			}
		}
	}

	itemModuleA := new(ModuleAConf)
	c.Get("/ModuleA", itemModuleA)

	fmt.Println(*itemModuleA)

	// One item
	itemIntArray := []int64{}
	c.Get("/ModuleA/SubModule1/itemIntArray", &itemIntArray)

	fmt.Println(itemIntArray)

	// Output
	// {{2 [1 2 3 4 5] 2 [1 2 3 4 5] /home/abc/httpd.conf [admin1 admin2] false} {{[abc@360.cn def@360.cn] [130123456789 150123456789]} [{127.0.0.1 0 9000} {127.0.0.2 0 9001} {127.0.0.3 0 9002}]}}
	// [1 2 3 4 5]
}
