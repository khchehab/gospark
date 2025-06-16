package spark

type Config struct {
	BgColor   string
	FgColor   string
	ShowSum   bool
	ShowStats bool
}

func (c *Config) Validate() error {
	var err error
	if err = ValidateColor(c.BgColor); err != nil {
		return err
	}
	if err = ValidateColor(c.FgColor); err != nil {
		return err
	}
	return nil
}
