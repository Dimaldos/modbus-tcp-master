import os
import yaml
from django.conf import settings

class DeviceManager:
    @staticmethod
    def get_devices():
        devices = []
        config_dir = os.path.join(settings.BASE_DIR, 'configs')
        if not os.path.exists(config_dir):
            os.makedirs(config_dir)
            return devices
            
        for filename in os.listdir(config_dir):
            if filename.endswith('.yaml') or filename.endswith('.yml'):
                with open(os.path.join(config_dir, filename), 'r') as f:
                    device_data = yaml.safe_load(f)
                    devices.append(device_data)
        return devices

    @staticmethod
    def save_device(device_data):
        filename = f"{device_data['name']}.yaml"
        config_dir = os.path.join(settings.BASE_DIR, 'configs')
        if not os.path.exists(config_dir):
            os.makedirs(config_dir)
            
        filepath = os.path.join(config_dir, filename)
        with open(filepath, 'w') as f:
            yaml.dump(device_data, f)