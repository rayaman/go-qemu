# Base Image Setups

Defaults
---
User: `baseimguser`
Pass: `baseimgpass`

All base images should use these credentials. Once a new image is built off of these they should create new users and delete the old base user.

- Install [openssh-client/server](https://www.tecmint.com/install-openssh-server-in-linux/)
- Install net-tools (apt install net-tools)
- Install inetutils-ping (apt install inetutils-ping)
- Install firewall (apt install ufw)
- Install nmcli (apt install network-manager)
- Install yq (apt install yq)

Configure TAP networking:
---
refer to [windows setup](https://wonghoi.humgar.com/blog/2021/05/03/qemu-for-windows-host-quirks/)



Base Images To support:
---
- Ubuntu/Debian
- Rocky/AlmaLinux
Strech goal
- Windows
- MacOS

`qemu-system-x86_64.exe -m 4096 -accel whpx -boot c -smp 4 -nic tap,ifname=qemu-tap -hda .\imgs\UbuntuTest.img -nographic`