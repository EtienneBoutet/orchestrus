import requests

def is_host_reachable(host):
  """Verify if a host is reachable"""
  response = requests.get("http://" + host + ":5000/")
  if response.status_code == '204':
    return True
  return False