import requests, json

NETBOX_URL = "http://10.39.70.169:8000/api/"
HEADERS = {
    "Authorization": "Token 53bccdc4d527945c9b24b0a5cc5a558e212b3def",
    "Content-Type": "application/json",
}

def get_device_details(device_name):
    # Get device details by name
    url = f"{NETBOX_URL}dcim/devices/?name={device_name}"
    response = requests.get(url, headers=HEADERS)

    if response.status_code == 200:
        device_details = response.json()["results"][0]
        return device_details
    else:
        print(f"Error: {response.status_code}")
        return None
    
def get_devices_details():
    # Get device details by name
    url = f"{NETBOX_URL}dcim/devices"
    response = requests.get(url, headers=HEADERS)

    if response.status_code == 200:        
        device_details = response.json()["results"]
        device_lists = [device['name'] for device in device_details]
        return device_lists
    else:
        print(f"Error: {response.status_code}")
        return None

def get_interfaces_details():
    url = f"{NETBOX_URL}dcim/interfaces"
    response = requests.get(url, headers=HEADERS)

    if response.status_code == 200:        
        interface_details = response.json()["results"]
        # interface_lists = [device['name'] for device in interface_details]
        return interface_details
    else:
        print(f"Error: {response.status_code}")
        return None

def main():
    device_lists = get_devices_details()
    # device_name = input("Enter the device name: ")
    list_of_device_dicts = []
    for device_name in device_lists:
        interface_dict = []
        device_details = get_device_details(device_name)
        if device_details['interface_count'] > 0:
            interface_details = get_interfaces_details()
            for interface in interface_details:
                if interface['device']['name'] == device_name:
                    interface_dict.append({"name": interface['name'], "speed": interface['speed'],})                   
        if device_details:
            if (device_details['custom_fields']['State']):
                pass
            else:
                device_details['custom_fields']['State'] = "None"
            if interface_dict == []:
                interface_dict = [{}]
            device_data = {'Id': device_details['id'], 'Name': device_details['name'], "DeviceType": device_details['device_type']['model'], 'Manufacturer': device_details['device_type']['manufacturer']['name'], 'State': device_details['custom_fields']['State'], "interfaces": interface_dict}
            list_of_device_dicts.append(device_data)
        else:
            print(f"Device '{device_name}' not found.")
    print(json.dumps(list_of_device_dicts, indent=2))
if __name__ == "__main__":
    main()
