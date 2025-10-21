package cmd

import (
	"flag"
	"fmt"
)

var (
	ip      string
	port    int
	id      int
	read    int
	write   string
	timeout int // Таймаут в секундах
)

func init() {
	flag.StringVar(&ip, "ip", "", "IP-адрес устройства Modbus")
	flag.IntVar(&port, "p", 502, "Порт Modbus TCP")
	flag.IntVar(&id, "id", 1, "Slave ID устройства")
	flag.IntVar(&read, "read", 0, "Адрес регистра для чтения")
	flag.StringVar(&write, "write", "", "Запись в регистр в формате 'адрес:значение'")
	flag.IntVar(&timeout, "timeout", 5, "Таймаут в секундах")
}

func Execute() error {
	flag.Parse()

	if ip == "" {
		return fmt.Errorf("не указан IP-адрес устройства")
	}

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
		return parseAndWrite(write)
	}

	return nil
}

func parseAndWrite(writeStr string) error {
	var addr, value int
	_, err := fmt.Sscanf(writeStr, "%d:%d", &addr, &value)
	if err != nil {
		return fmt.Errorf("неверный формат для -write. Ожидается 'адрес:значение', например '470:1'")
	}

	return writeRegister(addr, value)
}
