# xb

A simple BIRD CLI tool to show BGP peers in a pretty table.

![xb](example.png)

## instructions

Download one of the binaries from the releases page.

```bash
wget https://github.com/x6c-co/xb/releases/download/v0.1.0-alpha/xb_Linux_x86_64.tar.gz
tar -xf xb_Linux_x86_64.tar.gz
sudo install ./xb /usr/local/bin
xb
```

By default, `xb`, looks for the `bird.ctl` file at `/var/run/bird/bird.ctl`. You
can tell `xb` to use a different path by either passing the `-socket` argument
or using the `BIRD_SOCKET` environment variable.
