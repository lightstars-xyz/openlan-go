package point

import (
	"flag"
	"fmt"
	"strings"

	"github.com/lightstar-dev/openlan-go/libol"
)

type Config struct {
	Addr     string `json:"VsAddr"`
	Auth     string `json:"VsAuth"`
	Verbose  int    `json:"Verbose"`
	Ifmtu    int    `json:"IfMtu"`
	Ifaddr   string `json:"IfAddr"`
	Brname   string `json:"IfBridge"`
	Iftun    bool   `json:"IfTun"`
	Ifethsrc string `json:"IfEthSrc"`
	Ifethdst string `json:"IfEthDst"`

	saveFile string
	name     string
	password string
}

var Default = Config{
	Addr:     "openlan.net",
	Auth:     "hi:hi@123$",
	Verbose:  libol.INFO,
	Ifmtu:    1518,
	Ifaddr:   "",
	Iftun:    false,
	Brname:   "",
	saveFile: ".point.json",
	name:     "",
	password: "",
	Ifethdst: "2e:4b:f0:b7:6d:ba",
	Ifethsrc: "",
}

func RightAddr(listen *string, port int) {
	values := strings.Split(*listen, ":")
	if len(values) == 1 {
		*listen = fmt.Sprintf("%s:%d", values[0], port)
	}
}

func NewConfig() (this *Config) {
	this = &Config{}

	flag.StringVar(&this.Addr, "vs:addr", Default.Addr, "the server connect to")
	flag.StringVar(&this.Auth, "vs:auth", Default.Auth, "the auth login to")
	flag.IntVar(&this.Verbose, "verbose", Default.Verbose, "open verbose")
	flag.IntVar(&this.Ifmtu, "if:mtu", Default.Ifmtu, "the interface MTU include ethernet")
	flag.StringVar(&this.Ifaddr, "if:addr", Default.Ifaddr, "the interface address")
	flag.StringVar(&this.Brname, "if:br", Default.Brname, "the bridge name")
	flag.BoolVar(&this.Iftun, "if:tun", Default.Iftun, "using tun device as interface, otherwise tap")
	flag.StringVar(&this.Ifethdst, "if:ethdst", Default.Ifethdst, "ethernet destination for tun device")
	flag.StringVar(&this.Ifethsrc, "if:ethsrc", Default.Ifethsrc, "ethernet source for tun device")
	flag.StringVar(&this.saveFile, "conf", Default.SaveFile(), "The configuration file")

	flag.Parse()
	libol.SetLog(this.Verbose)

	this.Load()
	this.Default()
	this.Save(fmt.Sprintf("%s.cur", this.saveFile))
	str, err := libol.Marshal(this, false)
	if err != nil {
		libol.Error("NewConfig.json error: %s", err)
	}
	libol.Info("NewConfig.json: %s", str)

	return
}

func (this *Config) Default() {
	if this.Auth != "" {
		values := strings.Split(this.Auth, ":")
		this.name = values[0]
		if len(values) > 1 {
			this.password = values[1]
		}
	}

	RightAddr(&this.Addr, 10002)

	//reset zero value to default
	if this.Addr == "" {
		this.Addr = Default.Addr
	}
	if this.Auth == "" {
		this.Auth = Default.Auth
	}
	if this.Ifmtu == 0 {
		this.Ifmtu = Default.Ifmtu
	}
	if this.Ifaddr == "" {
		this.Ifaddr = Default.Ifaddr
	}
}

func (this *Config) Name() string {
	return this.name
}

func (this *Config) Password() string {
	return this.password
}

func (this *Config) SaveFile() string {
	return this.saveFile
}

func (this *Config) Save(file string) error {
	if file == "" {
		file = this.saveFile
	}

	return libol.MarshalSave(this, file, true)
}

func (this *Config) Load() error {
	if err := libol.UnmarshalLoad(this, this.saveFile); err != nil {
		return err
	}

	return nil
}
