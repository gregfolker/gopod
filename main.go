// Project: gopod
// Author: Greg Folker
//
// A simple Go program designed to create a running container, inspired by this article:
// https://www.infoq.com/articles/build-a-container-in-golang/

package main

import (
	"fmt"
   "os"
   "os/exec"
   "syscall"
)

func main() {
   switch os.Args[1] {
   case "run":
      parent()
   case "child":
      child()
   default:
      fmt.Printf("Unexpected option. Expected either 'run' or 'child'")
   }
}

func parent() {
   // The parent method runs `/proc/self/exe` which is a special file containing an
   // in-memory image of the current executable. In other words, we re-run ourselves but
   // pass the child process as the first argument
   cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)

   // This is required to create the new namespace for the child process that this parent is going to execute
   cmd.SysProcAttr = &syscall.SysProcAttr {
      Cloneflags: (syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS),
   }

   cmd.Stdin = os.Stdin
   cmd.Stdout = os.Stdout
   cmd.Stderr = os.Stderr

   if err := cmd.Run(); err != nil {
      fmt.Printf("Error: ", err)
      os.Exit(1)
   }
}

func child() {
   // Initial mounts are inherited from creating the namespace in the parent process
   // This `syscall.Mount()` is required to create the root filesystem inside the child process
   assert(syscall.Mount("rootfs", "rootfs", "", syscall.MS_BIND, ""))
   assert(os.MkdirAll("rootfs/oldrootfs", 0700))
   assert(os.Chdir("/"))

   cmd := exec.Command(os.Args[2], os.Args[3:]...)

   cmd.Stdin = os.Stdin
   cmd.Stdout = os.Stdout
   cmd.Stderr = os.Stderr

   if err := cmd.Run(); err != nil {
      fmt.Printf("Error: ", err)
      os.Exit(1)
   }
}

func assert(err error) {
   if err != nil {
      panic(err)
   }
}
