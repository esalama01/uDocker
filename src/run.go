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
			syscall.CLONE_NEWUSER,
		Unshareflags: syscall.CLONE_NEWNS, //unshare ensures mounts are private to this namespace

		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID  : 0,
				HostID : os.Getuid(),
				Size: 1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID  : 0,
				HostID : os.Getuid(),
				Size: 1,
			},
		},
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
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	//Setting the hostname
	err := syscall.Sethostname([]byte("Container"))
	//setting the root directory.
	syscall.Chroot("/home/esalama01/projects/uDocker/alpinefs") //example usage after this change: sudo ./uDocker run /bin/busybox pwd.
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
