//go:build kandelo

package unix

import "syscall"

type Signal = syscall.Signal
type Errno = syscall.Errno
type SysProcAttr = syscall.SysProcAttr
type Termios = syscall.Termios
type Winsize = syscall.Winsize
type Rlimit = syscall.Rlimit

const (
	O_DSYNC = syscall.O_DSYNC
	O_RDWR  = syscall.O_RDWR

	SOL_SOCKET   = syscall.SOL_SOCKET
	SO_KEEPALIVE = syscall.SO_KEEPALIVE
	SO_REUSEPORT = syscall.SO_REUSEPORT

	PROT_READ      = syscall.PROT_READ
	PROT_WRITE     = syscall.PROT_WRITE
	MAP_SHARED     = syscall.MAP_SHARED
	MS_SYNC        = syscall.MS_SYNC
	MADV_NORMAL    = syscall.MADV_NORMAL
	MADV_RANDOM    = syscall.MADV_RANDOM
	LOCK_SH        = syscall.LOCK_SH
	LOCK_EX        = syscall.LOCK_EX
	LOCK_NB        = syscall.LOCK_NB
	LOCK_UN        = syscall.LOCK_UN
	TCGETS         = syscall.TCGETS
	TCSETS         = syscall.TCSETS
	TIOCGWINSZ     = syscall.TIOCGWINSZ
	RLIMIT_MEMLOCK = syscall.RLIMIT_MEMLOCK
)

func Access(path string, mode uint32) error {
	return syscall.Access(path, mode)
}

func SetsockoptInt(fd, level, opt, value int) error {
	return syscall.SetsockoptInt(fd, level, opt, value)
}

func Mmap(fd int, offset int64, length, prot, flags int) ([]byte, error) {
	return syscall.Mmap(fd, offset, length, prot, flags)
}

func Munmap(data []byte) error {
	return syscall.Munmap(data)
}

func Msync(data []byte, flags int) error {
	return syscall.Msync(data, flags)
}

func Madvise(data []byte, advice int) error {
	return syscall.Madvise(data, advice)
}

func Mlock(data []byte) error {
	return syscall.Mlock(data)
}

func Munlock(data []byte) error {
	return syscall.Munlock(data)
}

func Flock(fd int, how int) error {
	return syscall.Flock(fd, how)
}

func IoctlGetTermios(fd int, request uint) (*Termios, error) {
	if request != TCGETS {
		return nil, syscall.EINVAL
	}
	return syscall.GetTermios(fd)
}

func IoctlSetTermios(fd int, request uint, termios *Termios) error {
	if request != TCSETS {
		return syscall.EINVAL
	}
	return syscall.SetTermios(fd, termios)
}

func IoctlGetWinsize(fd int, request uint) (*Winsize, error) {
	if request != TIOCGWINSZ {
		return nil, syscall.EINVAL
	}
	return syscall.GetWinsize(fd)
}

func Getrlimit(resource int, limit *Rlimit) error {
	return syscall.Getrlimit(resource, limit)
}
