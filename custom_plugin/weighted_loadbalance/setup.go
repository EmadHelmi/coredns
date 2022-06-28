package weighted_loadbalance

import (
	"io/ioutil"
	"log"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"gopkg.in/yaml.v2"
)

type Content struct {
	Content string `yaml:"content"`
	Weight  int16  `yaml:"weight"`
}

type RecordTypes struct {
	A map[string][]Content `yaml:"A"`
}
type LoadBalancesFile struct {
	RecordTypes RecordTypes `yaml:"record_types"`
}

func init() { plugin.Register("weighted_loadbalance", setup) }

func setup(c *caddy.Controller) error {
	if !c.NextArg() {
		return plugin.Error("weighted_loadbalance", c.ArgErr())
	}

	args := c.RemainingArgs()
	switch args[0] {
	case "round_robin":
		return nil
	case "percentage":
		lb, err := fileParse(c, args[1:])
		if err != nil {
			return plugin.Error("weighted_loadbalance", err)
		}
		dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
			return Percentage{
				Next:    next,
				Configs: lb,
			}
		})
	}

	return nil
}

func fileParse(c *caddy.Controller, args []string) (LoadBalancesFile, error) {
	yamlFile, err := ioutil.ReadFile(args[0])
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		return LoadBalancesFile{}, err
	}

	lb := LoadBalancesFile{}
	err = yaml.Unmarshal(yamlFile, &lb)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return LoadBalancesFile{}, err
	}

	return lb, nil
}
