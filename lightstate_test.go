package litra

import (
    "testing"
    "gotest.tools/assert"
)

func TestInit(t *testing.T) {
    s := NewLightState(1234, true, 0x80, 5000)

    assert.Equal(t, s.Id, uint64(1234))
    assert.Equal(t, s.Power.Value, true)
    assert.Equal(t, s.Brightness.Value, uint8(0x80))
    assert.Equal(t, s.Temperature.Value, uint16(5000))
}

func TestSetters(t *testing.T) {
    s := LightState{}

    s.SetPower(true)
    s.SetBrightness(0x80)
    s.SetTemperature(5000)

    assert.Equal(t, s.Power.Value, true)
    assert.Equal(t, s.Brightness.Value, uint8(0x80))
    assert.Equal(t, s.Temperature.Value, uint16(5000))
}

func TestApplyStateDoesNothingIfNoFieldsToApply(t *testing.T) {

    s := LightState{}

    s.ApplyState(
        func(uint64, bool) { t.Fail() },
        func(uint64, uint8) { t.Fail() },
        func(uint64, uint16) { t.Fail() },
    )
}

func TestApplyStateTriggersCallbacksWhenThereAreFieldsToApply(t *testing.T) {

    s := NewLightState(1234, true, 0x80, 5000)

    counter := 0

    s.ApplyState(
        func(uint64, bool) { counter++ },
        func(uint64, uint8) { counter++ },
        func(uint64, uint16) { counter++ },
    )

    assert.Equal(t, counter, 3)
}

func TestIsEmptyWhenNothingIsSet(t *testing.T) {

    s := LightState{}

    assert.Assert(t, s.IsEmpty())
}

func TestIsNotEmptyWhenSomethingIsSet(t *testing.T) {
    p := LightState{}
    p.SetPower(true)

    b := LightState{}
    b.SetBrightness(0x80)

    w := LightState{}
    w.SetTemperature(5000)

    assert.Assert(t, !p.IsEmpty())
    assert.Assert(t, !b.IsEmpty())
    assert.Assert(t, !w.IsEmpty())
}
