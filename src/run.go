package src

import (
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
			syscall.CLONE_NEWNS | //The mount (NEWNS) namespace call isolates the mount points. --> leads to changing the root filesystem.
			syscall.CLONE_NEWPID | //The PID namespace isolates the process id.
			syscall.CLONE_NEWUSER, //the user namespace isolates the security related identifiers.
		Unshareflags: syscall.CLONE_NEWNS, //unshare ensures mounts are private to this namespace
		//Chroot:       "/home/esalama01/projects/uDocker/alpinefs",
		/* Must grant the container root privileges within itelf but looks like a normal process from the host os's perspective.*/
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},

		//UseCgroupFD: true,

		//area for improvement: add cap_sys_admin.
	}
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func Child() {
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	//CombinedOutput runs the command and returns its combined standard output and standard error.
	Configure_cgroups() //cgroups are v2 not v1, the step descriptiopn is outdated tho
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	//Setting the hostname
	err := syscall.Sethostname([]byte("Container"))
	//setting the root directory.
	syscall.Chroot("/home/esalama01/projects/uDocker/output") //example usage after this change: sudo ./uDocker run /bin/busybox pwd.
	os.Chdir("/")
	//mounting the virtual fs /proc
	syscall.Mount("proc", "proc", "proc", 0, "")
	if err != nil {
		panic(err)
	}
	err = cmd.Run()
	if err != nil {
		panic(err)
	}

	syscall.Unmount("proc", 0)
}
