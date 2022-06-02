import os, sys
import threading

def run_node(ip: str, port: int):
  os.system(f"go run cmd/main.go {ip} {port}")

def create_node(ip: str, port: int) -> threading.Thread:
  os.system(f"bash setup.sh {ip}:{port}")

  thread = threading.Thread(target=run_node, args=(ip, port,))
  thread.start()

  return thread

if __name__ == "__main__":
  if len(sys.argv) != 2:
    print("Not enough arguments: ./local_development <number_nodes>")
    exit(1)

  number = int(sys.argv[1])

  ip = "127.0.0.1"
  port = 8000

  threads = []

  for i in range(number):
    t = create_node(ip, port)
    threads.append(t)

    port += 2
