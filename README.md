# go-varnish

Collection of experimental VMODs written in Go.
The main reason for writing VMODs in Go, is to quickly prototype new ideas and experiment with existing Go libraries.
However, the VMODs are not production ready and should not be used in production environments.


Current VMODs:

- [etcd](vmods/etcd) - VMOD for etcd v3 client
- [dynamic](vmods/dynamic) - VMOD for dynamic backends
- [simple](vmods/simple) - Simple example VMOD

## Requirements

- Go 1.24+
- Varnish 7.7
- Make

## Building

In VMOD directory:

```bash
make
```

## Running

In VMOD directory:

```bash
sudo varnishd -F -f $(pwd)/example.vcl -j none -p vmod_path=$(pwd)
```

This will start varnish with the example VCL file and the VMODs in the current directory.


