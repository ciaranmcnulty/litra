package main

import (
    "github.com/ciarancmnulty/litra-go"
    "log"
    "time"
)

func main() {
    litra.Start(logDeviceStatus, logLightState)

    for {
        time.Sleep(time.Second)
        litra.Request(litra.NewLightState(litra.ALL_LIGHTS, true, 0x10, 5000))
        time.Sleep(time.Second)
        litra.Request(litra.NewLightState(litra.ALL_LIGHTS, false, 0x10, 5000))
    }
}

func logDeviceStatus (status litra.Device) {
     if (status.Connected) {
         log.Printf("Device %d connected", status.Id);
     } else {
         log.Printf("Device %d disconnected", status.Id);
     }
}

func updatePower (id uint64, power bool) {
    if (power) {
        log.Printf("Device %d On", id)
    } else {
        log.Printf("Device %d Off", id)
    }
}

func updateBrightness (id uint64, brightness uint8) {
    log.Printf("Device %d set to brightness %d", id, brightness)
}

func updateTemperature (id uint64, temperature uint16) {
    log.Printf("Device %d set to temperature %d", id, temperature)
}

func logLightState (status litra.LightState) {
    status.ApplyState(updatePower, updateBrightness, updateTemperature)
}
