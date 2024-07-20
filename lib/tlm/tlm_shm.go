package tlm

import (
	"github.com/AarC10/GSW-V2/lib/ipc"
	"github.com/AarC10/GSW-V2/proc"
)

func TlmServiceInit(packet proc.TelemetryPacket) (ipc.IpcServiceSide, error) {
	var tlmService ipc.ShmServiceSide
	err := tlmService.Setup(packet)
	if err != nil {
		return nil, err
	}

	return &tlmService, nil
}

func TlmClientInit(packet proc.TelemetryPacket) (ipc.IpcClientSide, error) {
	var tlmClient ipc.ShmClientSide
	err := tlmClient.Setup(packet)
	if err != nil {
		return nil, err
	}

	return &tlmClient, nil
}
