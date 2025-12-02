import requests
import json

URL = "http://localhost:8080/articles/"   # adjust your route if different
TOKEN = "ZHTORVMVDJDJCFH2FHXQ2RR2R3GI77GQQW7ANSSJGSUTC7QEUIDQ"  # replace with a valid token

headers = {
    "Authorization": f"Bearer {TOKEN}",
    "Content-Type": "application/json",
}

payload = {
    "title": "How to Scale Logistics SaaS in 2025",
    "image": "https://example.com/article-cover.jpg",
    "paragraphs": [
        {
            "headline": "Introduction",
            "body": "In this article, we explore best practices for scaling logistics-tech platforms.",
            "order": 1
        },
        {
            "headline": "Understanding the Market",
            "body": "The trucking and logistics sector is rapidly modernizing with data-driven solutions.",
            "order": 2
        },
        {
            "headline": "Building Systems that Scale",
            "body": "Designing reliable backend architecture is essential for growth.",
            "order": 3
        }
    ]
}

response = requests.post(URL, headers=headers, json=payload)

print("Status Code:", response.status_code)
try:
    print(json.dumps(response.json(), indent=4))
except:
    print(response.text)
