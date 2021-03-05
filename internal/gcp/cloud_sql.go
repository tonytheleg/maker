package gcp

import (
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

// CreateSQLService creates a new client to interact with GCP
func CreateSQLService(keyfile string) (*sqladmin.Service, error) {
	ctx := context.Background()
	sqlService, err := sqladmin.NewService(
		ctx, option.WithCredentialsFile(keyfile))
	if err != nil {
		return nil, errors.Errorf("Failed to create client: ", err)
	}
	return sqlService, nil
}

// CreateSQLInstance creates a compute instance with provided specs
func CreateSQLInstance(sqlService *sqladmin.Service, name, project, zone, machineType string) error {
	ctx := context.Background()

	db := &sqladmin.DatabaseInstance{
		ConnectionName:  name,
		DatabaseVersion: "POSTGRES_12",
		GceZone:         zone,
		InstanceType:    "CLOUD_SQL_INSTANCE",
		Name:            name,
		Project:         project,
		RootPassword:    "cloudsqltemp",
		Settings:        &sqladmin.Settings{Tier: machineType},
	}

	_, err := sqlService.Instances.Insert(project, db).Context(ctx).Do()
	if err != nil {
		return errors.Errorf("Failed to create SQL Instance: ", err)
	}
	fmt.Printf("SQL Instance %s is being created\n", name)
	return nil
}

// PrintSQLDbStatus outputs instance info
func PrintSQLDbStatus(sqlService *sqladmin.Service, name, project, zone string) error {
	ctx := context.Background()

	resp, err := sqlService.Instances.Get(project, name).Context(ctx).Do()
	if err != nil {
		return errors.Errorf("Failed to retreive GCE Instance %s: ", name, err)
	}
	fmt.Printf(
		"Name: %s\nConnection Name: %s\nDB Version: %s\n\nMaster Name: %s\nInstance Type: %s\nTier: %s\n\nIP Address: %s\nProject: %s\nRegion: %s\nZone: %s\nState: %s\n",
		resp.Name,
		resp.ConnectionName,
		resp.DatabaseVersion,
		resp.MasterInstanceName,
		resp.InstanceType,
		resp.Settings.Tier,
		resp.IpAddresses[0].IpAddress,
		resp.Project,
		resp.Region,
		resp.GceZone,
		resp.State,
	)
	return nil
}

// DeleteSQLInstance delets a droplet with the provided ID
func DeleteSQLInstance(sqlService *sqladmin.Service, name, project, zone string) error {
	ctx := context.Background()

	_, err := sqlService.Instances.Delete(project, name).Context(ctx).Do()
	if err != nil {
		return errors.Errorf("Failed to delete SQL Instance %s: ", name, err)
	}
	fmt.Printf("SQL Instance %s has been deleted\n", name)
	return nil
}
