package command

import (
	"strings"

	"fmt"
	"github.com/ikaven1024/bolt-cli/cli/framework"
)

type rootContext struct {
	*framework.BaseContext
}

func NewRootContext(fw framework.Framework) framework.Context {
	c := &rootContext{
		BaseContext: framework.NewBaseContext(fw),
	}
	c.init()
	return c
}

func (c *rootContext) Exit() {
}

func (c *rootContext) init() {
	c.AddCommand(
		framework.Command{
			Name:        "create",
			Description: "create a bucket",
			Usage:       "<bucketname>",
			Args:        framework.ExactArgs(1),
			Action:      c.create,
		},
		framework.Command{
			Name:        "list",
			Aliases:     []string{"ls"},
			Description: "list buckets",
			Usage:       "",
			Args:        framework.NoArgs,
			Action:      c.list,
		},
		framework.Command{
			Name:        "use",
			Description: "enter a bucket",
			Usage:       "<bucketname>",
			Args:        framework.ExactArgs(1),
			Action:      c.use,
		},
		framework.Command{
			Name:        "drop",
			Aliases:     []string{"rm"},
			Description: "delete a bucket",
			Usage:       "<bucketname>",
			Args:        framework.ExactArgs(1),
			Action:      c.drop,
		},
	)
}

func (c *rootContext) create(args []string) error {
	bucketName := args[0]
	if err := c.DB().CreateBucket(bucketName); err != nil {
		return fmt.Errorf("%s: %v", bucketName, err)
	}
	return nil
}

func (c *rootContext) list(args []string) error {
	buckets, err := c.DB().Buckets()
	if err != nil {
		return err
	}
	c.EchoLine(strings.Join(buckets, "\n"))
	return nil
}

func (c *rootContext) use(args []string) error {
	bucketName := args[0]
	if err := c.DB().UseBucket(bucketName); err != nil {
		return fmt.Errorf("%s: %v", bucketName, err)
	}
	ctx := NewBucketContext(c.Framework, bucketName)
	c.Framework.EnterContext(ctx)
	return nil
}

func (c *rootContext) drop(args []string) error {
	bucketName := args[0]
	if err := c.DB().DeleteBucket(bucketName); err != nil {
		return fmt.Errorf("%s: %v", bucketName, err)
	}
	return nil
}
