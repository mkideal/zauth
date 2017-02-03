package model

func Init() error {
	// TODO
	return nil
}

func ValidateClint(client *Client, clientSecret string) bool {
	return clientSecret == client.Secret
}
