package litra

/*
#cgo CFLAGS: -I.//hidapi/hidapi
#cgo darwin CFLAGS: -DOS_DARWIN
#cgo darwin LDFLAGS: -framework CoreFoundation -framework IOKit -framework AppKit
#ifdef OS_DARWIN
	#include "./hidapi/mac/hid.c"
#endif
*/
import "C"

import (
    "time"
)

const DEVICE_DETECT_MS = 500

const VENDOR_ID = 1133
const PRODUCT_ID = 51456

type ReadMessage struct {
    Id uint64
    Data [6]byte
}

type WriteMessage struct {
    Id uint64
    Data [20]byte
}

type HidUsbProvider struct {
    writeMessageCh chan WriteMessage
    onDeviceConnect func(uint64)
    onDeviceDisconnect func(uint64)
    onBytesFromDevice func(uint64, [6]byte)
}

func (h *HidUsbProvider) Start() {
    deviceConnectCh := make(chan uint64)
    deviceDisconnectCh := make(chan uint64)
    readMessageCh := make(chan ReadMessage)
    h.writeMessageCh = make(chan WriteMessage)

    nextDeviceId := uint64(1) // 0 is reserved
    connectedDevices := make(map[string]uint64)
    deviceWriteChannels := make(map[uint64]chan [20]byte)

    // routes writes to correct devices
    go func() {
        for {
            message := <- h.writeMessageCh

            for id, deviceWriteMessageCh := range(deviceWriteChannels) {
                if (message.Id == 0) || (message.Id == id) {
                   deviceWriteMessageCh <- message.Data
                }
            }
        }
    }()

    // client callbacks
    go func() {
        for {
            select{
                case id := <- deviceConnectCh:
                    if (h.onDeviceConnect != nil) {
                        h.onDeviceConnect(id)
                    }
                case id := <- deviceDisconnectCh:
                    if (h.onDeviceDisconnect != nil) {
                        h.onDeviceDisconnect(id)
                    }
                case message := <- readMessageCh:
                    if (h.onBytesFromDevice != nil) {
                        h.onBytesFromDevice(message.Id, message.Data)
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
                                deviceDisconnectCh <- connectedDevices[path]
                                delete(connectedDevices, path)
                                delete(deviceWriteChannels, id)
                                return
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

func (h *HidUsbProvider) SendBytesToDevice(id uint64, bytes [20]byte) {
    h.writeMessageCh <- WriteMessage{id, bytes}
}

func (h *HidUsbProvider) SetOnDeviceConnect(f func(uint64)) {
    h.onDeviceConnect = f
}
func (h *HidUsbProvider) SetOnDeviceDisconnect(f func(uint64)) {
    h.onDeviceDisconnect = f
}
func (h *HidUsbProvider) SetOnBytesFromDevice(f func(uint64, [6]byte)) {
    h.onBytesFromDevice = f
}
