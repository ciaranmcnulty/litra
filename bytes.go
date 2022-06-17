package litra

const BYTE_PREFIX_1 = 0x11
const BYTE_PREFIX_2 = 0xff
const BYTE_PREFIX_3 = 0x04

const BYTE_POWER_CHANGED = 0x00
const BYTE_BRIGHTNESS_CHANGED = 0x10
const BYTE_TEMPERATURE_CHANGED = 0x20
const BYTE_POWER_SET = 0x1c
const BYTE_BRIGHTNESS_SET = 0x4c
const BYTE_TEMPERATURE_SET = 0x9c

const OFFSET_MESSAGE_TYPE = 3
const OFFSET_POWER_VALUE = 4
const OFFSET_BRIGHTNESS_VALUE = 5
const OFFSET_TEMPERATURE_VALUE_HIGH = 4
const OFFSET_TEMPERATURE_VALUE_LOW = 5

func lightStateFromBytes(id uint64, bytes [6]byte) *LightState {

    if bytes[0] != BYTE_PREFIX_1 || bytes[1] != BYTE_PREFIX_2 || bytes[2] != BYTE_PREFIX_3 {
        return nil
    }

    lightState := LightState{}
    lightState.Id = id

    switch (bytes[OFFSET_MESSAGE_TYPE]) {
        case BYTE_POWER_CHANGED:
            lightState.SetPower(bytes[OFFSET_POWER_VALUE] == 0x01)
        case BYTE_BRIGHTNESS_CHANGED:
            lightState.SetBrightness(bytes[OFFSET_BRIGHTNESS_VALUE])
        case BYTE_TEMPERATURE_CHANGED:
            lightState.SetTemperature(uint16(
                256 * int(bytes[OFFSET_TEMPERATURE_VALUE_HIGH]) + int(bytes[OFFSET_TEMPERATURE_VALUE_LOW],
            )))

       default:
            return nil
    }

    return &lightState
}

func bytesFromLightState (state LightState) [][20]byte {
    var responses [][20]byte

    var prefix [20]byte
    prefix[0] = BYTE_PREFIX_1
    prefix[1] = BYTE_PREFIX_2
    prefix[2] = BYTE_PREFIX_3

    state.ApplyState(
        func (_ uint64, power bool) {
            bytes := prefix
            bytes[OFFSET_MESSAGE_TYPE] = BYTE_POWER_SET
            if (power) {
                bytes[OFFSET_POWER_VALUE] = 0x01
            } else {
                bytes[OFFSET_POWER_VALUE] = 0x00
            }

            responses = append(responses, bytes)
        },
        func (_ uint64, brightness uint8) {
            bytes := prefix
            bytes[OFFSET_MESSAGE_TYPE] = BYTE_BRIGHTNESS_SET
            switch {
                case brightness <=20:
                    bytes[OFFSET_BRIGHTNESS_VALUE] = uint8(20)
                case brightness >=250:
                    bytes[OFFSET_BRIGHTNESS_VALUE] = uint8(250)
                default:
                    bytes[OFFSET_BRIGHTNESS_VALUE] = uint8(brightness)
            }

            responses = append(responses, bytes)
        },
        func (_ uint64, temperature uint16) {
            bytes := prefix
            bytes[OFFSET_MESSAGE_TYPE] = BYTE_TEMPERATURE_SET
            actualTemp := temperature
            switch {
                case temperature <=2700:
                    actualTemp = 2700
                case temperature >=6500:
                    actualTemp = 6500
            }
            bytes[OFFSET_TEMPERATURE_VALUE_LOW] = uint8(actualTemp)
            bytes[OFFSET_TEMPERATURE_VALUE_HIGH] = uint8((actualTemp - uint16(bytes[5]))/256)

            responses = append(responses, bytes)
        },
    )

    return responses
}
