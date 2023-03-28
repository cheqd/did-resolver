package utils

func PrintConfig() error {
	config := MustLoadConfig()
	configJson := config.MustMarshalJson()

	println(configJson)

	return nil
}
