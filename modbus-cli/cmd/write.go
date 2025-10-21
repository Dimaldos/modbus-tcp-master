package cmd

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

func writeRegister(register, value int) error {
	addr := net.JoinHostPort(ip, fmt.Sprintf("%d", port))
	conn, err := net.DialTimeout("tcp", addr, time.Duration(timeout)*time.Second)
	if err != nil {
		return fmt.Errorf("ошибка подключения к %s: %v", addr, err)
	}
	defer conn.Close()

	// Устанавливаем таймаут для операций чтения/записи
	conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))

	request := make([]byte, 12)
	binary.BigEndian.PutUint16(request[0:], 1)                // Transaction ID
	binary.BigEndian.PutUint16(request[2:], 0)                // Protocol
	binary.BigEndian.PutUint16(request[4:], 6)                // Length
	request[6] = byte(id)                                     // Unit ID
	request[7] = 0x06                                         // Function Code
	binary.BigEndian.PutUint16(request[8:], uint16(register)) // Address
	binary.BigEndian.PutUint16(request[10:], uint16(value))   // Value

	_, err = conn.Write(request)
	if err != nil {
		return fmt.Errorf("ошибка отправки запроса: %v", err)
	}

	response := make([]byte, 256)
	n, err := conn.Read(response)
	if err != nil {
		// Проверяем, это таймаут или другая ошибка
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return fmt.Errorf("таймаут при чтении ответа от устройства")
		}
		return fmt.Errorf("ошибка чтения ответа: %v", err)
	}

	if n < 12 {
		return fmt.Errorf("слишком короткий ответ от устройства")
	}

	if response[7] != 0x06 {
		if response[7] == 0x86 {
			exceptionCode := response[8]
			return fmt.Errorf("modbus исключение: код %d", exceptionCode)
		}
		return fmt.Errorf("неверный код функции в ответе: %d", response[7])
	}

	respAddr := binary.BigEndian.Uint16(response[8:10])
	respValue := binary.BigEndian.Uint16(response[10:12])

	if respAddr != uint16(register) || respValue != uint16(value) {
		return fmt.Errorf("несоответствие эхо-ответа: ожидался адрес %d, значение %d, получен адрес %d, значение %d",
			register, value, respAddr, respValue)
	}

	fmt.Printf("Successfully written %d to register %d\n", value, register)
	return nil
}
