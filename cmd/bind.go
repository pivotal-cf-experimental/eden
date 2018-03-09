package cmd

import (
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/pborman/uuid"
	"github.com/pivotal-cf-experimental/eden/apiclient"
)

// BindOpts represents the 'bind' command
type BindOpts struct {
}

// Execute is callback from go-flags.Commander interface
func (c BindOpts) Execute(_ []string) (err error) {
	instanceNameOrID := Opts.Instance.NameOrID
	if instanceNameOrID == "" {
		return fmt.Errorf("bind command requires --instance [NAME|GUID], or $SB_INSTANCE")
	}

	broker := apiclient.NewOpenServiceBroker(
		Opts.Broker.URLOpt,
		Opts.Broker.ClientOpt,
		Opts.Broker.ClientSecretOpt,
		Opts.Broker.APIVersion,
	)

	instance, err := broker.GetServiceInstance(instanceNameOrID)
	if err != nil {
		return errwrap.Wrapf("Failed to get service instance {{err}}", err)
	}

	bindingID := uuid.New()
	bindingName := fmt.Sprintf("%s-%s", instance.ServiceName, bindingID)

	// TODO - store allocated bindingIDs into local DB
	bindingResp, err := broker.Bind(instance.ServiceID, instance.PlanID, instance.ID, bindingID)
	if err != nil {
		return errwrap.Wrapf("Failed to bind to service instance {{err}}", err)
	}

	fmt.Println("Success")
	fmt.Println("")
	fmt.Printf("Created binding:\n\n%+v", bindingResp)
	fmt.Println("")
	fmt.Printf("Run 'eden credentials -i %s -b %s' to see credentials\n", instance.Name, bindingName)
	return
}
