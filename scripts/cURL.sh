#!/bin/bash

set -x

# In curl >= 7.18.0, use
# curl http://www.baidu.com --socks5-hostname localhost:1080

# In curl >= 7.21.7, use
curl http://www.baidu.com -x socks5h://localhost:1080

# for username/password, use
# curl http://www.baidu.com -x socks5h://username:password@localhost:1080

#curl --insecure -v -X GET  https://localhost:8090/v1/connection/list