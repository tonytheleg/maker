package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/pkg/errors"
)

// CreateRdsInstance creates a Postgres RDS instance in AWS
func CreateRdsInstance(sess *session.Session, name, size string) error {
	svc := rds.New(sess)
	input := &rds.CreateDBInstanceInput{
		AllocatedStorage:     aws.Int64(5),
		DBInstanceClass:      aws.String(size),
		DBInstanceIdentifier: aws.String(name),
		Engine:               aws.String("postgres"),
		MasterUserPassword:   aws.String("rdsadmin"),
		MasterUsername:       aws.String("rdsadmintemp"),
	}

	_, err := svc.CreateDBInstance(input)
	if err != nil {
		return errors.Errorf("Failed to create database %s:", name, err)
	}
	fmt.Println("Database", name, "creating")
	return nil
}

// DeleteRdsInstance deletes a Postgres RDS instance in AWS
func DeleteRdsInstance(sess *session.Session, name string) error {
	svc := rds.New(sess)
	input := &rds.DeleteDBInstanceInput{
		DBInstanceIdentifier: aws.String(name),
		SkipFinalSnapshot:    aws.Bool(true),
	}

	_, err := svc.DeleteDBInstance(input)
	if err != nil {
		return errors.Errorf("Failed to create database %s:", name, err)
	}
	fmt.Println("Database", name, "is being deleted")
	return nil
}

// PrintRdsStatus prints the status of a RDS DB instance
func PrintRdsStatus(sess *session.Session, name string) error {
	svc := rds.New(sess)
	input := &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(name),
	}

	result, err := svc.DescribeDBInstances(input)
	if err != nil {
		return errors.Errorf("Failed to create database %s:", name, err)
	}
	fmt.Printf("Name: %s\nARN: %s\nAZ: %s\nSize: %s\n\nDB Username: %s\nEndpoint: %s\nDB Engine: %s\nDB Version: %s\n\nCreated: %s\nStatus: %s\n",
		*result.DBInstances[0].DBInstanceIdentifier,
		*result.DBInstances[0].DBInstanceArn,
		*result.DBInstances[0].AvailabilityZone,
		*result.DBInstances[0].DBInstanceClass,
		*result.DBInstances[0].MasterUsername,
		*result.DBInstances[0].Endpoint.Address,
		*result.DBInstances[0].Engine,
		*result.DBInstances[0].EngineVersion,
		*result.DBInstances[0].InstanceCreateTime,
		*result.DBInstances[0].DBInstanceStatus,
	)
	return nil
}
