package shell

func (c *Command) Subscribe(regx string, message string) {

}

func (c *Command) GetOutput() []string {
	if c.outputWriter != nil {
		return c.outputWriter.GetLines()
	}
	return make([]string, 0)
}
