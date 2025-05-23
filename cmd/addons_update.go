package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var addonsUpdateCmd = &cobra.Command{
	Use:     "update [slug]",
	Aliases: []string{"upgrade", "up"},
	Short:   "Upgrades a Home Assistant add-on to the latest version",
	Long: `
This command can upgrade a Home Assistant add-on to its latest version.
It is currently not possible to upgrade/downgrade to a specific version.
`,
	Example: `
  ha addons update core_ssh
`,
	ValidArgsFunction: addonsCompletions,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("addons update")

		section := "addons"
		command := "{slug}/update"

		url, err := helper.URLHelper(section, command)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
			return
		}

		ProgressSpinner.Start()
		request := helper.GetJSONRequestTimeout(helper.ContainerDownloadTimeout)
		ProgressSpinner.Stop()

		slug := args[0]

		request.SetPathParams(map[string]string{
			"slug": slug,
		})

		options := make(map[string]any)

		backup, _ := cmd.Flags().GetBool("backup")
		if cmd.Flags().Changed("backup") {
			request.SetBody(options)
			options["backup"] = backup
		}

		if len(options) > 0 {
			log.WithField("options", options).Debug("Request body")
			request.SetBody(options)
		}

		resp, err := request.Post(url)
		resp, err = helper.GenericJSONErrorHandling(resp, err)

		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	addonsUpdateCmd.Flags().Bool("backup", false, "Create partial backup before update")
	addonsUpdateCmd.RegisterFlagCompletionFunc("backup", boolCompletions)
	addonsCmd.AddCommand(addonsUpdateCmd)
}
