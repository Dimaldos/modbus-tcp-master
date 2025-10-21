from django.shortcuts import render, redirect
from django.http import JsonResponse
from .models import DeviceManager
from .management.commands.modbus_client import ModbusClient
import os
from django.conf import settings

def device_list(request):
    devices = DeviceManager.get_devices()
    return render(request, 'devices/device_list.html', {
        'devices': devices
    })

def create_device(request):
    if request.method == 'POST':
        device_data = {
            'name': request.POST['name'],
            'ip': request.POST['ip'],
            'port': int(request.POST['port']),
            'slave_id': int(request.POST['slave_id']),
            'write_register': int(request.POST['write_register']),
            'write_mask': int(request.POST['write_mask']),
            'read_register': int(request.POST['read_register']),
            'read_mask': int(request.POST['read_mask']),
        }
        DeviceManager.save_device(device_data)
        return redirect('device_list')
    
    return render(request, 'devices/create_device.html')

def edit_device(request, device_name):
    devices = DeviceManager.get_devices()
    device = next((d for d in devices if d['name'] == device_name), None)
    
    if not device:
        return redirect('device_list')
    
    if request.method == 'POST':
        # Удаляем старый файл, если имя изменилось
        if device_name != request.POST['name']:
            old_filename = f"{device_name}.yaml"
            old_filepath = os.path.join(settings.BASE_DIR, 'configs', old_filename)
            if os.path.exists(old_filepath):
                os.remove(old_filepath)
        
        # Сохраняем новую конфигурацию
        device_data = {
            'name': request.POST['name'],
            'ip': request.POST['ip'],
            'port': int(request.POST['port']),
            'slave_id': int(request.POST['slave_id']),
            'write_register': int(request.POST['write_register']),
            'write_mask': int(request.POST['write_mask']),
            'read_register': int(request.POST['read_register']),
            'read_mask': int(request.POST['read_mask']),
        }
        DeviceManager.save_device(device_data)
        return redirect('device_list')
    
    return render(request, 'devices/edit_device.html', {'device': device})

def delete_device(request, device_name):
    if request.method == 'POST':
        filename = f"{device_name}.yaml"
        filepath = os.path.join(settings.BASE_DIR, 'configs', filename)
        if os.path.exists(filepath):
            os.remove(filepath)
        return JsonResponse({'success': True})
    
    return JsonResponse({'success': False, 'error': 'Invalid request method'})

def toggle_device(request, device_name):
    if request.method == 'POST':
        devices = DeviceManager.get_devices()
        device = next((d for d in devices if d['name'] == device_name), None)
        
        if device:
            try:
                # Получаем текущее состояние для переключения
                current_state = ModbusClient.read_state(
                    device['ip'],
                    device['port'],
                    device['slave_id'],
                    device['read_register'],
                    device['read_mask']
                )
                new_state = not current_state
                
                success = ModbusClient.write_state(
                    device['ip'],
                    device['port'],
                    device['slave_id'],
                    device['write_register'],
                    device['write_mask'],
                    new_state
                )
                
                return JsonResponse({'success': success, 'new_state': new_state})
            except Exception as e:
                return JsonResponse({'success': False, 'error': str(e)})
        
        return JsonResponse({'success': False, 'error': 'Device not found'})

def get_device_statuses(request):
    """API endpoint для получения статусов всех устройств"""
    devices = DeviceManager.get_devices()
    statuses = {}
    
    for device in devices:
        try:
            # Пытаемся прочитать состояние - если успешно, значит соединение есть
            state_status = ModbusClient.read_state(
                device['ip'],
                device['port'],
                device['slave_id'],
                device['read_register'],
                device['read_mask']
            )
            connection_status = True
                
        except Exception as e:
            # В случае ошибки - нет соединения
            connection_status = False
            state_status = False
        
        statuses[device['name']] = {
            'connection': connection_status,
            'state': state_status,
            'device': device
        }
    
    return JsonResponse(statuses)