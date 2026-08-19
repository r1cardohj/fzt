"""Drive ./fzt (Go binary) in a pty to test end-to-end."""
import fcntl
import os
import pty
import select
import struct
import sys
import termios
import time

CMD = ["./fzt", "demo"]


def set_winsize(fd, rows=24, cols=80):
    fcntl.ioctl(fd, termios.TIOCSWINSZ, struct.pack("HHHH", rows, cols, 0, 0))


def drain(fd, wait=0.4):
    out = b""
    while True:
        r, _, _ = select.select([fd], [], [], wait)
        if not r:
            break
        try:
            data = os.read(fd, 65536)
        except OSError:
            break
        if not data:
            break
        out += data
    return out


def run_session(keys_steps, settle=1.0, grace=5.0):
    pid, fd = pty.fork()
    if pid == 0:
        os.execv(CMD[0], CMD)
    set_winsize(fd)
    time.sleep(settle)
    out = drain(fd)
    for delay, keys in keys_steps:
        time.sleep(delay)
        os.write(fd, keys)
        out += drain(fd)
    # wait for exit, kill on timeout
    deadline = time.time() + grace
    while time.time() < deadline:
        done, status = os.waitpid(pid, os.WNOHANG)
        if done:
            out += drain(fd)
            return out, os.waitstatus_to_exitcode(status)
        out += drain(fd, wait=0.2)
    os.kill(pid, 9)
    os.waitpid(pid, 0)
    out += drain(fd)
    return out, "KILLED"


def tail_path(out):
    text = out.decode(errors="replace").replace("\r", "")
    lines = [l.strip() for l in text.split("\n") if l.strip()]
    return lines[-1] if lines else ""


if __name__ == "__main__":
    mode = sys.argv[1]
    if mode == "nav":
        # j -> src, Enter -> expand, j j -> main.py, Enter -> select
        out, code = run_session([
            (0.4, b"j"), (1.0, b"\r"), (0.5, b"j"), (0.5, b"j"), (0.5, b"\r"),
        ])
        print("exit:", code, "| final:", tail_path(out))
    elif mode == "search":
        # '/' opens search box, type query, Enter accepts top match
        out, code = run_session([
            (0.3, b"/"),
            (0.2, b"h"), (0.2, b"e"), (0.2, b"l"),
            (0.2, b"p"), (0.2, b"e"), (0.2, b"r"),
            (0.6, b"\r"),
        ])
        print("exit:", code, "| final:", tail_path(out))
    elif mode == "escback":
        # '/' search mode, type, esc -> back to tree, esc -> quit
        out, code = run_session([
            (0.3, b"/"), (0.2, b"s"), (0.2, b"r"), (0.2, b"c"),
            (0.5, b"\x1b"), (0.5, b"\x1b"),
        ])
        print("exit:", code, "(expect 130)")
