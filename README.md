# Lab Reservation Service End-To-End Workflow.
## Prerequisites
* Testbed API SDK (Ondatra/Cafy)
    * import goopentestbed for Ondatra ```github.com/open-traffic-generator/opentestbed/goopentestbed```
    * import opentestbed For Cafy ```pip install opentestbed```
* Docker Netbox should be running, to get the Inventory data (Dynamic updates have not yet been implemented, currently using manual feed inventory data)
* Lab Reservation Service should be running.

## Running the Client for Ondatra
* Require testbed input in JSON.
* Import goopentestbed
* Run the Go file and the file data should have the below content.
* For the test run use the dynamically generated testbed_output.
    ```
    package main
    import (
    	"fmt"
    	"os"
    	"github.com/open-traffic-generator/opentestbed/goopentestbed"
    	bindpb "github.com/openconfig/featureprofiles/topologies/proto/binding"
    	"google.golang.org/protobuf/encoding/prototext"
    )
    func main() {
    
    	api := goopentestbed.NewApi()
    	api.NewHttpTransport().SetLocation("http://127.0.0.1:8080")
    	testbed := goopentestbed.NewTestbed()	
    	in, err := os.ReadFile("testbed.json")
    	if err != nil {
    		fmt.Printf("could not read file: %v\n", err)
    		return
    	}
    	testbed.Unmarshal().FromJson(string(in))
    	reservationResult, _ := api.Reserve(testbed)
    	testbed_output := &bindpb.Binding{}
    	if err := prototext.Unmarshal([]byte(*reservationResult), testbed_output); err != nil {
    		fmt.Printf("Error unmarshalling Prototext: %v\n", err)
    		return
    	}
    	fmt.Println("Unmarshalled binding:", testbed_output)	
    }
    ```
* Testbed.json content.
    ```
    {
        "devices": [
            {
                "id": "d1",
                "ports": [
                    {
                        "id": "intf1",
                        "pmd": "PMD_UNSPECIFIED",
                        "speed": "S_400GB"
                    },
                    {
                        "id": "intf2",
                        "pmd": "PMD_UNSPECIFIED",
                        "speed": "S_400GB"
                    }
                ],
                "role": "DUT"
            },
            {
                "id": "d2",
                "ports": [
                    {
                        "id": "intf1",
                        "pmd": "PMD_UNSPECIFIED",
                        "speed": "S_400GB"
                    },
                    {
                        "id": "intf2",
                        "pmd": "PMD_UNSPECIFIED",
                        "speed": "S_400GB"
                    }
                ],
                "role": "ATE"
            }
        ],
        "links": [
            {
                "src": {
                    "device": "d1",
                    "port": "intf1"
                },
                "dst": {
                    "device": "d2",
                    "port": "intf1"
                }
            },
            {
                "src": {
                    "device": "d1",
                    "port": "intf2"
                },
                "dst": {
                    "device": "d2",
                    "port": "intf2"
                }
            }            
        ]
    }
    ```
## Running the Client for Cafy
* Pip install opentestbed
* Import opentestbed package
* Run the python file and the file has the below content.
* For the test run use the dynamically generated testbed_output.
    ```
    import opentestbed
    api = opentestbed.api(location="http://127.0.0.1:8080", transport="http")    
    testbed = opentestbed.Testbed()
    d1, d2 = testbed.devices.add(), testbed.devices.add()
    
    d1.id = "d1"
    # d1.name = "R2"
    d1.role = "DUT"
    d2.id = "d2"
    # d2.name = "TGEN2"
    d2.role = "ATE"
    
    d1_port1 = d1.ports.add()
    d1_port1.id = "intf1"
    d1_port1.speed = d1_port1.S_100GB
    
    d2_port1 = d2.ports.add()
    d2_port1.id = "intf1"
    d2_port1.speed = d2_port1.S_100GB    
    link1 = testbed.links.add()    
    link1.src.device = d1.id
    link1.src.port = d1_port1.id
    link1.dst.device = d2.id
    link1.dst.port = d2_port1.id
       
    testbed_output = api.reserve(testbed)
    print(testbed_output)
    ```
## Lab Reservation Service Setup
* Run reservation service (server) using docker run.
* Pull the latest version from the ghrc.
    ```docker pull ghcr.io/open-traffic-generator/lab-reservation-service:0.0.2​```
* Use the below command to run the server and ensure the Netbox is available.
    ```
    docker run -d -p 8080:8080 --name laas -e VERSION=0.0.2 lab_reservation_service:0.0.2 -netbox-host "netbox-host/IP" -netbox-port "netbox-port" -netbox-user-token "netbox-token" -framework-name cafy (generic/cafy/ondatra)
    ```
* Execute the client-side app (the above ondatra/cafy) to obtain the testbed reservation once the server is up and running.
## Setup Netbox Docker.
* Clone the Netbox repository.
    ```git clone -b release https://github.com/netbox-community/netbox-docker.git```
* Move to netbox-docker directory.
* Create a new file that defines the port under which NetBox will be available. The file name must be docker-compose.override.yml and its content should be as follows:
    ```
    version: '3.4'
    services:
      netbox:
        ports:
        - 8000:8080
    ```
* Need to pull all the containers from the Docker registry and may take a while, depending on your internet connection.
    ```docker compose pull```
* Finally, start all the required Docker containers.
    ```docker compose up```
* To create the first admin user run this command:
    ```docker compose exec netbox /opt/netbox/netbox/manage.py createsuperuser```
* If you want to just stop NetBox to continue work later on, use the following command.
    ```
    # Stop all the containers
    docker compose stop
    
    # Start the containers again
    docker compose start
    ```
## Workflow Diagram.
    ![image](https://github.com/keysightgems/Lab_service/assets/40664949/4296659d-d9d7-4bd9-b540-32ee5b45ddfc)
