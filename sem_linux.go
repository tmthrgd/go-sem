// Copyright (C) 2003-2014 Free Software Foundation, Inc.
// This file is part of the GNU C Library.
// Contributed by Paul Mackerras <paulus@au.ibm.com>, 2003.
//
// The GNU C Library is free software; you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 2.1 of the License, or (at your option) any later version.
//
// The GNU C Library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with the GNU C Library; if not, see
// <http://www.gnu.org/licenses/>.

// +build linux,!386,!amd64

package sem

/*
#include <semaphore.h>      // For sem_*
#include <linux/futex.h>    // For FUTEX_*
#include <bits/local_lim.h> // For SEM_VALUE_MAX

// This is pulled from glibc-2.17/nptl/sysdeps/unix/sysv/linux/internaltypes.h
// The case of field names has been changed to be consistent with cgo -godefs
struct new_sem
{
	unsigned int Value;
	int Private;
	unsigned long int NWaiters;
};
*/
import "C"

import (
	"golang.org/x/sys/unix"
	"sync/atomic"
	"syscall"
	"unsafe"
)

type Semaphore C.sem_t

type newSem C.struct_new_sem

func New(value uint) (*Semaphore, error) {
	sem := new(Semaphore)

	if err := sem.Init(value); err != nil {
		return nil, err
	}

	return sem, nil
}

// This mirrors atomic_decrement_if_positive from glibc-2.17/include/atomic.h
func atomicDecrementIfPositive(mem *uint32) uint32 {
	for {
		if old := atomic.LoadUint32(mem); old == 0 || atomic.CompareAndSwapUint32(mem, old, old-1) {
			return old
		}
	}
}

// This (mostly?) mirrors __new_sem_wait from glibc-2.17/nptl/sysdeps/unix/sysv/linux/sem_wait.c
func (sem *Semaphore) Wait() error {
	isem := (*newSem)(unsafe.Pointer(sem))

	if atomicDecrementIfPositive((*uint32)(&isem.Value)) > 0 {
		return nil
	}

	atomic.AddUintptr((*uintptr)(unsafe.Pointer(&isem.NWaiters)), 1)

	for {
		//err = do_futex_wait(isem);
		if _, _, err := unix.Syscall6(unix.SYS_FUTEX, uintptr(unsafe.Pointer(&isem.Value)), uintptr(C.FUTEX_WAIT), 0, 0, 0, 0); err != 0 && err != syscall.EWOULDBLOCK {
			atomic.AddUintptr((*uintptr)(unsafe.Pointer(&isem.NWaiters)), ^uintptr(0))
			return err
		}

		if atomicDecrementIfPositive((*uint32)(&isem.Value)) > 0 {
			atomic.AddUintptr((*uintptr)(unsafe.Pointer(&isem.NWaiters)), ^uintptr(0))
			return nil
		}
	}
}

// This (loosely?) mirrors __new_sem_trywait from glibc-2.17/nptl/sysdeps/unix/sysv/linux/sem_trywait.c
func (sem *Semaphore) TryWait() error {
	isem := (*newSem)(unsafe.Pointer(sem))

	if atomicDecrementIfPositive((*uint32)(&isem.Value)) > 0 {
		return nil
	}

	return syscall.EAGAIN
}

// This mirrors __new_sem_post from glibc-2.17/nptl/sysdeps/unix/sysv/linux/sem_post.c
func (sem *Semaphore) Post() error {
	isem := (*newSem)(unsafe.Pointer(sem))

	for {
		cur := atomic.LoadUint32((*uint32)(&isem.Value))

		if cur == C.SEM_VALUE_MAX {
			return syscall.EOVERFLOW
		}

		if atomic.CompareAndSwapUint32((*uint32)(&isem.Value), cur, cur+1) {
			break
		}
	}

	// atomic_full_barrier ();

	if atomic.LoadUintptr((*uintptr)(unsafe.Pointer(&isem.NWaiters))) <= 0 {
		return nil
	}

	if _, _, err := unix.Syscall6(unix.SYS_FUTEX, uintptr(unsafe.Pointer(&isem.Value)), uintptr(C.FUTEX_WAKE), 1, 0, 0, 0); err != 0 {
		return err
	}

	return nil
}

// This mirrors __new_sem_init from glibc-2.17/nptl/sem_init.c
func (sem *Semaphore) Init(value uint) error {
	if value > C.SEM_VALUE_MAX {
		return syscall.EINVAL
	}

	isem := (*newSem)(unsafe.Pointer(sem))
	isem.Value = C.uint(value)
	isem.Private = 0
	isem.NWaiters = 0

	return nil
}

// This mirrors __new_sem_destroy from glibc-2.17/nptl/sem_destroy.c
func (sem *Semaphore) Destroy() error {
	return nil
}
