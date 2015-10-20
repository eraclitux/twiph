package main

// FIXME add mongodb as backend.

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/eraclitux/cfgp"
)

type myConf struct {
	OutDir            string `cfgp:",output directory where to store graph,"`
	Csv               string `cfgp:",path to csv file,"`
	Name              string `cfgp:",graph name,"`
	Verified          bool   `cfgp:",get only verified profiles,"`
	ConsumerKey       string `cfgp:",Twitter ConsumerKey,"`
	ConsumerSecret    string `cfgp:",Twitter ConsumerSecret,"`
	AccessToken       string `cfgp:",Twitter AccessToken,"`
	AccessTokenSecret string `cfgp:",Twitter AccessTokenSecret,"`
}

var conf myConf

// ErrorLogger is used to log error messages.
var ErrorLogger *log.Logger

// InfoLogger is used to log general info events like access log.
var InfoLogger *log.Logger

func SetupLoggers(o io.Writer) {
	ErrorLogger = log.New(o, "[ERROR] ", log.Ldate|log.Ltime)
	InfoLogger = log.New(o, "[INFO] ", log.Ldate|log.Ltime)
}

func main() {
	conf = myConf{
		OutDir: "./out",
	}
	err := cfgp.Parse(&conf)
	if err != nil {
		log.Fatal("parsing configuration:", err)
	}
	SetupLoggers(os.Stdout)
	InfoLogger.Println("starting...")
	var r io.ReadCloser
	if conf.Csv != "" {
		r, err = os.Open(conf.Csv)
		if err != nil {
			log.Fatal("opening csv file:", err)
		}
		defer r.Close()
	} else {
		r = os.Stdin
	}
	nodes, err := GetData(r)
	if err != nil {
		log.Fatal("retrieving data:", err)
	}
	var dirName string
	// FIXME use filepath.Join
	if conf.Name != "" {
		dirName = fmt.Sprintf("%s/%s", conf.OutDir, conf.Name)
	} else {
		dirName = fmt.Sprintf("%s/%s", conf.OutDir, time.Now().Format("02-01-2006_15_04_05"))
	}
	err = os.MkdirAll(dirName, 0755)
	if err != nil {
		log.Fatal("creating folder:", err)
	}
	jsonW, err := os.Create(dirName + "/data.json")
	if err != nil {
		log.Fatal("writing json file:", err)
	}
	defer jsonW.Close()
	err = WriteData(jsonW, nodes)
	if err != nil {
		log.Fatal("writing json file:", err)
	}
	indexAvatarW, err := os.Create(dirName + "/index_avatar.html")
	if err != nil {
		log.Fatal("writing index_avatar.html:", err)
	}
	defer indexAvatarW.Close()
	err = createIndexWithAvatar(indexAvatarW)
	if err != nil {
		log.Fatal("writing index_avatar.html:", err)
	}
	indexGroupsW, err := os.Create(dirName + "/index_groups.html")
	if err != nil {
		log.Fatal("writing index_groups.html:", err)
	}
	defer indexGroupsW.Close()
	err = createIndexWithGroups(indexGroupsW)
	if err != nil {
		log.Fatal("writing index_groups.html:", err)
	}
}
