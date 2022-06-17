package litra

import (
    "testing"
    "gotest.tools/assert"
)

func TestPowerOn(t *testing.T) {
    s := LightState{}

    s.SetPower(true)

    expected := [20]byte{
        0x11, 0xff, 0x04,
        0x1c, // power
        0x01, // value
        0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,    0x00, 0x00, 0x00, 0x00, 0x00,
    }

    assert.Equal(t, bytesFromLightState(s)[0], expected)
}

func TestPowerOff(t *testing.T) {
    s := LightState{}

    s.SetPower(false)

    expected := [20]byte{
        0x11, 0xff, 0x04,
        0x1c, // power
        0x00, // value
        0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,    0x00, 0x00, 0x00, 0x00, 0x00,
    }

    assert.Equal(t, bytesFromLightState(s)[0], expected)
}

func TestSetBrightness(t *testing.T) {
    s := LightState{}

    s.SetBrightness(0x80)

    expected := [20]byte{
        0x11, 0xff, 0x04, // prefix
        0x4c, // brightness
        0x00, 0x80, // value
        0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
    }

    assert.Equal(t, bytesFromLightState(s)[0], expected)
}

func TestSetTemperature(t *testing.T) {
    s := LightState{}

    s.SetTemperature(4097)

    expected := [20]byte{
        0x11, 0xff, 0x04, // prefix
        0x9c, // temperature
        0x10, 0x01, // value 4097
        0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
    }

    assert.Equal(t, bytesFromLightState(s)[0], expected)
}

func TestSetMultiple(t *testing.T) {
    s := NewLightState(0, true, 0x80, 5000)

    assert.Equal(t, len(bytesFromLightState(s)), 3)
}

func TestReadPower(t *testing.T) {
    bytes := [6]byte {
        0x11, 0xff, 0x04, //prefix
        0x00, // power
        0x01, // value
        0x00,
    }

    s := lightStateFromBytes(0, bytes)

    assert.Equal(t, s.Power.Value, true)
}

func TestReadBrightness(t *testing.T) {
    bytes := [6]byte {
        0x11, 0xff, 0x04, //prefix
        0x10, // brightness
        0x00, 0x80, // value
    }

    s := lightStateFromBytes(0, bytes)

    assert.Equal(t, s.Brightness.Value, uint8(0x80))
}

func TestReadTemperature(t *testing.T) {
    bytes := [6]byte {
        0x11, 0xff, 0x04, //prefix
        0x20, // temperature
        0x10, 0x01, // value
    }

    s := lightStateFromBytes(0, bytes)

    assert.Equal(t, s.Temperature.Value, uint16(4097))
}
