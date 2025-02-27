package main

/*
#include <stdlib.h>
#include <vdef.h>
#include <vrt.h>

struct vmod_etcd_client {
	uintptr_t priv;
};
*/
import "C"
import (
	"context"
	"fmt"
	"runtime/cgo"
	"time"
	"unsafe"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdClient struct {
	client *clientv3.Client
}

//export vmod_client__init
func vmod_client__init(ctx *C.struct_vrt_ctx, client **C.struct_vmod_etcd_client, vclName C.VCL_STRING, endpointStrands C.VCL_STRANDS, timeoutSecs C.VCL_INT) {
	if endpointStrands.n == 0 {
		fail(ctx, "No endpoints provided")
		return
	}

	clientPointer := C.malloc(C.size_t(unsafe.Sizeof(C.struct_vmod_etcd_client{})))
	etcdClient := new(EtcdClient)

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   strandsToStrings(endpointStrands),
		DialTimeout: time.Duration(timeoutSecs) * time.Second,
	})
	if err != nil {
		fail(ctx, fmt.Sprintf("Failed to create etcd client: %v", err))
		return
	}

	etcdClient.client = cli
	handle := cgo.NewHandle(etcdClient)
	c := (*C.struct_vmod_etcd_client)(clientPointer)
	c.priv = C.uintptr_t(handle)
	*client = c
}

func strandsToStrings(strands C.VCL_STRANDS) []string {
	tmp := unsafe.Slice(strands.p, strands.n)
	strings := make([]string, strands.n)
	for i := range strands.n {
		strings[i] = C.GoString(tmp[i])
	}
	return strings
}

//export vmod_client__fini
func vmod_client__fini(clientPtr **C.struct_vmod_etcd_client) {
	client := *clientPtr
	handle := cgo.Handle(client.priv)
	val := handle.Value().(*EtcdClient)
	val.client.Close()
	handle.Delete()
	C.free(unsafe.Pointer(client))
	*clientPtr = nil
}

//export vmod_client_set
func vmod_client_set(ctx *C.struct_vrt_ctx, client *C.struct_vmod_etcd_client, key C.VCL_STRING, value C.VCL_STRING) C.VCL_INT {
	etcdClient := cgo.Handle(client.priv).Value().(*EtcdClient)

	_, err := etcdClient.client.Put(context.Background(), C.GoString(key), C.GoString(value))
	if err != nil {
		fail(ctx, fmt.Sprintf("Failed to set key: %v", err))
		return 0
	}
	return 1
}

//export vmod_client_get
func vmod_client_get(ctx *C.struct_vrt_ctx, client *C.struct_vmod_etcd_client, key C.VCL_STRING) C.VCL_STRING {
	etcdClient := cgo.Handle(client.priv).Value().(*EtcdClient)

	resp, err := etcdClient.client.Get(context.Background(), C.GoString(key))
	if err != nil {
		fail(ctx, fmt.Sprintf("Failed to get key: %v", err))
		return nil
	}

	if len(resp.Kvs) == 0 {
		return nil
	}

	return C.CString(string(resp.Kvs[0].Value))
}

func main() {}
