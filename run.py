import requests
import threading
import time
import os

# 测试后端是否可用
res = requests.get("http://127.0.0.1:8080/api/v1/files/dir")
if res.status_code != 200:
    print("后端不可用")
    exit()
request_size = len(res.content)

def job():
    requests.get("http://127.0.0.1:8080/api/v1/files/dir")

start = time.time()
duration = 60  # 持续时间为60秒
end_time = start + duration
request_count = 0

# 获取初始CPU占用
cpu_usage_start = os.getloadavg()[0]
cpu_usage_peak = cpu_usage_start

print("开始持续测试接口:")

while time.time() < end_time:
    threading.Thread(target=job).start()
    request_count += 1
    # 获取当前CPU占用并更新峰值
    current_cpu_usage = os.getloadavg()[0]
    if current_cpu_usage > cpu_usage_peak:
        cpu_usage_peak = current_cpu_usage

end = time.time()
print("结束测试:")
print(f"总耗时: {end - start} 秒")
print(f"总请求数: {request_count} 个")
print(f"平均每秒请求数: {request_count / duration} 个")
print(f"平均请求大小: {request_size} 字节")
print(f"初始CPU占用: {cpu_usage_start}")
print(f"峰值CPU占用: {cpu_usage_peak}")