package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"time"

	lab3 "github.com/jamoreno22/dns_1/pkg/proto"
	"google.golang.org/grpc"
)

//DNSServer unimplemented
type DNSServer struct {
	lab3.UnimplementedDNSServer
}

//clock  := time.Now()

func main() {

	// create a listener on TCP port 8000
	lis, err := net.Listen("tcp", "10.10.28.18:8000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// create a server instance
	dnss := DNSServer{}                          // create a gRPC server object
	grpcDNSServer := grpc.NewServer()            // attach the Ping service to the server
	lab3.RegisterDNSServer(grpcDNSServer, &dnss) // start the server

	log.Println("DNSServer_1 running ...")
	if err := grpcDNSServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

}

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
}

// Action server side
func (s *DNSServer) Action(ctx context.Context, cmd *lab3.Command) (*lab3.VectorClock, error) {

	switch cmd.Action {
	case 1: //Create
		// check if file exists
		var _, err = os.Stat("ZF/" + cmd.Domain)

		// create file if not exists
		if os.IsNotExist(err) {
			var file, err1 = os.Create("ZF/" + cmd.Domain)
			if isError(err1) {
				fmt.Printf("File creation error")
			}
			defer file.Close()
		}

		var file, err2 = os.OpenFile("ZF/"+cmd.Domain, os.O_RDWR, 0644)
		if isError(err2) {
			fmt.Printf("File opening error")

		}
		defer file.Close()

		_, err = file.WriteString(cmd.Name + cmd.Domain + " IN A " + cmd.Ip)
		if isError(err) {
			fmt.Printf("File writing error")

		}

	case 2: //Update
		input, err := ioutil.ReadFile("ZF/" + cmd.Domain)
		if err != nil {
			log.Fatalln(err)
		}

		lines := strings.Split(string(input), "\n")

		for i, line := range lines {
			if strings.Contains(line, cmd.Name) {
				if cmd.Option == "Name" {
					lines[i] = cmd.Parameter + cmd.Domain + " IN A " + cmd.Ip
				} else {
					lines[i] = cmd.Name + cmd.Domain + " IN A " + cmd.Parameter
				}
			}
		}
		output := strings.Join(lines, "\n")
		err = ioutil.WriteFile("ZF/"+cmd.Domain, []byte(output), 0644)
		if err != nil {
			log.Fatalln(err)
		}

	case 3: //Delete
		input, err := ioutil.ReadFile("ZF/" + cmd.Domain)
		if err != nil {
			log.Fatalln(err)
		}

		lines := strings.Split(string(input), "\n")
		deleted := false
		for i, line := range lines {
			if deleted == true {
				lines[i-1] = lines[i]
			}
			if strings.Contains(line, cmd.Name) {
				deleted = true
			}
		}
		lines = lines[:len(lines)-1]
		output := strings.Join(lines, "\n")
		err = ioutil.WriteFile("ZF/"+cmd.Domain, []byte(output), 0644)
		if err != nil {
			log.Fatalln(err)
		}
	}

	return &lab3.VectorClock{Rv1: 1, Rv2: 0, Rv3: 0}, nil
}

//Spread server side
func (s *DNSServer) Spread(ctx context.Context, lg *lab3.Log) (*lab3.Message, error) {

	// Client connections to anothers dns servers

	for {
		time.Sleep(5 * time.Minute)
		break
		//send log to dns servers
		//check vector clocks
		//erase log
	}
	return &lab3.Message{Text: "asdf"}, nil
}

//GetIP server side
func (s *DNSServer) GetIP(ctx context.Context, cmd *lab3.Command) (*lab3.PageInfo, error) {

	return &lab3.PageInfo{}, nil
}

/*
	Crear el archivo log para registrar los cambios en el servidor y añadirlo a una variable global
	Añadir los relojes de vector de cada registro ZF
	Solucionar el problema de las conexiones entre DNS sin llamar a la función desde un servidor no DNS
*/
