#!/usr/bin/python
#
# Rubrik Prometheus Client
#
# Requirements: Python 2.7
#               rubrik_cdm module
#               prometheus_client module
#               Rubrik CDM 3.0+
#               Environment variables for RUBRIK_IP (IP of Rubrik node), RUBRIK_USER (Rubrik username), RUBRIK_PASS (Rubrik password)
#
# NOTE: 2019-09-11: the Python version of this client has been deprecated in favour of the GoLang client. The GoLang client can be found in this same repository in the `src/golang` folder. The Python version of this agent will no longer be maintained.
#

from prometheus_client import start_http_server
from prometheus_client import Gauge
import rubrik_cdm
import time
import json
import os

# disable warnings
import urllib3
urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)

# Define our metrics
RUBRIK_TOTAL_STORAGE = Gauge('rubrik_total_storage_bytes', 'Total storage in Rubrik cluster')
RUBRIK_USED_STORAGE = Gauge('rubrik_used_storage_bytes', 'Used storage in Rubrik cluster')
RUBRIK_AVAILABLE_STORAGE = Gauge('rubrik_available_storage_bytes', 'Available storage in Rubrik cluster')
RUBRIK_SNAPSHOT_STORAGE = Gauge('rubrik_snapshot_storage_bytes', 'Snapshot storage in Rubrik cluster')
RUBRIK_LIVEMOUNT_STORAGE = Gauge('rubrik_livemount_storage_bytes', 'Live Mount storage in Rubrik cluster')
RUBRIK_MISC_STORAGE = Gauge('rubrik_misc_storage_bytes', 'Miscellaneous storage in Rubrik cluster')

def get_rubrik_stats():
    rubrik = rubrik_cdm.Connect(node_ip=os.environ['RUBRIK_IP'],username=os.environ['RUBRIK_USER'],password=os.environ['RUBRIK_PASS'])
    stats = rubrik.get('internal','/stats/system_storage')
    RUBRIK_TOTAL_STORAGE.set(stats['total'])
    RUBRIK_USED_STORAGE.set(stats['used'])
    RUBRIK_AVAILABLE_STORAGE.set(stats['available'])
    RUBRIK_SNAPSHOT_STORAGE.set(stats['snapshot'])
    RUBRIK_LIVEMOUNT_STORAGE.set(stats['liveMount'])
    RUBRIK_MISC_STORAGE.set(stats['miscellaneous'])

if __name__ == '__main__':
    # Start up the server to expose the metrics.
    start_http_server(9477)
    while True:
        get_rubrik_stats()
        time.sleep(600)
