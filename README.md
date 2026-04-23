# uDocker

A custom, minimal container runtime written in Go, built to understand the core Linux primitives that make containers possible. The project follows the [Build Your Own Docker](https://codingchallenges.fyi/challenges/challenge-docker/) coding challenge and implements process isolation through Linux namespaces, resource limiting via cgroups, and image pulling from the Docker Registry API.

It is **not** production software. It is a learning tool that should be run inside a Linux virtual machine or on a bare-metal Linux host.

## Completed Steps 

The challenge is split into eight steps. The table below shows which steps are complete and where the relevant code lives.

| Step | Feature | Status | Location |
|------|---------|--------|----------|
| 0 | Environment setup | Done | – |
| 1 | Run an arbitrary command | Done | `src/run.go` (`Run`) |
| 2 | Private hostname (UTS namespace) | Done | `src/run.go` (`Run` – `CLONE_NEWUTS`, `sethostname`) |
| 3 | Chroot to a filesystem | Done | `src/run.go` (`Child` – `chroot` + `chdir`) |
| 4 | PID and mount isolation (PID + mount namespaces) | Done | `src/run.go` (`Run` – `CLONE_NEWPID`, `CLONE_NEWNS`; `Child` – mounts `/proc`) |
| 5 | User namespace (rootless container) | Done | `src/run.go` (`Run` – `CLONE_NEWUSER` + UID/GID mappings) |
| 6 | Resource limits (cgroups v2) | Done | `src/cgrps.go` (`Configure_cgroups` – memory.max, cpu.max) |
| 7 | Pull image from Docker Hub | Done | `src/pull.go` (auth, manifest, layer download, extraction) |
| 8 | Run the pulled image | **Partially complete** | `src/run.go` uses `/home/esalama01/projects/uDocker/output` as root |
| 8 | **Environment variables from image config** | **Not yet implemented** | – |

## What is missing

The challenge specification asks that after pulling an image the runtime should parse the container’s configuration (the `Config` blob saved during Step 7) and set the environment variables and working directory before executing the command. The image configuration is fully fetched and parsed in `src/pull.go` (`Config_structure`, `confiiiig_struct`), and the runtime already chroots into `/home/esalama01/projects/uDocker/output`, but:

* The `Env` slice and the `WorkingDir` field from the config are **not read** at container startup.
* The runtime does **not** set any environment variables inside the container.
* The runtime does **not** adjust the working directory to match the image’s `WorkingDir`.

Because of this, containers behave as if they start with an empty environment and in `/`, which may cause images that rely on environment variables (for example `PATH`, `HOME`, `JAVA_HOME`, etc.) to fail.

This is the **only** deliberate omission from the core challenge.

## Project structure

```text
.
├── src
│ ├── run.go       # container lifecycle (namespaces, chroot, /proc)
│ ├── pull.go      # Docker Registry API (auth, manifest, layers, config)
│ └── cgrps.go     # cgroups v2 resource controls
├── go.mod
├── .gitignore
└── README.md
```

The main entry point is expected to be a file at the repository root that calls into the `src` package. (In earlier commits the file was called `uDocker.go`; if it is not present in your clone you can add a small `main.go` that switches on the `run` subcommand and invokes `src.Run`.)

## Building and running

### Prerequisites

* A Linux environment with **root access** (namespaces and cgroups require privileged operations).
* Go 1.26.1 (or later).
* Docker Hub access (public images only; authentication uses a token flow).

### Build

From the repository root:

```bash
# If a main package file exists:
go build -o uDocker .
```

If there is no `main.go`, create one with the content shown below and then build. 

**Example minimal `main.go` (place it in the repository root):**

```go
package main

import (
    "os"
    "uDocker/src"
)

func main() {
    if len(os.Args) < 2 {
        panic("usage: uDocker run <command>")
    }
    
    switch os.Args[1] {
    case "run":
        src.Run()
    default:
        panic("unknown command")
    }
}
```

### Run

The binary must be executed as **root** because it creates namespaces and writes to cgroup files.

```bash
sudo ./uDocker run /bin/sh
```

After pulling an image (see “Pulling an image” below), you can run commands inside that image by chrooting to the output directory. The current code hardcodes the chroot path to `/home/esalama01/projects/uDocker/output`. Make sure that directory contains a valid root filesystem before executing a command.

### Pulling an image

Pulling is not exposed as a separate CLI command. The pull logic lives in `src/pull.go` and must be invoked from a test or a small helper. To use it, write a short Go program that calls the exported pull functions:

```go
package main

import "uDocker/src"

func main() {
    // downloads and extracts to ..., the current code writes to ./output
    src.Pull("ubuntu") 
}
```

Then run:

```bash
go run .
```

This will fetch the `ubuntu` image (or any other public image) from Docker Hub and unpack its layers.

### Cgroups

The cgroup implementation is written for **cgroups v2** (unified hierarchy). The program writes to `/sys/fs/cgroup/init.scope/`. On systems that use cgroups v1 the writes will fail. If you see errors related to `memory.max` or `cpu.max`, check which cgroup version your kernel supports and adjust the path accordingly.

## Limitations (beyond the challenge)

* **No command-line arguments for image name or tag.** The runtime uses hardcoded values when pulling.
* **No networking namespace.** Containers share the host network stack.
* **No image caching.** Every pull re-downloads all layers.
* **No cleanup of temporary directories.** Downloaded layers remain on disk.

## References

* [Original challenge – Build Your Own Docker](https://codingchallenges.fyi/challenges/challenge-docker/)
* [Docker Registry HTTP API V2](https://docs.docker.com/registry/spec/api/)
* [Linux namespaces man page](https://man7.org/linux/man-pages/man7/namespaces.7.html)
* [cgroups v2 documentation](https://www.kernel.org/doc/Documentation/cgroup-v2.txt)
* [Containers From Scratch • Liz Rice • GOTO 2018](https://youtu.be/8fi7uSYlOdc)

## License

This project is created for educational purposes. No license is applied.
