import requests, json

async def displayData(sheet:str="", char:str="", elem:str=""):
    resp = requests.get(f'http://localhost:8080/api/{sheet}?character={char}&element={elem}')
    data = json.loads(resp.text)
    print(data[0])

displayData(char="イレイナ")