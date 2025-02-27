package main

/*
#include <stdlib.h>
#include <sys/socket.h>
#include <vsa.h>
#include <netinet/in.h>
#include <cache/cache.h>

struct vmod_dynamic_director {
    uintptr_t priv;
};

*/
import "C"
import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"runtime/cgo"
	"strconv"
	"time"
	"unsafe"
)

type Director struct {
	backends map[string]C.VCL_BACKEND
}

//export vmod_director__init
func vmod_director__init(ctx *C.struct_vrt_ctx, director **C.struct_vmod_dynamic_director, vclName C.VCL_STRING) {
	directorPointer := C.malloc(C.size_t(unsafe.Sizeof(C.struct_vmod_dynamic_director{})))
	dir := new(Director)
	dir.backends = make(map[string]C.VCL_BACKEND)
	handle := cgo.NewHandle(dir)
	d := (*C.struct_vmod_dynamic_director)(directorPointer)
	d.priv = C.uintptr_t(handle)
	*director = d
}

//export vmod_director__fini
func vmod_director__fini(directorPtr **C.struct_vmod_dynamic_director) {
	director := *directorPtr
	handle := cgo.Handle(director.priv)
	val := handle.Value().(*Director)
	for _, b := range val.backends {
		C.free(unsafe.Pointer(b))
	}
	clear(val.backends)
	handle.Delete()
	C.free(unsafe.Pointer(director))
	*directorPtr = nil
}

//export vmod_director_backend
func vmod_director_backend(ctx *C.struct_vrt_ctx, director *C.struct_vmod_dynamic_director, key C.VCL_STRING, port C.VCL_INT) C.VCL_BACKEND {
	if port < 0 || port > 65535 {
		fail(ctx, "Invalid port")
		return nil
	}
	domain := C.GoString(key)
	d, ok := cgo.Handle(director.priv).Value().(*Director)
	if !ok {
		fail(ctx, "Failed to get director")
		return nil
	}
	b, ok := d.backends[domain+":"+strconv.Itoa(int(port))]
	if ok {
		fmt.Println("Using cached backend for", domain+":"+strconv.Itoa(int(port)))
		return b
	}
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * 1000,
			}
			return d.DialContext(ctx, network, "8.8.8.8:53")
		},
	}
	ips, err := resolver.LookupIP(context.Background(), "ip4", domain)
	if err != nil {
		fail(ctx, "Failed to resolve domain")
		return nil
	}
	var resolvedIP net.IP
	for _, ip := range ips {
		if ip.To4() != nil {
			resolvedIP = ip
			break
		}
	}

	backendPointer := C.malloc(C.size_t(unsafe.Sizeof(C.struct_vrt_backend{})))
	defer C.free(backendPointer)
	backend := (*C.struct_vrt_backend)(backendPointer)
	endpointPointer := C.malloc(C.size_t(unsafe.Sizeof(C.struct_vrt_endpoint{})))
	defer C.free(endpointPointer)
	endpoint := (*C.struct_vrt_endpoint)(endpointPointer)

	sockaddr := C.malloc(C.size_t(unsafe.Sizeof(C.struct_sockaddr_in{})))
	defer C.free(sockaddr)
	sa := (*C.struct_sockaddr_in)(sockaddr)
	sa.sin_family = C.AF_INET
	sa.sin_port = C.in_port_t(portToNetworkOrder(uint16(port)))
	sa.sin_addr.s_addr = C.uint32_t(ipToNetworkOrder(resolvedIP))

	endpoint.ipv4 = C.VSA_Malloc(sockaddr, 16)
	defer C.free(unsafe.Pointer(endpoint.ipv4))
	endpoint.magic = C.VRT_ENDPOINT_MAGIC
	backend.endpoint = (*C.struct_vrt_endpoint)(endpointPointer)
	backend.magic = C.VRT_BACKEND_MAGIC
	name := C.CString("boot")
	defer C.free(unsafe.Pointer(name))
	backend.vcl_name = name
	backend.connect_timeout = 1000
	backend.first_byte_timeout = 1000
	backend.between_bytes_timeout = 1000
	backend.backend_wait_timeout = 1000
	b = C.VRT_new_backend(ctx, (*C.struct_vrt_backend)(backendPointer), nil)
	d.backends[domain+":"+strconv.Itoa(int(port))] = b
	return b
}

func ipToNetworkOrder(ip net.IP) uint32 {
	ip = ip.To4()
	return binary.LittleEndian.Uint32(ip)
}

func portToNetworkOrder(port uint16) uint16 {
	return binary.LittleEndian.Uint16([]byte{byte(port >> 8), byte(port)})
}

func main() {}
