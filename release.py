import os
import shutil

from webdav3 import client

FrontendPath = os.environ.get("FRONTEND_PATH")
BackendPath = os.environ.get("BACKEND_PATH")


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
    os.system(os.path.join(BackendPath, "build.bat"))
    print("OK")


def BuildFrontend():
    print("构建前端中...", end="")
    os.chdir(FrontendPath)
    os.system(os.path.join(FrontendPath, "build.bat"))
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
        'webdav_hostname': os.environ.get("WEBDAV_HOST"),
        'webdav_login': os.environ.get("WEBDAV_USER"),
        'webdav_password': os.environ.get("WEBDAV_PASSWD")
    }
    webdav = client.Client(options)

    webdav.upload_sync(remote_path="/CDN/Project/LoongPanel/bin/dist.zip",
                       local_path=os.path.join(BackendPath, "bin", "dist.zip"))
    webdav.upload_sync(remote_path="/CDN/Project/LoongPanel/bin/scripts.zip",
                       local_path=os.path.join(BackendPath, "bin", "scripts.zip"))
    print("上传完成")


if __name__ == '__main__':
    Clean()
    BuildBackend()
    BuildFrontend()
    PackScript()
    UploadWebDav()
    print("发布完成")
    print("https://cdn1.d5v.cc/CDN/Project/LoongPanel/bin/dist.zip")
    print("https://cdn1.d5v.cc/CDN/Project/LoongPanel/bin/scripts.zip")
