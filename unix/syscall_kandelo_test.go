package unix_test

import (
	"os"
	"syscall"
	"testing"

	"golang.org/x/sys/unix"
)

func TestKandeloUnixDup(t *testing.T) {
	file, err := os.CreateTemp("", "kandelo-xsys-dup-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	defer file.Close()
	if _, err := file.WriteString("abc"); err != nil {
		t.Fatal(err)
	}
	if _, err := file.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	fd := int(file.Fd())
	if flags, err := syscall.Fcntl(fd, syscall.F_GETFD, 0); err != nil || flags&syscall.FD_CLOEXEC == 0 {
		t.Fatalf("original descriptor close-on-exec: flags=%d err=%v", flags, err)
	}
	duplicate, err := unix.Dup(fd)
	if err != nil {
		t.Fatal(err)
	}
	defer syscall.Close(duplicate)
	if duplicate == fd {
		t.Fatal("Dup returned the original descriptor")
	}
	if flags, err := syscall.Fcntl(duplicate, syscall.F_GETFD, 0); err != nil || flags&syscall.FD_CLOEXEC != 0 {
		t.Fatalf("duplicate close-on-exec must be clear: flags=%d err=%v", flags, err)
	}
	var data [1]byte
	if n, err := syscall.Read(duplicate, data[:]); err != nil || n != 1 || data[0] != 'a' {
		t.Fatalf("duplicate read: n=%d data=%q err=%v", n, data, err)
	}
	if n, err := file.Read(data[:]); err != nil || n != 1 || data[0] != 'b' {
		t.Fatalf("shared file offset: n=%d data=%q err=%v", n, data, err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	if n, err := syscall.Read(duplicate, data[:]); err != nil || n != 1 || data[0] != 'c' {
		t.Fatalf("duplicate after original close: n=%d data=%q err=%v", n, data, err)
	}
	if _, err := unix.Dup(-1); err != syscall.EBADF {
		t.Fatalf("invalid descriptor: got %v, want EBADF", err)
	}
}

func TestKandeloUnixMemory(t *testing.T) {
	file, err := os.CreateTemp("", "kandelo-xsys-mmap-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	defer file.Close()
	pageSize := syscall.Getpagesize()
	if err := file.Truncate(int64(pageSize)); err != nil {
		t.Fatal(err)
	}
	data, err := unix.Mmap(int(file.Fd()), 0, pageSize, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_SHARED)
	if err != nil {
		t.Fatal(err)
	}
	if err := unix.Madvise(data, unix.MADV_RANDOM); err != nil {
		t.Fatal(err)
	}
	if err := unix.Mlock(data); err != nil {
		t.Fatal(err)
	}
	if err := unix.Munlock(data); err != nil {
		t.Fatal(err)
	}
	data[0] = 'X'
	if err := unix.Msync(data, unix.MS_SYNC); err != nil {
		t.Fatal(err)
	}
	if err := unix.Munmap(data); err != nil {
		t.Fatal(err)
	}
	stored := make([]byte, 1)
	if _, err := file.ReadAt(stored, 0); err != nil {
		t.Fatal(err)
	}
	if stored[0] != 'X' {
		t.Fatalf("mapped writeback: got %q, want X", stored[0])
	}
}

func TestKandeloUnixSocketOption(t *testing.T) {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer syscall.Close(fd)
	if err := unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_REUSEPORT, 1); err != nil {
		t.Fatal(err)
	}
	value, err := syscall.GetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_REUSEPORT)
	if err != nil {
		t.Fatal(err)
	}
	if value != 1 {
		t.Fatalf("SO_REUSEPORT value: got %d, want 1", value)
	}
}

func TestKandeloUnixFlockAccess(t *testing.T) {
	first, err := os.CreateTemp("", "kandelo-xsys-flock-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(first.Name())
	defer first.Close()
	second, err := os.OpenFile(first.Name(), os.O_RDWR, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer second.Close()
	if err := unix.Access(first.Name(), 0); err != nil {
		t.Fatal(err)
	}
	if err := unix.Flock(int(first.Fd()), unix.LOCK_EX|unix.LOCK_NB); err != nil {
		t.Fatal(err)
	}
	if err := unix.Flock(int(second.Fd()), unix.LOCK_EX|unix.LOCK_NB); err != syscall.EWOULDBLOCK {
		t.Fatalf("contending lock: got %v, want EWOULDBLOCK", err)
	}
	if err := unix.Flock(int(first.Fd()), unix.LOCK_UN); err != nil {
		t.Fatal(err)
	}
	if err := unix.Flock(int(second.Fd()), unix.LOCK_EX|unix.LOCK_NB); err != nil {
		t.Fatal(err)
	}
	if err := unix.Flock(int(second.Fd()), unix.LOCK_UN); err != nil {
		t.Fatal(err)
	}
}
