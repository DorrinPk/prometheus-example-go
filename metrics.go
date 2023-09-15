package main

import (
	"fmt"
	"math/rand"
	"os"
	"text/template"
	"time"
)

var tmpl1, _ = template.New("one").Parse(`{{.Mac}}({{.Host}}.ec2.internal)   {{.Count}} IPs ({{.Perc}}% of total)`)
var tmpl2, _ = template.New("two").Parse(`{{.Mac}}({{.Host}}.ec2.internal)   {{.Count}} IPs ({{.Perc}}% of total) ({{.Special}} {{.State}})`)
var tmpls = []*template.Template{tmpl1, tmpl2}
var states = []string{"active", "unreachable"}
var maxCount = 100000

type Input struct {
	Mac     string
	Host    string
	Perc    string
	Count   int
	Special int
	State   string
}

var hosts = map[string]string{
	"3a:2e:91:fa:c8:7b": "ip-10-202-72-5",
	"94:1d:e9:63:fe:42": "ip-10-202-53-8",
	"b8:27:eb:9e:24:cd": "ip-10-202-16-12",
}

func main() {
	rand.Seed(time.Now().Unix())
	for k, v := range hosts {
		count := rand.Intn(maxCount)
		i := Input{
			Mac:     k,
			Host:    v,
			Perc:    fmt.Sprintf("%0.4f", float64(count)/(float64(maxCount*2))),
			Count:   count,
			Special: rand.Intn(200),
			State:   states[rand.Intn(len(states))],
		}
		tmpl := tmpls[rand.Intn(len(tmpls))]
		tmpl.Execute(os.Stdout, i)
		fmt.Print("\n")
	}
}
