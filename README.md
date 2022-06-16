# Litra driver for Go

This is a library for interacting with Logitech Litra Glow devices

It is offered with absolutely no warranty

## Device status

Litra defines a struct to receive updates on device status

```golang
type Device struct {
    Id uint64
    Connected bool
}
```
 * `id` is an id identifying a light. This may change if the light is disconnected and reconnected
 * `connected` where `true` means this device just connected, `false` means this device just disconnected

The client should start Litra with a callback to receive device connection updates.

If a device is physically connected when Litra starts there should always be a `connected = true` status when it's found.

```golang
onDeviceStatus := func(status Device) {
    if (status.Connected) {
        log.Printf("Device %d connected", status.Id);
    } else {
        log.Printf("Device %d disconnected", status.Id);
    }
}

litra.Start(onDeviceStatus, nil)
```

If clients are expecting multiple lights to be in use, it's their responsibility to track the state of which are connected
or disconnected at any given time using the `id` field

## Reading Light state

Litra defines the following structs to receive status updates:

```golang
type Power struct { Value bool }

type Brightness struct { Value uint8 }

type Temperature struct { Value uint16 }

type LightState struct {
    Id uint64
    Power *Power
    Brightness *Brightness
    Temperature *Temperature
}
```
Not all of the pointers will be present in the message and `nil` means 'unchanged'

 * `id` is the id of the relevant device
 * `Power` where `true` means power on, `false` means power off
 * `Brightness` where `0x00` (0) is minimum (which isn't 'off') and `0xFF` (255) is maximum
 * `Temperature` is the temperature of the light in the range 2700-6500K. Out of bounds values will be snapped into range

The client should start Litra with a callback as second argument, to receive light status updates. To avoid extra nil checking
the Lightstate has a convenience method `ApplyState`.

```golang
updatePower := func(id uint64, power bool) {
    if (power) {
        log.Printf("Device %d On", id)
    } else {
        log.Printf("Device %d Off", id)
    }
}

updateBrightness := func(id uint64, brightness uint8) {
    log.Printf("Device %d set to brightness %d", id, brightness)
}

updateTemperature := func(id uint64, temperature uint16) {
    log.Printf("Device %d set to temperature %d", id, temperature)
}

onLightState := func(status LightState) {
    status.ApplyState(updatePower, updateBrightness, updateTemperature)
}

litra.Start(nil, onLightState)
```

It's possible to receive a status update before the connect message if buttons are being pressed as it's plugged in.

Due to the underlying USB protocol it is not possible to provide a full state at startup. The first messages will only 
be recieved when we successfully request the new state, or the physical buttons are used.

## Requesting Light state

Litra uses the same struct to request state updates. Not all fields needs to be provided and `Id` will default to `litra.ALL_LIGHTS`
and doesn't need to be provided unless you need to address individual lights.

Convenience setters allow the raw values to be passed without worrying about pointers:

```golang
litra.Start(nil, nil)

s := litra.LightState{}
s.SetPower(true)
s.SetBrightness(0xFF)
s.SetTemperature(5000)

litra.Request(s)
```

It's important to note that there is no return from this function. Clients should wait for their light state callback
to be invoked, (or maybe receive device disconnected message).

A single `Request` may result in multiple state updates. For instance if all fields were set, you may receive three state 
updates each with 1 of the fields set. Note also that the `Id` will be set in the updates even if you didn't set it in the 
request.
