# lxcstats

Show statistics about LXC containers in Proxmox VE.

## Usage

```console
# ./lxcstats 
Usage of ./lxcstats:
  ./lxcstats <disk>
  ./lxcstats [-cp] [-mp] [-ip] [-df]

  -cp
        list LXC with highest CPU pressures
  -df
        list LXC with highest root filesystem space usage
  -ip
        list LXC with highest I/O pressures
  -mp
        list LXC with highest memory pressures
# ./lxcstats -cp
ID    Avg10  Avg60  Avg300
1111  50.0   50.0   50.0
1112  1.4    1.4    1.5
1113  0.9    0.5    0.5
1114  0.3    0.2    0.2
1115  0.2    0.2    0.2
# ./lxcstats -ip
ID    Avg10  Avg60  Avg300
1111  0.0    0.0    0.0
1112  0.0    0.0    0.0
1113  0.0    0.0    0.0
1114  0.0    0.0    0.0
1115  0.0    0.0    0.0
# ./lxcstats -mp
ID    Avg10  Avg60  Avg300
1111  0.0    0.0    0.0
1112  0.0    0.0    0.0
1113  0.0    0.0    0.0
1114  0.0    0.0    0.0
1115  0.0    0.0    0.0
# ./lxcstats -df
Top stats of rootfs space
ID    RootFS used / total
1111  17.6 GB / 23.3 GB
1112  4.0 GB / 15.6 GB
# ./lxcstats /dev/sdd
2024-01-14 18:57:49
ID  Rios  Wios  Rbytes  Wbytes

2024-01-14 18:57:50
ID  Rios  Wios  Rbytes  Wbytes

2024-01-14 18:57:51
ID  Rios  Wios  Rbytes  Wbytes

2024-01-14 18:57:52
ID  Rios  Wios  Rbytes  Wbytes

2024-01-14 18:57:53
ID  Rios  Wios  Rbytes  Wbytes

2024-01-14 18:57:54
ID  Rios  Wios  Rbytes  Wbytes

2024-01-14 18:57:55
ID  Rios  Wios  Rbytes  Wbytes

2024-01-14 18:57:56
ID    Rios  Wios  Rbytes  Wbytes
1111  0     2     0 B     4.0 kB

2024-01-14 18:57:57
ID    Rios  Wios  Rbytes  Wbytes
1111  0     1     0 B     4.0 kB
^C
```
