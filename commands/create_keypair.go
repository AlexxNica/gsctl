package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/giantswarm/gsclientgen"
	"github.com/giantswarm/microerror"
	"github.com/spf13/cobra"

	"github.com/giantswarm/gsctl/client"
	"github.com/giantswarm/gsctl/config"
	"github.com/giantswarm/gsctl/util"
)

var (
	// CreateKeypairCommand performs the "create keypair" function
	CreateKeypairCommand = &cobra.Command{
		Use:     "keypair",
		Short:   "Create key pair",
		Long:    `Creates a new key pair for a cluster`,
		PreRunE: checkAddKeypair,
		Run:     addKeypair,
	}
)

const (
	addKeyPairActivityName = "add-keypair"
)

func init() {
	CreateKeypairCommand.Flags().StringVarP(&cmdClusterID, "cluster", "c", "", "ID of the cluster to create a key pair for")
	CreateKeypairCommand.Flags().StringVarP(&cmdDescription, "description", "d", "", "Description for the key pair")
	CreateKeypairCommand.Flags().StringVarP(&cmdCNPrefix, "cn-prefix", "", "", "The common name prefix for the issued certificates 'CN' field.")
	CreateKeypairCommand.Flags().StringVarP(&cmdCertificateOrganizations, "certificate-organizations", "", "", "A comma separated list of organizations for the issued certificates 'O' fields.")
	CreateKeypairCommand.Flags().IntVarP(&cmdTTLDays, "ttl", "", 30, "Duration until expiry of the created key pair in days")

	CreateKeypairCommand.MarkFlagRequired("cluster")

	CreateCommand.AddCommand(CreateKeypairCommand)
}

func checkAddKeypair(cmd *cobra.Command, args []string) error {

	endpoint := config.Config.ChooseEndpoint(cmdAPIEndpoint)
	token := config.Config.ChooseToken(endpoint, cmdToken)

	if endpoint == "" {
		return microerror.Mask(endpointMissingError)
	}
	if token == "" {
		return errors.New("You are not logged in. Use '" + config.ProgramName + " login' to log in.")
	}
	if cmdClusterID == "" {
		// use default cluster if possible
		clusterID, _ := config.GetDefaultCluster(requestIDHeader, addKeyPairActivityName, cmdLine, cmdAPIEndpoint)
		if clusterID != "" {
			cmdClusterID = clusterID
		} else {
			return errors.New("No cluster given. Please use the -c/--cluster flag to set a cluster ID.")
		}
	}
	if cmdDescription == "" {
		return errors.New("No description given. Please use the -d/--description flag to set a description.")
	}
	return nil
}

func addKeypair(cmd *cobra.Command, args []string) {
	if cmdDescription == "" {
		cmdDescription = "Added by user " + config.Config.Email + " using 'gsctl create keypair'"
	}

	endpoint := config.Config.ChooseEndpoint(cmdAPIEndpoint)
	token := config.Config.ChooseToken(endpoint, cmdToken)

	clientConfig := client.Configuration{
		Endpoint:  endpoint,
		UserAgent: config.UserAgent(),
	}
	apiClient, clientErr := client.NewClient(clientConfig)
	if clientErr != nil {
		fmt.Println(color.RedString("Error: %s", clientErr))
		os.Exit(1)
	}

	authHeader := "giantswarm " + token
	ttlHours := int32(cmdTTLDays * 24)
	addKeyPairBody := gsclientgen.V4AddKeyPairBody{Description: cmdDescription, TtlHours: ttlHours, CnPrefix: cmdCNPrefix, CertificateOrganizations: cmdCertificateOrganizations}
	keypairResponse, apiResponse, err := apiClient.AddKeyPair(authHeader, cmdClusterID, addKeyPairBody, requestIDHeader, addKeyPairActivityName, cmdLine)

	if err != nil {
		fmt.Println(color.RedString("Error: %s", err))
		dumpAPIResponse(*apiResponse)
		os.Exit(1)
	}

	if apiResponse.StatusCode == 200 {
		msg := fmt.Sprintf("New key pair created with ID %s and expiry of %v",
			util.Truncate(util.CleanKeypairID(keypairResponse.Id), 10),
			util.DurationPhrase(int(keypairResponse.TtlHours)))
		fmt.Println(color.GreenString(msg))

		// store credentials to file
		caCertPath := util.StoreCaCertificate(config.CertsDirPath, cmdClusterID, keypairResponse.CertificateAuthorityData)
		fmt.Println("CA certificate stored in:", caCertPath)

		clientCertPath := util.StoreClientCertificate(config.CertsDirPath, cmdClusterID, keypairResponse.Id, keypairResponse.ClientCertificateData)
		fmt.Println("Client certificate stored in:", clientCertPath)

		clientKeyPath := util.StoreClientKey(config.CertsDirPath, cmdClusterID, keypairResponse.Id, keypairResponse.ClientKeyData)
		fmt.Println("Client private key stored in:", clientKeyPath)

	} else {
		fmt.Println(color.RedString("Unhandled response code: %v", apiResponse.StatusCode))
		dumpAPIResponse(*apiResponse)
	}
}
