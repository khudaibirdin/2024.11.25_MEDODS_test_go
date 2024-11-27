import requests


headers = {
    "Content-Type": "application/json"
}

# тестирование получения access и refresh
print("тестирование получения access и refresh")
url = "http://localhost:8000/auth/get"
data = {
    "guid": "1",
    "ip": "192.168.0.91"
}
response = requests.post(url, json=data, headers=headers)
print("Статус код:", response.status_code)
print("Ответ:", response.json())
print("\n")

access_token = response.json()["access_token"]
refresh_token = response.json()["refresh_token"]

# тестирование refresh с неверным ip
print("тестирование refresh с неверным ip")
url = "http://localhost:8000/auth/refresh"
data = {
    "ip": "192.168.0.92",
    'access_token': access_token, 
    'refresh_token': refresh_token}

response = requests.post(url, json=data, headers=headers)
print("Статус код:", response.status_code)
print("Ответ:", response.json())
print("\n")

# тестирование с неверным refresh
print("тестирование с неверным refresh")
url = "http://localhost:8000/auth/refresh"
data = {
    "ip": "192.168.0.92",
    'access_token': access_token, 
    'refresh_token': "nub"}

response = requests.post(url, json=data, headers=headers)
print("Статус код:", response.status_code)
print("Ответ:", response.json())
print("\n")

# тестирование с верными данными
print("тестирование с верными данными")
url = "http://localhost:8000/auth/refresh"
data = {
    "ip": "192.168.0.91",
    'access_token': access_token, 
    'refresh_token': refresh_token}

response = requests.post(url, json=data, headers=headers)
print("Статус код:", response.status_code)
print("Ответ:", response.json())