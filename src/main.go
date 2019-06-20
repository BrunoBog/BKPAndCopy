package main

import (
	"log"
	"os"
	"sync"

	"github.com/brunobog/CopyFiles/src/services"
	util "github.com/brunobog/CopyFiles/src/util"
)

func main() {

	var orquestrador sync.WaitGroup
	conf := services.Config{}
	conf.LoadConfig("config.json")

	log.Println("Copy files from ", conf.From)
	log.Println("to path:", conf.To)
	log.Println("Do backup? ", conf.DoBkp)

	if conf.DoBkp {
		util.ZipThis(conf.To, conf.PathBkp)
		os.RemoveAll(conf.To)
	}

	err := util.CopyDirAsync(conf.From, conf.To, &orquestrador)

	if err != nil {
		log.Println(err)
	}
	log.Println("Done!")
}

// func doBkp(pathDir string) (err error) {

// 	if _, err := os.Stat(pathDir); os.IsNotExist(err) {
// 		return err
// 	}

// 	progress := func(archivePath string) {
// 		fmt.Println(archivePath)
// 	}
// outFilePath := filepath.Join(tmpDir, "foo.zip")
// ArchiveFile(pathDir, "bkp."+pathDir+".zip", progress)
// }
