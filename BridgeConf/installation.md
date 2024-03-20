## Install Webtunnel
- Clone this project go to the `` /main/server `` and run `` go build ``
- You will get `` server `` file and put it into `` /user/local/bin ``, it is name can be webtunnel
- Install Nginx
- `` apt install apt-transport-https ``
- In `` /etc/apt/sources.list.d/tor.list ``  put
```bash
deb     [signed-by=/usr/share/keyrings/tor-archive-keyring.gpg] https://deb.torproject.org/torproject.org <DISTRIBUTION> main
deb-src [signed-by=/usr/share/keyrings/tor-archive-keyring.gpg] https://deb.torproject.org/torproject.org <DISTRIBUTION> main

```
- Find distribution with `` lsb_release -c ``
- At the end run `` wget -qO- https://deb.torproject.org/torproject.org/A3C4F0F979CAA22CDBA8F512EE8CBC9E886DDD89.asc | gpg --dearmor | tee /usr/share/keyrings/tor-archive-keyring.gpg >/dev/null ``
- `` apt update && apt install tor deb.torproject.org-keyring ``
