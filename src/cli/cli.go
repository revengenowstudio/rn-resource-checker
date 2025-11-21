package cli

import (
	"context"

	"github.com/urfave/cli/v3"

	. "rn-resource-checker/src"
	"rn-resource-checker/src/log"
)

type ArgTypes string

const (
	Args_Path   string = "path"
	Args_Suffix string = "suffix"
)

func Run(argv []string) (string, error) {
	cmd := &cli.Command{
		Name:  "rn-resource-checker",
		Usage: "input path list and suffix list to specify the path and suffix to find , or input nothing to use default settings",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{Name: Args_Path, Aliases: []string{"p"}, Value: []string{"./"}},
			&cli.StringSliceFlag{Name: Args_Suffix, Aliases: []string{"s"}, Value: []string{"mix", "exe", "dll", "ext"}},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			targetDirs := cmd.StringSlice(Args_Path)
			whiteListSuffix := cmd.StringSlice(Args_Suffix)
			err := DoHashJob(targetDirs, whiteListSuffix)
			if err != nil {
				return err
			}
			return nil
		},
	}

	if err := cmd.Run(context.Background(), argv); err != nil {
		log.Fatal(err)
		return "", err
	}
	return "", nil
}
