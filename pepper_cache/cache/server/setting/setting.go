package setting

import (
	"flag"
	"log"
)

func Init() {
	flag.StringVar(&Address, "node", DefaultAddress, "node address")
	flag.StringVar(&Cluster, "cluster", "", "cluster address")
	flag.IntVar(&Port, "port", DefaultPort, "port")
	flag.IntVar(&AdminPort, "admin-port", DefaultAdminPort, "port")
	flag.Parse()

	log.Println("Node is", Address)
	if Cluster == "" {
		log.Println("Cluster is", Address)
	} else {
		log.Println("Cluster is", Cluster)
	}
}
