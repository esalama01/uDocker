package src

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)




func Run() {
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...) //the /proc/.. is used to spawn the child process that initializes the UTS mnamespace before executing this (main Run) command.
	//cmd := exec.Command(os.Args[2], os.Args[3:]...)


	 
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | //the UTS clone call isolates the hostname
		
		 syscall.CLONE_NEWNS, //The mount (NEWNS) cllonbe sys call isolates the mount points. --> leads to changiing teh root filesystem.

		//area for improvement: add cap_sys_admin.
	}

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func Child() {
	err := syscall.Sethostname([]byte("new_container"))

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	//CombinedOutput runs the command and returns its combined standard output and standard error.
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	//Setting the hostname
	if err != nil {
		panic(err)
	}
	err = cmd.Run()
	if err != nil {
		panic(err)
	}

	fmt.Println("Hostname updated successfully!")
}
