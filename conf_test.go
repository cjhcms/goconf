package goconf_test

import (
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {
	// Load configure file
	c, err := conf.New("configure.json")
	if err != nil {
		t.Error(err)
	}

	// Whole configure into map
	RootMap := make(map[string]interface{})
	c.Get("/", &RootMap)

	fmt.Println("Whole configure:\t", RootMap, "\n")

	// Configure into struct
	type ModuleA struct {
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

	moduleAConf := new(ModuleA)
	c.Get("/ModuleA", moduleAConf)

	fmt.Println("ModuleA:\t", *moduleAConf, "\n")
	fmt.Println("ItemIntArray:\t", moduleAConf.SubModule1.ItemIntArray, "\n")

	// One item
	itemIntArray := []int64{}
	c.Get("/ModuleA/SubModule1/itemIntArray", &itemIntArray)

	fmt.Println("ItemIntArray:\t", itemIntArray, "\n")

}
