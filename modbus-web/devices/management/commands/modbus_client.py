from django.core.management.base import BaseCommand
import random
import time

class ModbusClient:
    @staticmethod
    def read_state(ip, port, slave_id, read_register, read_mask):
        # Заглушка для чтения состояния
        print(f"Reading state from {ip}:{port} (slave {slave_id})")
        
        # Имитация случайных ошибок чтения (10% вероятность)
        if random.random() < 0.1:
            raise Exception("Read error")
            
        # Имитация задержки сети
        time.sleep(0.05)
        
        # Возвращаем случайное состояние для демонстрации
        return True

    @staticmethod
    def write_state(ip, port, slave_id, write_register, write_mask, state):
        # Заглушка для записи состояния
        print(f"Writing state {state} to {ip}:{port} (slave {slave_id})")
        
        # Имитация случайных ошибок записи (5% вероятность)
        if random.random() < 0.05:
            raise Exception("Write error")
            
        # Имитация задержки сети
        time.sleep(0.05)
        return True