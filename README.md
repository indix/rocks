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
      --threads int      Number of CPU threads to work with (default : 2 * #(CPU Cores))
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
      --threads int      Number of CPU threads to work with (default : 2 * #(CPU Cores))
```

### Remote Backup Trigger
```
$ rocks trigger http --help
Triggers a backup via HTTP

Usage:
  rocks trigger http [flags]

Flags:
  -H, --header value             HTTP Headers that needs to be passed with the request in key=value format (default [])
      --http2                    Make a HTTP2 request instead of the default HTTP1.1
  -X, --method string            HTTP Method to invoke on the url (default "GET")
  -D, --payload string           File contents to be passed as payload in the request. If you're pipeing use "stdin".
  -k, --skip-ssl-verificatioin   Ignore SSL errors - while connecting to self-signed HTTPS servers
```

### Compact
```
$ rocks compact --help
Does a compaction on rocksdb stores

Usage:
  rocks compact [flags]

Flags:
      --recursive     Trying to compact in recursive fashion for src
      --src string    Compact for
      --threads int   Number of CPU threads to work with (default : 2 * #(CPU Cores))
```

### Statistics
```
$ rocks statistics --help
Displays current statistics for a rocksdb store

Usage:
  rocks statistics [flags]

Flags:
      --recursive     Trying to generate statistics in recursive fashion for src
      --src string    Statistics for
      --threads int   Number of threads to generate statistics (default : 2 * #(CPU Cores)
```

### Consistency
```
Checks for the consistency between rocks store and it's corresponding restore

Usage:
  rocks consistency [flags]

Flags:
      --src string    Rocks store location
      --dest string   Restore location for Rocks store
      --recursive     Trying to check consistency between rocks store and and it's restore
```

## License
http://www.apache.org/licenses/LICENSE-2.0
