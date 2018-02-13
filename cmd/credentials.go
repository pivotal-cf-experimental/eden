package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/starkandwayne/eden/apiclient"
)

// CredentialsOpts represents the 'credentials' command
type CredentialsOpts struct {
	BindingID string `short:"b" long:"bind" description:"Binding to display"`
	Attribute string `short:"a" long:"attribute" description:"Only diplay a single attribute from credentials"`
}

// Execute is callback from go-flags.Commander interface
func (c CredentialsOpts) Execute(_ []string) (err error) {
	instanceNameOrID := Opts.Instance.NameOrID
	if instanceNameOrID == "" {
		return fmt.Errorf("credentials command requires --instance [NAME|GUID], or $SB_INSTANCE")
	}

	broker := apiclient.NewOpenServiceBroker(
		Opts.Broker.URLOpt,
		Opts.Broker.ClientOpt,
		Opts.Broker.ClientSecretOpt,
		Opts.Broker.APIVersion,
	)

	bindings, err := broker.GetBindingsByServiceInstanceID(instanceNameOrID)

	fmt.Println("")
	fmt.Println(bindings)
	fmt.Println("")
	return
}

func (c CredentialsOpts) displayBinding(credentials map[string]interface{}, attribute string) error {
	if attribute == "" {
		b, err := json.MarshalIndent(credentials, "", "  ")
		if err != nil {
			return errwrap.Wrapf("Could not marshal credentials: {{err}}", err)
		}
		os.Stdout.Write(b)
		return nil
	}
	if val, ok := credentials[attribute]; ok {
		fmt.Printf("%v\n", val)
		return nil
	}
	attributes := make([]string, 0, len(credentials))
	for key := range credentials {
		attributes = append(attributes, key)
	}

	return fmt.Errorf("credentials --attribute key was unknown; try: %s", strings.Join(attributes, ", "))
}
