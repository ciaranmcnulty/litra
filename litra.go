package litra

type UsbProvider interface {
    Start()
	SendBytesToDevice(uint64, [20]byte)
	SetOnDeviceConnect(func(uint64))
	SetOnDeviceDisconnect(func(uint64,))
	SetOnBytesFromDevice(func(uint64, [6]byte))
}

type Device struct {
    Id uint64
    Connected bool
}

type Litra struct {
   usbProvider UsbProvider
   onDevice func(Device)
   onLightState func(LightState)
}

func (l *Litra) Start(
    usbProvider UsbProvider,
    onDevice func(Device),
    onLightState func(LightState),
) {
    if (onDevice == nil) {
        l.onDevice = func(_ Device) {}
    } else {
        l.onDevice = onDevice
    }

    if (onLightState == nil) {
        l.onLightState = func(_ LightState) {}
    } else {
        l.onLightState = onLightState
    }

    usbProvider.Start()
    usbProvider.SetOnDeviceConnect(func(id uint64) {
        l.onDevice(Device{id, true})
    })
    usbProvider.SetOnDeviceDisconnect(func(id uint64) {
        l.onDevice(Device{id, false})
    })
    usbProvider.SetOnBytesFromDevice(func(id uint64, bytes [6]byte) {
        lightState := lightStateFromBytes(id, bytes)
        if lightState != nil {
            l.onLightState(*lightState)
        }
    })
    l.usbProvider = usbProvider
}

func (l *Litra) Request(s LightState) {
    for _, bytes := range bytesFromLightState(s) {
        l.usbProvider.SendBytesToDevice(s.Id, bytes)
    }
}
