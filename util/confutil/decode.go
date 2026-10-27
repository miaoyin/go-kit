package confutil

import "github.com/mitchellh/mapstructure"


//Decode map转struct
func Decode(m map[string]any, out any) error {
    decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
        TagName: "mapstructure",
        DecodeHook: mapstructure.ComposeDecodeHookFunc(
            mapstructure.StringToTimeDurationHookFunc(),
        ),
        Result: out,
    })
    if err != nil {
        return err
    }
    return decoder.Decode(m)
}
