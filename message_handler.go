package main

func (client *Client) Send(message string) error {
	return client.SendBytes([]byte(message))
}

func (client *Client) SendBytes(message []byte) error {
	_, err := client.connection.Write(message)

	if err != nil {
		client.connection.Close()
		client.Server.onClientDisconnect(client, err)
	}

	return err
}
