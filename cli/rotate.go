package cli

import (
	"fmt"

	"github.com/99designs/aws-vault/prompt"
	"github.com/99designs/aws-vault/vault"
	"github.com/99designs/keyring"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"gopkg.in/alecthomas/kingpin.v2"
)

type RotateCommandInput struct {
	Profile   string
	Keyring   keyring.Keyring
	MfaToken  string
	MfaPrompt prompt.PromptFunc
}

func ConfigureRotateCommand(app *kingpin.Application) {
	input := RotateCommandInput{}

	cmd := app.Command("rotate", "Rotates credentials")
	cmd.Arg("profile", "Name of the profile").
		Required().
		StringVar(&input.Profile)

	cmd.Flag("mfa-token", "The mfa token to use").
		Short('t').
		StringVar(&input.MfaToken)

	cmd.Action(func(c *kingpin.ParseContext) error {
		input.MfaPrompt = prompt.Method(GlobalFlags.PromptDriver)
		input.Keyring = keyringImpl
		RotateCommand(app, input)
		return nil
	})
}

func RotateCommand(app *kingpin.Application, input RotateCommandInput) {
	rotator := vault.Rotator{
		Keyring:   input.Keyring,
		MfaToken:  input.MfaToken,
		MfaPrompt: input.MfaPrompt,
		Config:    awsConfig,
	}

	fmt.Printf("Rotating credentials for profile %q (takes 10-20 seconds)\n", input.Profile)

	if err := rotator.Rotate(input.Profile); err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "AccessDenied" {
			app.Fatalf("AccessDenied: Check your credentials have permission to list and update access keys. See https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_examples_iam_credentials_console.html for more details.")
		} else {
			app.Fatalf(awsConfig.FormatCredentialError(err, input.Profile))
		}
		return
	}

	fmt.Printf("Done!\n")
}
