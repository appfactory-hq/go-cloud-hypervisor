## Hack

### macOS

For developping on macOS, you must using Docker context to connect to the Docker daemon running on the remote linux host.

Note: Your remote host must have Docker installed and running and have the KVM opperationnal (see: `kvm-ok` cli command).

```bash
# Create a context
docker context create my-remote-host --docker "host=ssh://username@host:22"
```

With VSCode follow this [guide](https://code.visualstudio.com/docs/containers/ssh) to connect to the remote Docker host.
