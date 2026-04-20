# Python Skill Runtime

This directory contains the Python runtime for executing Skill-Plus skills written in Python.

## Security

- Runs inside gVisor container
- Network egress restricted to allowlist
- Memory capped at 512MB (configurable)
- `/tmp` is writable but cleared after each execution
- Static analysis rejects `eval`, `exec`, `os.system`, `subprocess` without declared reason
