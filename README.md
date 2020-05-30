

# cronshell

多台服务器相同 `crontab` 只执行一次

# crontab

```conf
 # cat /etc/cronshell.conf
[Log]
logfile=/var/log/cronshell.log

[Redis]
host=192.168.15.100
port=6379
```

 * logfile 日志路径
 * host redis 地址
 * port redis 端口

### Server A01

```bash
# cat /etc/cron.d/cronshell
MAILTO=""
SHELL=/bin/cronshell
PATH=/sbin:/bin:/usr/sbin:/usr/bin

* * * * * root d=$(date); echo $d run ok >> /tmp/t.log

# cat /tmp/t.log
Sat May 30 22:47:01 CST 2020 run ok
Sat May 30 22:48:01 CST 2020 run ok
Sat May 30 22:51:01 CST 2020 run ok
Sat May 30 22:54:01 CST 2020 run ok
```

### Server B01

```bash
# cat /etc/cron.d/cronshell
MAILTO=""
SHELL=/bin/cronshell
PATH=/sbin:/bin:/usr/sbin:/usr/bin

* * * * * root d=$(date); echo $d run ok >> /tmp/t.log

# cat /tmp/t.log
Sat May 30 22:49:01 CST 2020 run ok
Sat May 30 22:50:01 CST 2020 run ok
Sat May 30 22:52:01 CST 2020 run ok
Sat May 30 22:53:01 CST 2020 run ok
Sat May 30 22:55:01 CST 2020 run ok
Sat May 30 22:56:01 CST 2020 run ok
```
