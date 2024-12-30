import os
import shutil

from webdav3 import client

WEBDAV_HOST = "https://webdav-1817443180.pd1.123pan.cn/webdav"
WEBDAV_PASSWD = "1u9302f1000d228h8d9ijga3y05wc2fw"
WEBDAV_USER = "17774369488"

if os.name == "nt":
    FrontendPath = os.environ.get("FRONTEND_PATH")
    BackendPath = os.environ.get("BACKEND_PATH")
else:
    BackendPath = "/mnt/d/code/golang/LoongPanel"
    FrontendPath = "/mnt/d/code/WebstormProjects/loong_panel_react"


def Clean():
    print("清理目录中...", end="")
    shutil.rmtree(os.path.join(BackendPath, "bin"), ignore_errors=True)
    shutil.rmtree(os.path.join(FrontendPath, "dist"), ignore_errors=True)
    os.makedirs(os.path.join(BackendPath, "bin"), exist_ok=True)
    os.makedirs(os.path.join(FrontendPath, "dist"), exist_ok=True)
    print("OK")


def BuildBackend():
    print("构建后端中...", end="")
    os.chdir(BackendPath)
    if os.name == "nt":
        os.system(os.path.join(BackendPath, "build.bat"))
    else:
        os.system(os.path.join(BackendPath, "build.sh"))
    print("OK")


def BuildFrontend():
    print("构建前端中...", end="")
    os.chdir(FrontendPath)
    # if os.name == "nt":
    #     os.system(os.path.join(BackendPath, "build.bat"))
    # else:
    #     os.system(os.path.join(FrontendPath, "build.sh"))
    os.system(f"7z a -tzip {os.path.join(FrontendPath, 'dist.zip')} {os.path.join(FrontendPath, 'dist', '*')}")
    os.rename(os.path.join(FrontendPath, "dist.zip"), os.path.join(BackendPath, "bin", "dist.zip"))
    print("OK")


def PackScript():
    print("打包脚本中...", end="")
    os.chdir(BackendPath)
    os.system(
        f"7z a -tzip {os.path.join(BackendPath, 'bin', 'scripts.zip')} {os.path.join(BackendPath, 'Script', '*')}")
    print("OK")


def UploadWebDav():
    options = {
        'webdav_hostname': WEBDAV_HOST,
        'webdav_login': WEBDAV_USER,
        'webdav_password': WEBDAV_PASSWD
    }
    webdav = client.Client(options)

    print("上传dist.zip")
    webdav.upload_sync(remote_path="/CDN/Project/LoongPanel/bin/dist.zip",
                       local_path=os.path.join(BackendPath, "bin", "dist.zip"))
    print("上传script.zip")
    webdav.upload_sync(remote_path="/CDN/Project/LoongPanel/bin/scripts.zip",
                       local_path=os.path.join(BackendPath, "bin", "scripts.zip"))
    print("上传二进制文件")
    if os.path.exists(os.path.join(BackendPath, "bin", "LoongPanel-linux-amd64")):
        webdav.upload_sync(remote_path="/CDN/Project/LoongPanel/bin/LoongPanel-linux-amd64",
                           local_path=os.path.join(BackendPath, "bin", "LoongPanel-linux-amd64"))
    if os.path.exists(os.path.join(BackendPath, "bin", "LoongPanel-linux-loong64")):
        webdav.upload_sync(remote_path="/CDN/Project/LoongPanel/bin/LoongPanel-linux-loong64",
                           local_path=os.path.join(BackendPath, "bin", "LoongPanel-linux-loong64"))
    print("上传完成")


if __name__ == '__main__':
    Clean()
    BuildBackend()
    BuildFrontend()
    PackScript()
    UploadWebDav()
    print("发布完成")
