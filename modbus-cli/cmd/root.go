package cmd

import (
	"flag"
	"fmt"
)

var (
	ip    string
	port  int
	id    int
	read  int
	write string
)

func init() {
	flag.StringVar(&ip, "ip", "", "IP-адрес устройства Modbus")
	flag.IntVar(&port, "p", 502, "Порт Modbus TCP")
	flag.IntVar(&id, "id", 1, "Slave ID устройства")
	flag.IntVar(&read, "read", 0, "Адрес регистра для чтения")
	flag.StringVar(&write, "write", "", "Запись в регистр в формате 'адрес:значение'")
}

func Execute() error {
	flag.Parse()

	if ip == "" {
		return fmt.Errorf("не указан IP-адрес устройства")
	}

	// Проверяем, что указана ровно одна операция
	if read != 0 && write != "" {
		return fmt.Errorf("можно указать только одну операцию: -read или -write")
	}

	if read == 0 && write == "" {
		return fmt.Errorf("не указана операция (используйте -read или -write)")
	}

	if read != 0 {
		return readRegister(read)
	}

	if write != "" {
		// Парсим строку вида "адрес:значение"
		return parseAndWrite(write)
	}

	return nil
}

func parseAndWrite(writeStr string) error {
	// Ожидаем формат "адрес:значение"
	var addr, value int
	_, err := fmt.Sscanf(writeStr, "%d:%d", &addr, &value)
	if err != nil {
		return fmt.Errorf("неверный формат для -write. Ожидается 'адрес:значение', например '470:1'")
	}

	return writeRegister(addr, value)
}

func usage() {
	fmt.Println("Использование:")
	fmt.Println("  Чтение:  mb -ip <ip> -read <адрес_регистра> [-p <port>] [-id <id>]")
	fmt.Println("  Запись:  mb -ip <ip> -write <адрес:значение> [-p <port>] [-id <id>]")
	fmt.Println("")
	fmt.Println("Примеры:")
	fmt.Println("  mb -ip 192.168.1.99 -read 470")
	fmt.Println("  mb -ip 192.168.1.99 -write 470:1")
	fmt.Println("  mb -ip 192.168.1.99 -p 502 -id 1 -write 470:1")
}
