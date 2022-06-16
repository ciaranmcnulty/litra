package litra

import (
    "github.com/ciarancmnulty/litra-go/hid"
    "log"
)

const ALL_LIGHTS=0

type Device struct {
    Id uint64
    Connected bool
}

type Power struct {
    Value bool
}

type Brightness struct {
    Value uint8
}

type Temperature struct {
    Value uint16
}

type LightState struct {
    Id uint64
    Power *Power
    Brightness *Brightness
    Temperature *Temperature
}

func (s LightState) ApplyState (p func(uint64, bool), b func (uint64, uint8), t func (uint64, uint16)) {
    if (s.Power != nil) { p(s.Id, s.Power.Value) }
    if (s.Brightness != nil) { b(s.Id, s.Brightness.Value) }
    if (s.Temperature != nil) { t(s.Id, s.Temperature.Value) }
}

func (s *LightState) SetPower (p bool) {
    s.Power = &Power{p}
}

func (s *LightState) SetBrightness (b uint8) {
    s.Brightness = &Brightness{b}
}

func (s *LightState) SetTemperature (t uint16) {
    s.Temperature = &Temperature{t}
}

func Start(onDevice func(Device), onLight func(LightState)) {

    if (onDevice == nil) {
        onDevice = func(_ Device) {}
    }

    if (onLight == nil) {
        onLight = func(_ LightState) {}
    }

    // call clients back with regular updates
    deviceConnectCh := make(chan uint64)
    deviceDisConnectCh := make(chan uint64)
    lightStateCh := make(chan hid.ReadBytes)

    go func() {
        for {
            select {
                case id := <- deviceConnectCh:
                    onDevice(Device{id, true})
                case id := <- deviceDisConnectCh:
                    onDevice(Device{id, false})
                case lightStateBytes := <- lightStateCh:
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
        lightStateCh,
    );
}

const BYTE_PREFIX_1 = 0x11
const BYTE_PREFIX_2 = 0xff
const BYTE_PREFIX_3 = 0x04

const BYTE_POWER_CHANGED = 0x00
const BYTE_BRIGHTNESS_CHANGED = 0x10
const BYTE_TEMPERATURE_CHANGED = 0x20

func lightStateFromBytes(id uint64, bytes [6]byte) *LightState {

    if bytes[0] != BYTE_PREFIX_1 || bytes[1] != BYTE_PREFIX_2 || bytes[2] != BYTE_PREFIX_3 {
        log.Print(bytes)
        return nil
    }

    lightState := LightState{}
    lightState.Id = id

    switch (bytes[3]) {

        // sent by button presses
        case BYTE_POWER_CHANGED:
            lightState.SetPower(bytes[4] == 0x01)
        case BYTE_BRIGHTNESS_CHANGED:
            lightState.SetBrightness(bytes[5])
        case BYTE_TEMPERATURE_CHANGED:
            lightState.SetTemperature(uint16(256 * int(bytes[4]) + int(bytes[5])))

       default:
            log.Print(bytes)
            return nil
    }

    return &lightState
}
