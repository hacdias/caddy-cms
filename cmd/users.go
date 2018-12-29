package cmd

import (
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/filebrowser/filebrowser/types"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	rootCmd.AddCommand(usersCmd)
}

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Users management utility",
	Long:  `Users management utility.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(0)
	},
}

func printUsers(users []*types.User) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tUsername\tScope\tLocale\tV. Mode\tAdmin\tExecute\tCreate\tRename\tModify\tDelete\tShare\tDownload\tPwd Lock")

	for _, user := range users {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%t\t%t\t%t\t%t\t%t\t%t\t%t\t%t\t%t\t\n",
			user.ID,
			user.Username,
			user.Scope,
			user.Locale,
			user.ViewMode,
			user.Perm.Admin,
			user.Perm.Execute,
			user.Perm.Create,
			user.Perm.Rename,
			user.Perm.Modify,
			user.Perm.Delete,
			user.Perm.Share,
			user.Perm.Download,
			user.LockPassword,
		)
	}

	w.Flush()
}

func usernameOrIDRequired(cmd *cobra.Command, args []string) error {
	username, _ := cmd.Flags().GetString("username")
	id, _ := cmd.Flags().GetUint("id")

	if username == "" && id == 0 {
		return errors.New("'username' of 'id' flag required")
	}

	return nil
}

func addUserFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("perm.admin", false, "admin perm for users")
	cmd.Flags().Bool("perm.execute", true, "execute perm for users")
	cmd.Flags().StringSlice("perm.commands", nil, "a list of the commands a user can execute")
	cmd.Flags().Bool("perm.create", true, "create perm for users")
	cmd.Flags().Bool("perm.rename", true, "rename perm for users")
	cmd.Flags().Bool("perm.modify", true, "modify perm for users")
	cmd.Flags().Bool("perm.delete", true, "delete perm for users")
	cmd.Flags().Bool("perm.share", true, "share perm for users")
	cmd.Flags().Bool("perm.download", true, "download perm for users")
	cmd.Flags().Bool("lockPassword", false, "lock password")
	cmd.Flags().String("scope", "", "scope for users")
	cmd.Flags().String("locale", "en", "locale for users")
	cmd.Flags().String("viewMode", string(types.ListViewMode), "view mode for users")
}

func getViewMode(cmd *cobra.Command) types.ViewMode {
	viewMode := types.ViewMode(mustGetString(cmd, "viewMode"))
	if viewMode != types.ListViewMode && viewMode != types.MosaicViewMode {
		checkErr(errors.New("view mode must be \"" + string(types.ListViewMode) + "\" or \"" + string(types.MosaicViewMode) + "\""))
	}
	return viewMode
}

func getUserDefaults(cmd *cobra.Command, defaults *types.UserDefaults, all bool) {
	visit := func(flag *pflag.Flag) {
		switch flag.Name {
		case "scope":
			defaults.Scope = mustGetString(cmd, "scope")
		case "locale":
			defaults.Locale = mustGetString(cmd, "locale")
		case "viewMode":
			defaults.ViewMode = getViewMode(cmd)
		case "perm.admin":
			defaults.Perm.Admin = mustGetBool(cmd, "perm.admin")
		case "perm.execute":
			defaults.Perm.Execute = mustGetBool(cmd, "perm.execute")
		case "perm.create":
			defaults.Perm.Create = mustGetBool(cmd, "perm.create")
		case "perm.rename":
			defaults.Perm.Rename = mustGetBool(cmd, "perm.rename")
		case "perm.modify":
			defaults.Perm.Modify = mustGetBool(cmd, "perm.modify")
		case "perm.delete":
			defaults.Perm.Delete = mustGetBool(cmd, "perm.delete")
		case "perm.share":
			defaults.Perm.Share = mustGetBool(cmd, "perm.share")
		case "perm.download":
			defaults.Perm.Download = mustGetBool(cmd, "perm.download")
		case "perm.commands":
			commands, err := cmd.Flags().GetStringSlice("perm.commands")
			checkErr(err)
			defaults.Perm.Commands = commands
		}
	}

	if all {
		cmd.Flags().VisitAll(visit)
	} else {
		cmd.Flags().Visit(visit)
	}
}
