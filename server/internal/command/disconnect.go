package command

const DisconnectMessage = "Disconnected!"

type DisconnectCommand struct {
}

func (d DisconnectCommand) Do() error {
	panic(DisconnectMessage)
}

