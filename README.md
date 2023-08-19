# MPTCPTunnelTest
Simple MPTCP Tunnel
## How to Use?
- run go build in this floder and then sudo MPTCPTun-Go -h
- Only Linux kernels >= v5.16 CAN USE MPTCP, Windows is NOT supported
- BBR can improve the Network Speed
## About Performance
- In any case, it is similar to ordinary TCP, the only advantage is failover
- If you want to improve your Connection in high Loss, Visit https://github.com/apernet/hysteria
## Stringline inspired from https://github.com/ginuerzh/gost

## Chinese(简体中文)
- Linux内核必须大于5.16 否则无法使用
- 建议配合BBR+FQ或者BBR+Cake使用
- 编译后 sudo MPTCPTun-Go -h 查看使用方法
- 如果追求速度 请使用hysteria(https://github.com/apernet/hysteria) MPTCP拯救不了高丢包 只能拯救容易断连,但有备用方案的网络
- 基本环境跑起来的速度和TCP基本一致 并未发现优势 甚至慢了点
- Stringline 方法从 https://github.com/ginuerzh/gost 获取