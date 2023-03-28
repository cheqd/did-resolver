package utils

func printConfig() error {
	config := MustLoadConfig()
	configJson := config.MustMarshalJson()

	println(configJson)

	return nil
}
