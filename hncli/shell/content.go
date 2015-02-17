package shell

import "fmt"

func (c *Cli) SetContent(content string) error {
	if c.Writer != nil {
		_, err := fmt.Fprint(c.Writer, content)

		if err != nil {
			return fmt.Errorf("Error writing to output: %v", err)
		}
	} else {
		c.Window.Print(content)
		c.Window.Refresh()
	}

	return nil
}
