package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/pivotal-cf/brokerapi"
	"github.com/pivotal-cf-experimental/eden/apiclient"
)

// ProvisionOpts represents the 'provision' command
type ProvisionOpts struct {
	ServiceNameOrID string `short:"s" long:"service-name" description:"Service name/ID from catalog" required:"true"`
	PlanNameOrID    string `short:"p" long:"plan-name" description:"Plan name/ID from catalog (default: first)"`
}

// Execute is callback from go-flags.Commander interface
func (c ProvisionOpts) Execute(_ []string) (err error) {
	broker := apiclient.NewOpenServiceBroker(
		Opts.Broker.URLOpt,
		Opts.Broker.ClientOpt,
		Opts.Broker.ClientSecretOpt,
		Opts.Broker.APIVersion,
	)

	service, err := broker.FindServiceByNameOrID(c.ServiceNameOrID)
	if err != nil {
		return errwrap.Wrapf("Could not find service in catalog: {{err}}", err)
	}
	plan, err := broker.FindPlanByNameOrID(service, c.PlanNameOrID)
	if err != nil {
		return errwrap.Wrapf("Could not find plan in service: {{err}}", err)
	}

	instanceName := Opts.Instance.NameOrID

	prexisting := Opts.config().FindServiceInstance(instanceName)
	if prexisting.ServiceName != "" {
		return fmt.Errorf("Service instance '%s' already exists", instanceName)
	}

	provisioningResp, isAsync, err := broker.Provision(service.ID, plan.ID, instanceName, instanceName)
	if err != nil {
		return errwrap.Wrapf("Failed to provision service instance: {{err}}", err)
	}

	fmt.Printf("provision:   %s/%s - name: %s\n", service.Name, plan.Name, instanceName)
	if isAsync {
		fmt.Println("provision:   in-progress")
		// TODO: don't pollute brokerapi back into this level
		lastOpResp := &brokerapi.LastOperationResponse{State: brokerapi.InProgress}
		for lastOpResp.State == brokerapi.InProgress {
			time.Sleep(5 * time.Second)
			lastOpResp, err = broker.LastOperation(service.ID, plan.ID, instanceName, provisioningResp.OperationData)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}
			fmt.Printf("provision:   %s - %s\n", lastOpResp.State, lastOpResp.Description)
		}
	}
	if provisioningResp.DashboardURL == "" {
		fmt.Println("provision:   done")
	} else {
		fmt.Printf("provision:   done - %s\n", provisioningResp.DashboardURL)
	}

	return
}
