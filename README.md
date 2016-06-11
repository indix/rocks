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

## License
http://www.apache.org/licenses/LICENSE-2.0
