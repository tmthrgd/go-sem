// Copyright (C) 2016  Tom Thorogood
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.

// +build !linux !amd64

package sem

//#include <semaphore.h>      // For sem_*
import "C"

type Semaphore C.sem_t

func New(value uint) (*Semaphore, error) {
	sem := new(Semaphore)

	if err := sem.Init(value); err != nil {
		return nil, err
	}

	return sem, nil
}

func (sem *Semaphore) Wait() error {
	_, err := C.sem_wait((*C.sem_t)(sem))
	return err
}

func (sem *Semaphore) TryWait() error {
	_, err := C.sem_trywait((*C.sem_t)(sem))
	return err
}

func (sem *Semaphore) Post() error {
	_, err := C.sem_post((*C.sem_t)(sem))
	return err
}

func (sem *Semaphore) Init(value uint) error {
	_, err := C.sem_init((*C.sem_t)(sem), 1, C.uint(value))
	return err
}

func (sem *Semaphore) Destroy() error {
	_, err := C.sem_destroy((*C.sem_t)(sem))
	return err
}
