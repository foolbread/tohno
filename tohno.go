// tohno project main.go
package main

import (
	"flag"
	"fmt"

	"github.com/foolbread/tohno/config"
)

func init() {
	flag.StringVar(&config_file, "f", "conf.ini", "config file path!")
	flag.Parse()

	config.InitConfig(config_file)
}

var config_file string = "conf.ini"

func main() {
	fmt.Println("Hello World!")
}
