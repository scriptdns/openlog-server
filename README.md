# openlog
Open source log or data aggregation server

## Design

The OpenLog server can take structured logs or data in various formats using input plugins, process it in different ways using processing plugins, then output it to various locations using different output plugins. 

Each log message is pushed to a "stream" with a dot-separated heirarchy (eg, `web.api.big-api`).