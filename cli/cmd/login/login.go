package login

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/docker/api/cli/dockerclassic"
	"github.com/docker/api/client"
)

// Command returns the login command
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login [OPTIONS] [SERVER] | login azure",
		Short: "Log in to a Docker registry",
		Long:  "Log in to a Docker registry or cloud backend.\nIf no registry server is specified, the default is defined by the daemon.",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runLogin,
	}
	// define flags for backward compatibility with docker-classic
	flags := cmd.Flags()
	flags.StringP("username", "u", "", "Username")
	flags.StringP("password", "p", "", "Password")
	flags.BoolP("password-stdin", "", false, "Take the password from stdin")

	return cmd
}

func runLogin(cmd *cobra.Command, args []string) error {
	if len(args) == 1 && !strings.Contains(args[0], ".") {
		backend := args[0]
		switch backend {
		case "azure":
			return cloudLogin(cmd, "aci")
		default:
			return errors.New("Unknown backend type for cloud login : " + backend)
		}
	}
	return dockerclassic.ExecCmd(cmd)
}

func cloudLogin(cmd *cobra.Command, backendType string) error {
	ctx := cmd.Context()
	cs, err := client.GetCloudService(ctx, backendType)
	if err != nil {
		return errors.Wrap(err, "cannot connect to backend")
	}
	return cs.Login(ctx, nil)

}