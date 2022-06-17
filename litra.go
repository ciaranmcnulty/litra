package litra

import (
    "github.com/ciarancmnulty/litra-go/hid"
)

var writeMessageCh chan hid.WriteMessage

func Start(onDevice func(Device), onLight func(LightState)) {

    if (onDevice == nil) {
        onDevice = func(_ Device) {}
    }

    if (onLight == nil) {
        onLight = func(_ LightState) {}
    }

    deviceConnectCh := make(chan uint64)
    deviceDisConnectCh := make(chan uint64)
    readMessageCh := make(chan hid.ReadMessage)
    writeMessageCh = make(chan hid.WriteMessage)

    go func() {
        for {
            select {
                case id := <- deviceConnectCh:
                    onDevice(Device{id, true})
                case id := <- deviceDisConnectCh:
                    onDevice(Device{id, false})
                case lightStateBytes := <- readMessageCh:
                    lightState := lightStateFromBytes(lightStateBytes.Id, lightStateBytes.Data)
                    if lightState != nil {
                        onLight(*lightState)
                    }
            }
        }
    }()

    hid.StartListening(
        deviceConnectCh,
        deviceDisConnectCh,
        readMessageCh,
        writeMessageCh,
    );
}

func Request(s LightState) {
    if writeMessageCh == nil {
        return
    }

    for _, bytes := range bytesFromLightState(s) {
        writeMessageCh <- hid.WriteMessage{s.Id, bytes}
    }
}
