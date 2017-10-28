package command

import (
	"github.com/ikaven1024/bolt-cli/cli/framework"
	"github.com/ikaven1024/bolt-cli/cli/util/color_str"
)

type bucketContext struct {
	*framework.BaseContext

	bucketName string
}

func NewBucketContext(fw framework.Framework, bucketName string) framework.Context {
	c := &bucketContext{
		BaseContext: framework.NewBaseContext(fw),
		bucketName:  bucketName,
	}
	fw.ReadLine().SetPrompt(color_str.Red(bucketName + "»"))
	c.init()
	return c
}

func (c *bucketContext) Exit() {

}

func (c *bucketContext) init() {
	c.AddCommand(
		framework.Command{
			Name:        "list",
			Aliases:     []string{"ls"},
			Description: "list all key-values",
			Usage:       "",
			Args:        framework.NoArgs,
			Action:      c.list,
		},
		framework.Command{
			Name:        "put",
			Description: "put key value",
			Usage:       "<key> <value>",
			Args:        framework.ExactArgs(2),
			Action:      c.put,
		},
		framework.Command{
			Name:        "get",
			Description: "get value of the key",
			Usage:       "<key>",
			Args:        framework.ExactArgs(1),
			Action:      c.get,
		},
		framework.Command{
			Name:        "rm",
			Aliases:     []string{"delete", "remove", "del"},
			Description: "get value of the key",
			Usage:       "<key>",
			Args:        framework.ExactArgs(1),
			Action:      c.get,
		},
	)
}

func (c *bucketContext) list(args []string) error {
	kv, err := c.DB().List()
	if err != nil {
		return err
	}

	tw := c.TabWriter().Init()
	defer tw.Flush()
	for k, v := range kv {
		tw.AppendRow(k, v)
	}
	return nil
}

func (c *bucketContext) put(args []string) error {
	return c.DB().Put(args[0], args[1])
}

func (c *bucketContext) get(args []string) error {
	v, err := c.DB().Get(args[0])
	if err != nil {
		return err
	}
	c.EchoLine(v)
	return nil
}

func (c *bucketContext) rm(args []string) error {
	return c.DB().Remove(args[0])
}
