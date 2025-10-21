package modbus

import (
	"encoding/binary"
	"fmt"
)

func (c *Client) ReadHoldingRegister(address uint16, unitID uint8) (uint16, error) {
	request := make([]byte, 12)
	
	// MBAP Header
	binary.BigEndian.PutUint16(request[0:], 1) // Transaction ID
	binary.BigEndian.PutUint16(request[2:], 0) // Protocol ID
	binary.BigEndian.PutUint16(request[4:], 6) // Length
	request[6] = unitID                        // Unit ID
	
	// PDU
	request[7] = FuncCodeReadHoldingRegisters  // Function Code
	binary.BigEndian.PutUint16(request[8:], address) // Address
	binary.BigEndian.PutUint16(request[10:], 1)      // Quantity
	
	// Отправка запроса
	if _, err := c.conn.Write(request); err != nil {
		return 0, err
	}
	
	// Чтение ответа
	response := make([]byte, 256)
	n, err := c.conn.Read(response)
	if err != nil {
		return 0, err
	}
	
	if n < 9 {
		return 0, fmt.Errorf("слишком короткий ответ")
	}
	
	// Проверка на ошибку Modbus
	if response[7] != FuncCodeReadHoldingRegisters {
		if response[7]&0x80 != 0 {
			exceptionCode := response[8]
			return 0, fmt.Errorf("Modbus исключение: код %d", exceptionCode)
		}
		return 0, fmt.Errorf("неверный код функции в ответе")
	}
	
	// Извлечение данных
	byteCount := response[8]
	if byteCount != 2 {
		return 0, fmt.Errorf("неверное количество байт данных")
	}
	
	return binary.BigEndian.Uint16(response[9:]), nil
}

func (c *Client) WriteSingleRegister(address uint16, value uint16, unitID uint8) error {
	request := make([]byte, 12)
	
	// MBAP Header
	binary.BigEndian.PutUint16(request[0:], 1) // Transaction ID
	binary.BigEndian.PutUint16(request[2:], 0) // Protocol ID
	binary.BigEndian.PutUint16(request[4:], 6) // Length
	request[6] = unitID                        // Unit ID
	
	// PDU
	request[7] = FuncCodeWriteSingleRegister   // Function Code
	binary.BigEndian.PutUint16(request[8:], address) // Address
	binary.BigEndian.PutUint16(request[10:], value)  // Value
	
	// Отправка запроса
	if _, err := c.conn.Write(request); err != nil {
		return err
	}
	
	// Чтение ответа
	response := make([]byte, 256)
	n, err := c.conn.Read(response)
	if err != nil {
		return err
	}
	
	if n < 12 {
		return fmt.Errorf("слишком короткий ответ")
	}
	
	// Проверка на ошибку Modbus
	if response[7] != FuncCodeWriteSingleRegister {
		if response[7]&0x80 != 0 {
			exceptionCode := response[8]
			return fmt.Errorf("Modbus исключение: код %d", exceptionCode)
		}
		return fmt.Errorf("неверный код функции в ответе")
	}
	
	// Проверка эхо-ответа
	if binary.BigEndian.Uint16(response[8:10]) != address ||
		binary.BigEndian.Uint16(response[10:12]) != value {
		return fmt.Errorf("несоответствие эхо-ответа")
	}
	
	return nil
}