package cmd

import (
	"encoding/binary"
	"fmt"
	"net"
)

func readRegister(register int) error {
	addr := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("ошибка подключения к %s: %v", addr, err)
	}
	defer conn.Close()

	request := make([]byte, 12)
	binary.BigEndian.PutUint16(request[0:], 1)                // Transaction ID
	binary.BigEndian.PutUint16(request[2:], 0)                // Protocol
	binary.BigEndian.PutUint16(request[4:], 6)                // Length
	request[6] = byte(id)                                     // Unit ID
	request[7] = 0x03                                         // Function Code
	binary.BigEndian.PutUint16(request[8:], uint16(register)) // Address
	binary.BigEndian.PutUint16(request[10:], 1)               // Quantity

	fmt.Print("request: ")
	fmt.Println(request)

	_, err = conn.Write(request)
	if err != nil {
		return fmt.Errorf("ошибка отправки запроса: %v", err)
	}

	response := make([]byte, 256)
	n, err := conn.Read(response)
	if err != nil {
		return fmt.Errorf("ошибка чтения ответа: %v", err)
	}

	if n < 9 {
		return fmt.Errorf("слишком короткий ответ от устройства")
	}

	fmt.Print("response: ")
	fmt.Println(response)

	if response[7] != 0x03 {
		if response[7] == 0x83 {
			exceptionCode := response[8]
			return fmt.Errorf("modbus исключение: код %d", exceptionCode)
		}
		return fmt.Errorf("неверный код функции в ответе: %d", response[7])
	}

	data := binary.BigEndian.Uint16(response[9:])
	fmt.Printf("Регистр %d: %d\n", register, data)
	return nil
}
