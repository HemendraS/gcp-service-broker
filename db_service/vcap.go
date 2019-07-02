package db_service

import (
	"code.cloudfoundry.org/lager"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"net/url"
	"os"
)

type VcapService struct {
	BindingName  string            `json:"binding_name"`  // The name assigned to the service binding by the user.
	InstanceName string            `json:"instance_name"` // The name assigned to the service instance by the user.
	Name         string            `json:"name"`          // The binding_name if it exists; otherwise the instance_name.
	Label        string            `json:"label"`         // The name of the service offering.
	Tags         []string          `json:"tags"`          // An array of strings an app can use to identify a service instance.
	Plan         string            `json:"plan"`          // The service plan selected when the service instance was created.
	Credentials  map[string]string `json:"credentials"`   // The service-specific credentials needed to access the service instance.
}

func useVcapServices(logger lager.Logger) error {
	vcapData, vcapExists := os.LookupEnv("VCAP_SERVICES")
	if vcapExists {
		vcapService, err := parseVcapServices(vcapData, logger)
		if err != nil {
			return fmt.Errorf("Error parsing VCAP_SERVICES: %s", err)
		}

		// if URI is supplied, we should parse it to fill any missing fields
		u, err := url.Parse(vcapService.Credentials["uri"])
		if err != nil {
			return fmt.Errorf("Error parsing VCAP_SERVICES credentials URI: %s", err)
		}

		logger.Info("Using MySQL database injected via VCAP_SERVICES environment variable")
		viper.Set(dbPathProp, u.Path)
		viper.Set(dbTypeProp, DbTypeMysql)
		viper.Set(dbHostProp, vcapService.Credentials["host"])
		viper.Set(dbUserProp, vcapService.Credentials["Username"])
		viper.Set(dbPassProp, vcapService.Credentials["Password"])
		viper.Set(dbNameProp, vcapService.Credentials["database_name"])

		//  if database is one provided by gcp service broker, use the client_cert, ca_cert and client_key fields
		if contains(vcapService.Tags, "gcp") {
			viper.Set(caCertProp, vcapService.Credentials["CaCert"])
			viper.Set(clientCertProp, vcapService.Credentials["ClientCert"])
			viper.Set(clientKeyProp, vcapService.Credentials["ClientKey"])
		}
	}
	return nil
}

func parseVcapServices(vcapServicesEnv string, logger lager.Logger) (VcapService, error) {
	var vcapServiceMap map[string]*json.RawMessage
	err := json.Unmarshal([]byte(vcapServicesEnv), &vcapServiceMap)
	if err != nil {
		return VcapService{}, fmt.Errorf("Error unmarshaling VCAP_SERVICES: %s", err)
	}
	var vcapServices []VcapService
	for _,v := range vcapServiceMap {
		err := json.Unmarshal(*v, &vcapServices)
		if err != nil {
			return VcapService{}, fmt.Errorf("Error unmarshaling VCAP_SERVICES: %s", err)
		}
	}
	index, err := findMySqlTag(vcapServices, "mysql")
	if err != nil {
		return VcapService{}, fmt.Errorf("Error finding MySQL tag: %s", err)
	}
	return vcapServices[index], nil
}

// whether a given string array arr contains string key
func contains(arr []string, key string) bool {
	for _, n := range arr {
		if key == n {
			return true
		}
	}
	return false
}

// return the index of the VcapService with a tag of "mysql" in the list of VcapServices, fail if we find more or fewer than 1
func findMySqlTag(VcapServices []VcapService, key string) (int, error) {
	index := -1
	count := 0
	for i, vcapService := range VcapServices {
		if contains(vcapService.Tags, key) {
			count += 1
			index = i
		}
	}
	if count != 1 {
		return -1, fmt.Errorf("The variable VCAP_SERVICES must have one VCAP service with a tag of %s. There are currently %d VCAP services with the tag %s.", "'mysql'", count, "'mysql'")
	}
	return index, nil
}