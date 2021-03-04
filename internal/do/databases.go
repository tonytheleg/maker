package do

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/pkg/errors"
)

// CreateDoDatabase creates a Postgres DB cluster on Digital Ocean
func CreateDoDatabase(client *godo.Client, name, size, region string) error {
	ctx := context.TODO()
	createRequest := &godo.DatabaseCreateRequest{
		Name:       name,
		EngineSlug: "pg",
		Version:    "10",
		Region:     region,
		SizeSlug:   size,
		NumNodes:   1,
	}

	cluster, _, err := client.Databases.Create(ctx, createRequest)
	if err != nil {
		return errors.Errorf("Failed to create database:", err)
	}
	fmt.Println(cluster.Name, "created")
	return nil
}

// GetDoDatabase grabs the database ID with the provided name
func GetDoDatabase(client *godo.Client, name string) (string, error) {
	var databaseID string
	ctx := context.TODO()
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	databases, _, err := client.Databases.List(ctx, opt)
	if err != nil {
		return "", errors.Wrapf(err, "Could not list databases to search for %s:", name)
	}
	for index := range databases {
		if databases[index].Name == name {
			databaseID = databases[index].ID
		}
	}
	if databaseID != "" {
		return databaseID, nil
	}
	return "", errors.Wrapf(err, "Could not find database with name %s:", name)

}

// PrintDatabaseStatus outputs some database info
func PrintDatabaseStatus(client *godo.Client, id string) {
	ctx := context.TODO()
	database, _, err := client.Databases.Get(ctx, id)
	if err != nil {
		fmt.Println("Could not fetch database status:", err)
	}
	fmt.Printf(
		"Name: %s\nUID: %s\nEngine: %s\n\nConnection URI: %s\nHost: %s\nPort: %d\n\nUsername: %s\nPassword: %s\n\nNumber of Nodes: %d\nNode Size: %s\n\nRegion: %s\nCreated: %v\nStatus: %s\n",
		database.Name,
		database.ID,
		database.EngineSlug,
		database.Connection.URI,
		database.Connection.Host,
		database.Connection.Port,
		database.Connection.User,
		database.Connection.Password,
		database.NumNodes,
		database.SizeSlug,
		database.RegionSlug,
		database.CreatedAt,
		database.Status,
	)
}

// DeleteDoDatabase delets a database with the provided ID
func DeleteDoDatabase(client *godo.Client, id string, name string) error {
	ctx := context.TODO()
	_, err := client.Databases.Delete(ctx, id)
	if err != nil {
		return errors.Errorf("Deleting database failed:", err)
	}
	fmt.Println("Database", name, "deleted")
	return nil
}
