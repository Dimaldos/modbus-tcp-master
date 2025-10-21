package modbus

// Modbus функции
const (
	FuncCodeReadHoldingRegisters  = 0x03
	FuncCodeWriteSingleRegister   = 0x06
	FuncCodeWriteMultipleRegisters = 0x10
)

// MBAP Header структура
type MBAPHeader struct {
	TransactionID uint16
	ProtocolID    uint16
	Length        uint16
	UnitID        uint8
}

// Modbus исключения
const (
	ExceptionIllegalFunction     = 0x01
	ExceptionIllegalDataAddress  = 0x02
	ExceptionIllegalDataValue    = 0x03
	ExceptionServerDeviceFailure = 0x04
)