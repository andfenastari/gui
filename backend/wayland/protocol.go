package wayland

const DisplayId ObjectId = 1

const (
	OpDisplaySync        Opcode = 0
	OpDisplayGetRegistry        = 1

	OpDisplayError = 0

	OpRegistryBind         = 0
	OpRegistryGlobal       = 0
	OpRegistryGlobalRemove = 1

	OpCallbackDone = 0

	OpCompositorCreateSurface = 0
	OpCompositonCreateRegion  = 1
)

type ObjectId uint32
type Opcode uint16

type DisplaySync struct {
	Callback ObjectId
}

type DisplayGetRegistry struct {
	Registry ObjectId
}

type DisplayError struct {
	ObjectId ObjectId
	Code     uint32
	Message  string
}

type RegistryBind struct {
	Name uint32
	Interface string
	Version uint32
	Id   ObjectId
}

type RegistryGlobal struct {
	Name      uint32
	Interface string
	Version   uint32
}

type CallbackDone struct {
	CallbackData uint32
}

type CompositorCreateSurface struct {
	Id ObjectId
}

type CompositorCreateRegion struct {
	Id ObjectId
}
