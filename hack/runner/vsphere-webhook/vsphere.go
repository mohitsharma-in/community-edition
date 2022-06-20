package main

import (
	"context"
	"fmt"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/types"
	klog "k8s.io/klog/v2"
	"net/url"
	"os"
)

const (
	envURL          = "GOVC_URL"
	envUserName     = "GOVC_USERNAME"
	envPassword     = "GOVC_PASSWORD"
	envDataStore    = "GOVC_DATASTORE"
	envFolder       = "GOVC_FOLDER"
	envNetwork      = "GOVC_NETWORK"
	envResourcePool = "GOVC_RESOURCE_POOL"
)

var ctx context.Context = context.Background()

func GetGovmomiClient() *vim25.Client {
	//TODO: To make use of common creds function or struct to avoid redundant vars
	envUserName := os.Getenv(envUserName)
	envPassword := os.Getenv(envPassword)
	envURL := os.Getenv(envURL)
	u := &url.URL{
		Scheme: "https",
		Host:   envURL,
		Path:   "/sdk",
	}
	u.User = url.UserPassword(envUserName, envPassword)
	client, err := govmomi.NewClient(ctx, u, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Login to vsphere failed, %v", err)
		os.Exit(1)
	}
	return client.Client
}

func DeleteVM(client *vim25.Client, vmName string) error {
	finder := find.NewFinder(client)
	vm, err := finder.VirtualMachine(context.TODO(), vmName)
	if err != nil {
		if _, ok := err.(*find.NotFoundError); ok {
			klog.Errorf("Unable To find VM")
			return err
		}
	}
	var (
		task  *object.Task
		state types.VirtualMachinePowerState
	)

	state, err = vm.PowerState(ctx)
	if err != nil {
		return err
	}

	if state == types.VirtualMachinePowerStatePoweredOn {
		task, err = vm.PowerOff(ctx)
		if err != nil {
			return err
		}
		// Ignore error since the VM may already been in powered off state.
		// vm.Destroy will fail if the VM is still powered on.
		_ = task.Wait(ctx)
	}

	task, err = vm.Destroy(ctx)
	if err != nil {
		return err
	}

	if err = task.Wait(ctx); err != nil {
		return err
	}
	return nil
}

func CreateVMGithubRunner(client *vim25.Client, uniqueID, runnerToken string) (string, error) {
	klog.V(6).Infof("uniqueID: %s\n", uniqueID)
	klog.V(6).Infof("runnerToken: %s\n", runnerToken)
	finder := find.NewFinder(client)
	envResourcePool := os.Getenv(envResourcePool)
	envFolder := os.Getenv(envFolder)
	envNetwork := os.Getenv(envNetwork)
	envDataStore := os.Getenv(envDataStore)

	return fmt.Sprintf("nil"), nil
}
