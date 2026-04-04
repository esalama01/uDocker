package src

import (
	"os/exec"
	"syscall"
)

func RunCommand(name string, args ...string) (string, error) {
	ChildArgs := append([]string{"Child_RunCommand", name}, args...)
	cmd := exec.Command("/proc/self/exe", ChildArgs...) //the /proc/.. is used to spawn the child process that initializes the UTS mnamespace before executing this (main Run) command. 
	out, err := cmd.CombinedOutput() //CombinedOutput runs the command and returns its combined standard output and standard error.
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS, // the UTS clone call isolates the hostname.
	}
	return string(out), err
}

func Child_RunCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput() //CombinedOutput runs the command and returns its combined standard output and standard error.
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS, // the UTS clone call isolates the hostname.
	}
	//Setting the hostname
	syscall.Sethostname([]byte("container"))
	return string(out), err
}