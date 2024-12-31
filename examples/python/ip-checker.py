import os
import requests

api_url = os.getenv('IP_API_ENDPOINT', 'https://api.ipify.org')

def get_public_ip():
    response = requests.get(api_url)
    if response.status_code == 200:
        return response.text.strip()
    else:
        raise RuntimeError(f"Unable to fetch IP. Status code: {response.status_code}")

if __name__ == "__main__":
    try:
        ip_address = get_public_ip()
        print(ip_address)
    except Exception as e:
        raise e
