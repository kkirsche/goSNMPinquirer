# Inquirer
Easy to use / configure SNMP poller written in Golang. Configurable using JSON like so:

```json
{
  "ip": "127.0.0.1",
  "community": "Example-SNMP-Community-String",
  "cron": {
    "save_via": "file",
    "save_filepath": "/Users/exampleUser/Desktop/",
    "save_filename": "results",
    "minute": {
      "get": [".1.3.6.1.2.1.1.5.0"],
      "get_readable": ["SNMPv2-MIB::sysName"],
      "getbulk": [],
      "getbulk_readable": [],
      "bulkwalk": [
        ".1.3.6.1.2.1.2.2.1.10", ".1.3.6.1.2.1.2.2.1.16",
        ".1.3.6.1.2.1.2.2.1.1", ".1.3.6.1.2.1.2.2.1.11",
        ".1.3.6.1.2.1.2.2.1.17", ".1.3.6.1.2.1.2.2.1.19"
      ],
      "bulkwalk_readable": [
        "IF-MIB::ifInOctets", "IF-MIB::ifOutOctets",
        "IF-MIB::ifIndex", "IF-MIB::ifInUcastPkts",
        "IF-MIB::ifOutUcastPkts", "IF-MIB::ifOutDiscards"
      ]
    },
    "hour": {
      "get": [".1.3.6.1.2.1.1.5.0"],
      "get_readable": ["SNMPv2-MIB::sysName"],
      "getbulk": [],
      "getbulk_readable": [],
      "bulkwalk": [
        ".1.3.6.1.2.1.2.2.1.1", ".1.3.6.1.2.1.31.1.1.1.1"
      ],
      "bulkwalk_readable": [
        "IF-MIB::ifIndex", "IF-MIB::ifName"
      ]
    },
    "day": {
      "get": [],
      "get_readable": [],
      "getbulk": [],
      "getbulk_readable": [],
      "bulkwalk": [],
      "bulkwalk_readable": []
    }
  }
}
```

This will query `127.0.0.1` using the community string `Example-SNMP-Community-String` for the OID's provided in `get`, `getbulk`, and `bulkwalk`. This file by default should be stored in `$HOME/.inquirer.json` but can be adjusted with the `--settings` / `-s` command line arguments.

## Values

### save_via

* `stdout`
* `file` — CSV format only
* `syslog`

### save_file

* Fully qualified path, including file name and extension, to store the results in.

## Like to Add
I'd love to add:

* MIB lookups for OID values (OID <-> MIB depending on situation and configuration)
