set -e
URLS=("https://github.com/tasselx/free-gemini/releases/download/latest/")
url=${URLS[0]}
lc_type=$(echo $LC_CTYPE | cut -c 1-2)
if [ -z $lc_type ] || [ "$lc_type" = "UT" ]; then
  lc_type=$(echo $LANG | cut -c 1-2)
fi

if [ "$lc_type" = "zh" ]; then
  echo "正在安装..."
else
  echo "Installing..."
fi

for url0 in ${URLS[@]}; do
  if curl -Is --connect-timeout 4 "$url0" | grep -q "HTTP/1.1 404"; then
    url=$url0
    break
  fi
done

os_name=$(uname -s | tr '[:upper:]' '[:lower:]')
if [[ $os_name == *"mingw"* ]]; then
  os_name="windows"
fi
raw_hw_name=$(uname -m)
case "$raw_hw_name" in
"amd64")
  hw_name="amd64"
  ;;
"x86_64")
  hw_name="amd64"
  ;;
"arm64")
  hw_name="arm64"
  ;;
"aarch64")
  hw_name="arm64"
  ;;
"i686")
  hw_name="386"
  ;;
"armv7l")
  hw_name="arm"
  ;;
*)
  echo "Unsupported hardware: $raw_hw_name"
  exit 1
  ;;
esac

if [ "$lc_type" = "zh" ]; then
  echo "当前系统为 ${os_name} ${hw_name}"
else
  echo "Current system is ${os_name} ${hw_name}"
fi


# 如果是mac或者linux系统
if [[ $os_name == "darwin" || $os_name == "linux" ]]; then
  # 打印下载地址
  echo "下载地址/Download URL: ${url}free-gemini_${os_name}_${hw_name}"
  if [ "$lc_type" = "zh" ]; then
    echo "请输入开机密码"
  else
    echo "Please enter the boot password"
  fi;
  # 停掉正在运行的free-gemini
  pkill free-gemini || true
  # 安装
  sudo mkdir -p /usr/local/bin
  sudo curl -Lko /usr/local/bin/free-gemini ${url}/free-gemini_${os_name}_${hw_name}
  sudo chmod +x /usr/local/bin/free-gemini
  if [ "$lc_type" = "zh" ]; then
    echo "安装完成！自动运行；下次可直接输入 free-gemini 并回车来运行程序"
  else
    echo "Installation completed! Automatically run; you can run the program by entering free-gemini and pressing Enter next time"
  fi;

  echo ""
  free-gemini
fi;
# 如果是windows系统
if [[ $os_name == "windows" ]]; then
  # 打印下载地址
  echo "下载地址/Download URL: ${url}free-gemini_${os_name}_${hw_name}.exe"
  # 停掉正在运行free-gemini
  taskkill -f -im free-gemini.exe || true
  # 安装
  curl -Lko ${USERPROFILE}/Desktop/free-gemini.exe ${url}/free-gemini_${os_name}_${hw_name}.exe
  if [ "$lc_type" = "zh" ]; then
    echo "安装完成！自动运行; 下次可直接输入 ./free-gemini.exe 并回车来运行程序"
    echo "运行后如果360等杀毒软件误报木马，添加信任后，重新输入./free-gemini.exe 并回车来运行程序"
  else
    echo "Installation completed! Automatically run; you can run the program by entering ./free-gemini.exe and press Enter next time"
    echo "After running, if 360 antivirus software reports a Trojan horse, add trust, and then re-enter ./free-gemini.exe and press Enter to run the program"
  fi

  echo ""
  chmod +x ${USERPROFILE}/Desktop/free-gemini.exe
  ${USERPROFILE}/Desktop/free-gemini.exe
fi
