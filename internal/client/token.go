package client

import (
	"fmt"
	"os"
	"text/tabwriter"

	log "github.com/sirupsen/logrus"

	"github.com/your-overtime/api/pkg"
)

func (c *Client) CreateTokenViaCli(name string) error {
	t, err := c.CreateToken(name)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, '.', tabwriter.TabIndent)
	fmt.Fprintf(w, "ID\t: %d\n", t.ID)
	fmt.Fprintf(w, "Token\t: %s\n", t.Token)
	fmt.Fprintf(w, "Name\t: %s\n\n", t.Name)
	w.Flush()

	return nil
}

func (c *Client) CreateToken(name string) (*pkg.Token, error) {
	t, err := c.ots.CreateToken(pkg.InputToken{Name: name})

	if err != nil {
		log.Debug(err)
		return nil, err
	}

	return t, nil
}

func (c *Client) DeleteToken(id uint) error {
	err := c.ots.DeleteToken(id)

	if err != nil {
		log.Debug(err)
		return err
	}

	fmt.Println("Token deleted")

	return nil
}

func (c *Client) GetTokens() error {
	ts, err := c.ots.GetTokens()

	if err != nil {
		log.Debug(err)
		return err
	}

	fmt.Println()
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, '.', tabwriter.TabIndent)
	for _, t := range ts {
		fmt.Fprintf(w, "ID\t: %d\n", t.ID)
		fmt.Fprintf(w, "Token\t: %s\n", t.Token)
		fmt.Fprintf(w, "Name\t: %s\n\n", t.Name)
	}

	w.Flush()

	return nil
}
