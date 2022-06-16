package hid

/*
#cgo CFLAGS: -I..//hidapi/hidapi
#cgo darwin CFLAGS: -DOS_DARWIN
#cgo darwin LDFLAGS: -framework CoreFoundation -framework IOKit -framework AppKit
#ifdef OS_DARWIN
	#include "../hidapi/mac/hid.c"
#endif
*/
import "C"

import (
    "time"
//     "log"
)

const DEVICE_DETECT_MS = 500

const VENDOR_ID = 1133
const PRODUCT_ID = 51456

const PREFIX_1 = 0x11
const PREFIX_2 = 0xff
const PREFIX_3 = 0x04


var nextDeviceId uint64

type ReadBytes struct {
    Id uint64
    Data [6]byte
}

func StartListening(
    deviceConnectCh chan uint64,
    deviceDisconnectCh chan uint64,
    lightStateCh chan ReadBytes,
) {

    C.hid_init()

    nextDeviceId = uint64(1) // 0 is reserved
    connectedDevices := make(map[string]uint64)

    // listening loop
    go func() {
        for {
            firstDevice := C.hid_enumerate(C.ushort(VENDOR_ID), C.ushort(PRODUCT_ID))

            for device := firstDevice; device != nil && device.next != nil; device = device.next {
                path := C.GoString(device.path);
                id, exists := connectedDevices[path]
                if(!exists) {
                    id = nextDeviceId
                    nextDeviceId++
                    connectedDevices[path] = id
                    deviceConnectCh <- id

                    go func() {
                        for {
                            out := [6]byte{}

                            device := C.hid_open_path(C.CString(path))

                            if (device == nil) || (-1 == C.hid_read(device, (*C.uchar)(&out[0]), C.size_t(len(out)))) {
                                deviceDisconnectCh <- connectedDevices[path]
                                delete(connectedDevices, path)
                                return
                            }

                            if device != nil {
                                C.hid_close(device)
                            }

                            lightStateCh <- ReadBytes{id, out}
                        }
                    }()
                }
            }
            C.hid_free_enumeration(firstDevice)

            time.Sleep(DEVICE_DETECT_MS * time.Millisecond)
        }
    }()
}
