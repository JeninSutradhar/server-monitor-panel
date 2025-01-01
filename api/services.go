package api

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// simple interface to work with generic operations.(if is used you also will call only on upperCamelCase as interface logic ). For now is only for documentation or use a template pattern
type Service interface {
	Start() error
	Stop() error
	Reload() error
	Uninstall() error
}

type ServiceStatus string // (Pascal case string to define at struct also, const, at public variables on code implementation (from files))
const (
	SERVICE_STATUS_STARTED      ServiceStatus = "STARTED"    // this defines the data variable if a string
	SERVICE_STATUS_STOPPED      ServiceStatus = "STOPPED"    // (all uppercase and separated to space) for exported variables with string type at local implementations also!. ( for all static const etc.. if export to be a valid public access code variables )
	SERVICE_STATUS_INSTALLING   ServiceStatus = "INSTALLING" //
	SERVICE_STATUS_UNINSTALLING ServiceStatus = "UNINSTALLING"
)

// general struct for services operation (exported types using Pascal Case). Also should to specify ` json tag`, or public struct will not return a valid  json value. (and are invisible to another structure using )
type ServiceInfo struct {
	Name   string        `json:"name"`
	Status ServiceStatus `json:"status"` // must keep the public declaration types also, since structs variables by other data source can call
	//we need to persist somewhere , but right now it would be volatile inside in mem structure of the struct
}
type ServiceStore struct {
	services map[string]*ServiceInfo

	sync.RWMutex
}

// map which holds running services and if running or not, volatile
var serviceList *ServiceStore

// initialize empty list of services in startup to use for volatile state tracking (all uppercase,  and also  export function name implementation also if data type exist to return from that variable , method/ or logic function call).
func InitServices() {

	serviceList = &ServiceStore{

		services: make(map[string]*ServiceInfo),
	}

}

// Simulates  instalation  and  create services that will persist volatile at all steps of application(  use pascal case for this specific implementation ) also as other public struct method ( export function by rules, uppercase to see those implementations methods with that scope, all those also, if need for use them external that specific files implementation.)
func Install(serviceName string) error {
	serviceList.Lock()
	defer serviceList.Unlock()

	_, ok := serviceList.services[serviceName]
	if ok {

		return fmt.Errorf("the service: %v exist and cannot re install it", serviceName)

	}

	serviceList.services[serviceName] = &ServiceInfo{
		Name: serviceName,

		Status: SERVICE_STATUS_INSTALLING,
	}

	go installService(serviceName)

	return nil
}

// simulates instalation after the create method was performed, and updates status accordingly to service (private func to avoid problems with export types using upper cases as parameters type on return) local implementation for that function usage  (as also structs data). So no type exports for those! and with all lower-cases names by convention for internal use . With a `time out`

func installService(serviceName string) {
	rand.Seed(rand.Int63())

	var sec int = rand.Intn(5) + 1
	fmt.Println("Installing service " + serviceName + "  takes  " + fmt.Sprintf("%v", sec))

	<-time.After(time.Duration(sec) * time.Second)

	serviceList.Lock()
	defer serviceList.Unlock()

	if service, exist := serviceList.services[serviceName]; exist {
		service.Status = SERVICE_STATUS_STOPPED
	}

}

func GetServiceStatus(name string) (ServiceInfo, error) {

	serviceList.RLock()

	defer serviceList.RUnlock()

	if service, ok := serviceList.services[name]; ok {

		return *service, nil

	}

	return ServiceInfo{}, fmt.Errorf("service %v   not found in the system", name)

}

// simulates running specific operations of service in a period.( public func or method using Camelcase!)

func Start(serviceName string) error {
	serviceList.Lock()

	defer serviceList.Unlock()

	if service, ok := serviceList.services[serviceName]; ok {
		if service.Status == SERVICE_STATUS_STARTED {

			return errors.New(fmt.Sprintf("The service: %v,  already started, and status = STARTED", serviceName))

		}
		go startService(service)

		return nil
	}

	return fmt.Errorf("service by that %s is not instaled!", serviceName)

}

// simulation function of doing an actual task and setting a timer to wait to respond.(all lower cases no exports local implementation types ). Since no return structs/ type .
func startService(service *ServiceInfo) {

	rand.Seed(rand.Int63())
	var sec = rand.Intn(3) + 1

	fmt.Printf("Starting   %s ,    waiting    %d   sec... \n", service.Name, sec)

	<-time.After(time.Duration(sec) * time.Second)

	serviceList.Lock()

	defer serviceList.Unlock()
	service.Status = SERVICE_STATUS_STARTED

}

// simulate to stop current services and update the volatile map info for every change, like setting it stopped in general ( all  methods for internal implementation with no exports, low cases also the data variables)
func Stop(serviceName string) error {
	serviceList.Lock()

	defer serviceList.Unlock()
	if service, ok := serviceList.services[serviceName]; ok {

		if service.Status == SERVICE_STATUS_STOPPED {

			return errors.New(fmt.Sprintf(" service:  %v   , already stop", service.Name))

		}

		go stopService(service)

		return nil
	}

	return fmt.Errorf("service :%v , not founded to stop the implementation for /services endpoint!", serviceName)

}

func stopService(service *ServiceInfo) { // if is local implementation

	rand.Seed(rand.Int63())

	var sec = rand.Intn(2) + 1

	fmt.Printf("Stopping  service :  %s   ,waiting :   %d  sec ...\n", service.Name, sec)

	<-time.After(time.Duration(sec) * time.Second)

	serviceList.Lock()

	defer serviceList.Unlock()
	service.Status = SERVICE_STATUS_STOPPED
}

// function that simulates service reloading for every specific call ( export implementation with pascal cases)
func Reload(serviceName string) error {

	serviceList.Lock()
	defer serviceList.Unlock()

	if service, ok := serviceList.services[serviceName]; ok {

		go reloadService(service)
		return nil

	}
	return fmt.Errorf("no services was not be relaod with name id: %s   , for endpoint of method '/api/services'. Please, input a valid value ", serviceName)

}
func reloadService(service *ServiceInfo) { //Local, use cases

	rand.Seed(rand.Int63())

	var sec = rand.Intn(5) + 1

	fmt.Printf("Reloading service  :  %s ,     waiting   %d   seconds  ... \n", service.Name, sec)

	<-time.After(time.Duration(sec) * time.Second)

	fmt.Println("Service :" + service.Name + "   reloaded  successful!")

}

// function to call an operation for uninstall operation for specified services.( public struct )
func Uninstall(serviceName string) error {

	serviceList.Lock()

	defer serviceList.Unlock()

	service, ok := serviceList.services[serviceName]

	if ok {
		service.Status = SERVICE_STATUS_UNINSTALLING

		go unInstallService(serviceName)

		return nil

	}
	return fmt.Errorf("service %s cannot be founded and not is available to removed !.", serviceName)
}
func unInstallService(name string) { //local  struct with impl
	rand.Seed(rand.Int63())
	var sec = rand.Intn(10) + 1

	fmt.Println("Unistalling " + name + ",  waiting for " + fmt.Sprintf("%v", sec) + "  sec...")
	<-time.After(time.Duration(sec) * time.Second)
	serviceList.Lock()

	defer serviceList.Unlock()
	delete(serviceList.services, name)

	fmt.Printf("Service   %s  ,  removed successful \n", name)

}

// get services available in current list in memory, that keeps track of install / uninstall status ( export public implementation by types)
func ListServices() map[string]*ServiceInfo {
	serviceList.RLock()

	defer serviceList.RUnlock()
	services := make(map[string]*ServiceInfo)
	for key, value := range serviceList.services {

		services[key] = value

	}
	return services

}
