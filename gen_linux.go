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

//go:generate sh -c "GOARCH=386 go tool cgo -godefs sem_linux.go | gofmt -r 'atomic.AddUintptr->atomic.AddUint32' | gofmt -r 'atomic.LoadUintptr->atomic.LoadUint32' | gofmt -r '(*uintptr)(unsafe.Pointer(x))->x' | gofmt -r '^uintptr(0)->^uint32(0)' > sem_linux_386.go"
//go:generate sh -c "GOARCH=amd64 go tool cgo -godefs sem_linux.go | gofmt -r 'atomic.AddUintptr->atomic.AddUint64' | gofmt -r 'atomic.LoadUintptr->atomic.LoadUint64' | gofmt -r '(*uintptr)(unsafe.Pointer(x))->x' | gofmt -r '^uintptr(0)->^uint64(0)' > sem_linux_amd64.go"

package sem
