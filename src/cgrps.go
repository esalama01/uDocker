package src

import (
	"os"
	"path/filepath"
)

func Configure_cgroups() {
	path := "/sys/fs/cgroup/init.scope/" //init.scope is the one responsible for pid #1 cgroup controllers and all sub process in the tree.

	//setting the memory max to 100gb
	full_path1 := filepath.Join(path, "memory.max")
	data := []byte("107374182400")
	err := os.WriteFile(full_path1, data, 0644)
	if err != nil {
		panic(err)
	}

	//setting the cpu max to 1/2 core limit
	full_path2 := filepath.Join(path, "cpu.max")
	data = []byte("50000 100000")
	err = os.WriteFile(full_path2, data, 0644)
	if err != nil {
		panic(err)
	}
}

//cgrroups v2 are a lots different than cgroups1.
