package main

import (
	"log"
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common"
	"os"
	"github.com/natefinch/lumberjack"
	"github.com/lavalamp-/ipv666/common/fs"
	"flag"
	"github.com/lavalamp-/ipv666/common/shell"
)

func setupLogging(conf *config.Configuration) {
	log.Print("Now setting up logging.")
	log.SetFlags(log.Flags() & (log.Ldate | log.Ltime))
  	log.SetOutput(&lumberjack.Logger{
  		Filename:   conf.LogFilePath,
  		MaxSize:    conf.LogFileMBSize,		// megabytes
  		MaxBackups: conf.LogFileMaxBackups,
  		MaxAge:     conf.LogFileMaxAge,		// days
  		Compress:   conf.CompressLogFiles,
  	})
	log.Print("Logging set up successfully.")
}
  
func initializeFilesystem(conf *config.Configuration) (error) {
	log.Print("Now initializing filesystem for IPv6 address discovery process...")
	for _, dirPath := range conf.GetAllDirectories() {
		err := fs.CreateDirectoryIfNotExist(dirPath)
		if err != nil {
			return err
		}
	}
	log.Printf("Initializing state file at '%s'.", conf.GetStateFilePath())
	if _, err := os.Stat(conf.GetStateFilePath()); os.IsNotExist(err) {
		log.Printf("State file does not exist at path '%s'. Creating now.", conf.GetStateFilePath())
		err = common.InitStateFile(conf.GetStateFilePath())
		if err != nil {
			return err
		}
	} else {
		log.Printf("State file already exists at path '%s'.", conf.GetStateFilePath())
	}
	log.Print("Local filesystem initialized for IPv6 address discovery process.")
	return nil
}
  
func main() {

	var configPath string

	flag.StringVar(&configPath, "config", "config.json", "Local file path to the configuration file to use.")

	conf, err := config.LoadFromFile(configPath)

	if err != nil {
		log.Fatal("Can't proceed without loading valid configuration file.")
	}

	setupLogging(&conf)

	err = initializeFilesystem(&conf)

	if err != nil {
		log.Fatal("Error thrown during initialization: ", err)
	}

	zmapAvailable, err := shell.IsZmapAvailable(&conf)

	if err != nil {
		log.Fatal("Error thrown when checking for Zmap: ", err)
	} else if !zmapAvailable {
		log.Fatal("Zmap not found. Please install Zmap.")
	}

	log.Printf("Zmap found and working at path '%s'.", conf.ZmapExecPath)

	log.Print("All systems are green. Entering state machine.")

	err = common.RunStateMachine(&conf)

	if err != nil {
		log.Fatal(err)
	}

	// Ping the router LAN IP address
	//count, err := ping.Ping("2606:6000:6008:AF00:921A:CAFF:FE59:437", time.Duration(100)*time.Millisecond, time.Duration(100)*time.Millisecond, 1, true, false)
	//log.Printf("Ping response count: %d\n", count)

	//fmt.Printf("Hello world\n")
	//addresses, err := common.GetAddressListFromBitStringsFile("/Users/lavalamp/Documents/Projects/IPv6/modeling/files/filtered_ipv6_addrs.dat")
	//addrs, err := addresses.GetAddressListFromBinaryFile("/Users/lavalamp/Documents/Projects/IPv6/modeling/files/ipv6_addresses_2.bin")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//addrModel := modeling.GenerateAddressModel(addrs, "My First Model")
	//addrModel.Save("firstmodel.model")
	//addrModel, err := modeling.GetProbablisticModelFromFile("firstmodel.model")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//addresses := addrModel.GenerateMulti(2, 10000000)
	//addresses, err := common.GetAddressListFromBinaryFile("generated_addys.bin")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//addresses.ToAddressesFile("generated_addys.txt")
	//addresses.ToBinaryFile("generated_addys.bin")
	//log.Printf("Woop woop!")
	//addresses.ToBinaryFile("/Users/lavalamp/Documents/Projects/IPv6/modeling/files/ipv6_addresses_2.bin")
	//log.Printf("Woop woop!")
}