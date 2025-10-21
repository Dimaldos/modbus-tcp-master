from django.core.management.base import BaseCommand
import subprocess
import re
import os
import platform

class ModbusClient:
    @staticmethod
    def read_state(ip, port, slave_id, read_register, read_mask):
        # Определяем путь к mb.exe в зависимости от ОС
        if platform.system() == "Windows":
            mb_path = "../modbus-cli/mb.exe"
        else:
            mb_path = "./mb"
        
        # Формируем команду для чтения регистра
        command = [
            mb_path,
            '-read', str(read_register),
            '-ip', ip,
            '-p', str(port),
            '-id', str(slave_id)
        ]
        
        try:
            # Проверяем существование файла mb.exe
            if not os.path.exists(mb_path):
                raise Exception(f"CLI tool not found at: {mb_path}")
            
            # Выполняем команду (без указания кодировки)
            result = subprocess.run(
                command, 
                capture_output=True, 
                text=True, 
                timeout=10,
                shell=False
            )
            
            # Проверяем успешность выполнения
            if result.returncode != 0:
                error_msg = result.stderr.strip() if result.stderr else "Unknown error"
                raise Exception(f"CLI error (code {result.returncode}): {error_msg}")
            
            # Парсим ответ (теперь на английском)
            output = result.stdout.strip()
            
            # Ищем формат "Register X: Y"
            match = re.search(r'Register \d+: (\d+)', output)
            if not match:
                # Альтернативный формат
                match = re.search(r'register \d+: (\d+)', output, re.IGNORECASE)
            if not match:
                # Просто ищем число
                match = re.search(r'(\d+)', output)
            
            if match:
                value = int(match.group(1))
                # Применяем битовую маску
                masked_value = value & read_mask
                # Преобразуем в булево значение
                state = masked_value != 0
                return state
            else:
                raise Exception(f"Unexpected response format: '{output}'")
                
        except subprocess.TimeoutExpired:
            raise Exception("CLI command timeout")
        except FileNotFoundError:
            raise Exception(f"CLI tool not found: {mb_path}")
        except Exception as e:
            raise Exception(f"Read failed: {str(e)}")

    @staticmethod
    def write_state(ip, port, slave_id, write_register, write_mask, state):
        # Определяем путь к mb.exe в зависимости от ОС
        if platform.system() == "Windows":
            mb_path = "../modbus-cli/mb.exe"
        else:
            mb_path = "./mb"
            
        # Преобразуем булево состояние в значение для записи
        value_to_write = write_mask if state else 0
        
        # Формируем команду для записи в регистр
        command = [
            mb_path,
            '-write', f"{write_register}:{value_to_write}",
            '-ip', ip,
            '-p', str(port),
            '-id', str(slave_id)
        ]
        
        try:
            # Проверяем существование файла mb.exe
            if not os.path.exists(mb_path):
                raise Exception(f"CLI tool not found at: {mb_path}")
            
            # Выполняем команду (без указания кодировки)
            result = subprocess.run(
                command, 
                capture_output=True, 
                text=True, 
                timeout=10,
                shell=False
            )
            
            # Проверяем успешность выполнения
            if result.returncode != 0:
                error_msg = result.stderr.strip() if result.stderr else "Unknown error"
                raise Exception(f"CLI error (code {result.returncode}): {error_msg}")
            
            # Проверяем ответ (теперь на английском)
            output = result.stdout.strip()
            
            # Проверяем успешный ответ на английском
            if "Successfully written" in output:
                return True
            else:
                raise Exception(f"Write failed: '{output}'")
                
        except subprocess.TimeoutExpired:
            raise Exception("CLI command timeout")
        except FileNotFoundError:
            raise Exception(f"CLI tool not found: {mb_path}")
        except Exception as e:
            raise Exception(f"Write failed: {str(e)}")