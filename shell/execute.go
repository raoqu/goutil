package shell

import (
	"io"
	"os/exec"
)

func (c *Command) Run() error {
	c.initExec()

	if c.Manager != nil {
		manager := c.Manager
		manager.Mutex.Lock()
		manager.Processes[c.UniqueId] = c
		manager.Mutex.Unlock()
	}

	if c.Async {
		go func() {
			c.execute()
		}()
	} else {
		c.execute()
	}
	return nil
}

func (c *Command) execute() {
	err := c.Exec.Start() // start shell command
	if err != nil {
		c.Err = err
	} else {
		c.Wait()
	}
}

func (c *Command) handleClose() {
	if c.OnClose != nil {
		c.OnClose()
	}
}

func (c *Command) handleOutput(line string) {
	if c.OnOutput != nil {
		c.OnOutput(line)
	}
}

func (c *Command) Wait() {
	if c.Exec != nil && c.Err == nil {
		c.Exec.Wait()
		if c.outputWriter != nil {
			c.outputWriter.Flush()
		}
	}
	c.handleClose()
}

func (c *Command) initExec() {
	isWindows := OSFeature(true, false)
	shell := "c:\\windows\\system32\\cmd.exe"
	if !isWindows {
		shell = getDefaultShell()
	}

	flag := OSFeature("/C", "-c")
	cmd := exec.Command(shell, flag, c.Command)
	cmd.Dir = c.Dir
	c.Exec = cmd

	if c.OnOutput != nil {
		c.initOutputWriter()
	}
}

func (c *Command) initOutputWriter() io.Writer {
	writer := NewLineBufferWriter(c.BufferLines)
	writer.Handler = func(line string) {
		c.handleOutput(line)
	}
	c.Exec.Stdout = writer
	c.outputWriter = writer
	return writer
}

func (c *Command) initErrorWriter() io.Writer {
	writer := NewLineBufferWriter(c.BufferLines)
	c.Exec.Stderr = writer
	return writer
}
