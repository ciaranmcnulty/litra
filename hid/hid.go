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
    "log"
)

const DEVICE_DETECT_MS = 500

const VENDOR_ID = 1133
const PRODUCT_ID = 51456

const PREFIX_1 = 0x11
const PREFIX_2 = 0xff
const PREFIX_3 = 0x04

type ReadMessage struct {
    Id uint64
    Data [6]byte
}

type WriteMessage struct {
    Id uint64
    Data [20]byte
}

//
// func Send (message WriteMessage) {
//     Init()
//
//     for id, WriteMessageCh := range(deviceWriteChannels) {
//         if (message.Id == 0) || (message.Id == id) {
//             WriteMessageCh <- message.Data
//         }
//     }
// }

func StartListening(
    deviceConnectCh chan uint64,
    deviceDisconnectCh chan uint64,
    readMessageCh chan ReadMessage,
    writeMessageCh chan WriteMessage,
) {
    nextDeviceId := uint64(1) // 0 is reserved
    connectedDevices := make(map[string]uint64)
    deviceWriteChannels := make(map[uint64]chan [20]byte)

    // write dispatcher
    go func() {
        for {
            message := <-writeMessageCh

            for id, deviceWriteMessageCh := range(deviceWriteChannels) {
                if (message.Id == 0) || (message.Id == id) {
                   deviceWriteMessageCh <- message.Data
                }
            }
        }
    }()

    // per-device loop
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
                    deviceWriteChannels[id] = make(chan [20]byte)
                    deviceConnectCh <- id

                    device := C.hid_open_path(C.CString(path))

                    go func() {
                        for {
                            bytes := <-deviceWriteChannels[id]

                            if (device == nil) || 20 != int(C.hid_write(device, (*C.uchar)(&bytes[0]), C.size_t(len(bytes)))) {
                                log.Print("Write error")
                            }
                        }
                    }()

                    go func() {
                        for {
                            out := [6]byte{}

                            if (device == nil) || (-1 == C.hid_read(device, (*C.uchar)(&out[0]), C.size_t(len(out)))) {
                                deviceDisconnectCh <- connectedDevices[path]
                                delete(connectedDevices, path)
                                delete(deviceWriteChannels, id)
                                return
                            }

                            readMessageCh <- ReadMessage{id, out}
                        }
                    }()
                }
            }
            C.hid_free_enumeration(firstDevice)

            time.Sleep(DEVICE_DETECT_MS * time.Millisecond)
        }
    }()
}
