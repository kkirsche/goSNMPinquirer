# inquirer
Easy to use / configure SNMP poller written in Golang. Configurable using JSON like so:

```
{
  "ip": "127.0.0.1",
  "community": "Example-SNMP-Community-String",
  "cron": {
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
    "day": {
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
    }
  }
} 
```

This will query `127.0.0.1` using the community string `Example-SNMP-Community-String` for the OID's provided in `get`, `getbulk`, and `bulkwalk`. This file by default should be stored in `$HOME/.inquirer.json` but can be adjusted with the `--settings` / `-s` command line arguments.
