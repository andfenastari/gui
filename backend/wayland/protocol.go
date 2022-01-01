package wayland

const DisplayId ObjectId = 1

const (
	OpDisplaySync        Opcode = 0
	OpDisplayGetRegistry        = 1

	OpRegistryBind         = 0
	OpRegistryGlobal       = 0
	OpRegistryGlobalRemove = 1

	OpCallbackDone = 0
)

type ObjectId uint32
type Opcode uint16

type RequestDisplaySync struct {
	Callback ObjectId
}

type RequestDisplayGetRegistry struct {
	Registry ObjectId
}

type EventRegistryGlobal struct {
	Name      uint32
	Interface string
	Version   uint32
}

type EventCallbackDone struct {
	CallbackData uint32
}
