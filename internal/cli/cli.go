package cli

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type AzureResourceVerifierCliError struct {
	Message string
	Err     error
}

func (e *AzureResourceVerifierCliError) Error() string {
	return e.Err.Error()
}

func AzureClientWrapRunE(
	runEFunc func(cmd *cobra.Command, args []string, cred *azidentity.DefaultAzureCredential, ctx context.Context) error,
) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// Silence usage so we don't print the usage when an error occurs
		cmd.SilenceUsage = true

		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return fmt.Errorf("error binding flags: %s", err)
		}

		cred, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			return err
		}

		ctx := context.Background()

		return runEFunc(cmd, args, cred, ctx)
	}
}

func CreateAzrErr(msg string, err error) error {
	return &AzureResourceVerifierCliError{Message: msg, Err: err}
}

func ExitOnError(err error, userMessage string) {
	var message string
	var details string
	exitCode := 1 // Default to 1

	if err != nil {
		// Print the user message first if provided
		if userMessage != "" {
			fmt.Fprintf(os.Stderr, "Message: %s\n", userMessage)
		}

		// Check if the error is wrapped
		var wrappedError *AzureResourceVerifierCliError
		if errors.As(err, &wrappedError) {
			err = wrappedError.Err
			message = wrappedError.Message
		}

		// Check if it's an Azure Authentication Error
		if azureErr, ok := err.(*azidentity.AuthenticationFailedError); ok {
			message = "It looks like you're not authenticated. Please run `az login` and try again."
			details = azureErr.Error()
		} else if azureErr, ok := err.(*azidentity.AuthenticationRequiredError); ok {
			message = "It looks like you're not authenticated. Please run `az login` and try again."
			details = azureErr.Error()

			/* TODO: credentialUnavailableError is not available in the current version of the SDK
			} else if azureErr, ok := err.(*azidentity.credentialUnavailableError); ok {
				message = "It looks like you're not authenticated. Please run `az login` and try again."
				details = azureErr.Error()
			*/
		} else {
			details = err.Error()
		}

		// Print the message, if any
		if message != "" {
			fmt.Fprintf(os.Stderr, "Message: %s\n", message)
		}
		// Print the details, if any
		if details != "" {
			fmt.Fprintf(os.Stderr, "Details: %s\n", details)
		}
		os.Exit(exitCode)
	}
}
