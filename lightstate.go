package litra

const ALL_LIGHTS=0

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

func NewLightState (id uint64, p bool, b uint8, t uint16) LightState {
    return LightState {
        id,
        &Power{p},
        &Brightness{b},
        &Temperature{t},
    }
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

func (s *LightState) IsEmpty () bool {
    return s.Power == nil && s.Brightness == nil && s.Temperature == nil
}
