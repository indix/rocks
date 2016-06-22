[![Build Status](https://snap-ci.com/ind9/rocks/branch/master/build_image)](https://snap-ci.com/ind9/rocks/branch/master)
# Rocks

Rocks is a RocksDB Ops CLI. It is portable and helps you do perform common administrative operations on one or many rocksdb instances.

## Usage
```
$ rocks --help
Perform common ops related tasks on one or many RocksDB instances.

Find more details at https://github.com/ind9/rocks

Usage:
  rocks [command]

Available Commands:
  restore     Restore backed up rocksdb files
  backup      Backup rocksdb files
  trigger     Triggers a backup on the remote system

Use "rocks [command] --help" for more information about a command.
```

### Restore
```
$ rocks restore --help
Restore backed up rocksdb files

Usage:
  rocks restore [flags]

Flags:
      --dest string      Restore to
      --keep-log-files   If true, restore won't overwrite the existing log files in wal_dir
      --recursive        Trying restoring in recursive fashion from src to dest
      --src string       Restore from
      --wal string       Restore WAL to (generally same as --dest)
```

### Backup
```
$ rocks backup --help
Backs up rocksdb files

Usage:
  rocks backup [flags]

Flags:
      --dest string      Backup to
      --recursive        Trying to backup in recursive fashion from src to dest
      --src string       Backup from
```

### Remote Backup Trigger
```
$ rocks trigger --help
Triggers a backup on the remote system

Usage:
  rocks trigger [command]

Available Commands:
  http     Triggers backup via http

Flags for http command:
  -H, --header value             HTTP Headers that needs to be passed with the request in key=value format (default [])
      --http2                    Make a HTTP2 request instead of the default HTTP1.1
  -X, --method string            HTTP Method to invoke on the url (default "GET")
  -D, --payload string           File contents to be passed as payload in the request. If you're pipeing use "stdin".
  -k, --skip-ssl-verificatioin   Ignore SSL errors - while connecting to self-signed HTTPS servers
```

## License
http://www.apache.org/licenses/LICENSE-2.0
