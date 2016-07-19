import json
import sys
from datetime import datetime, timedelta

with open(sys.argv[1]) as data_file:
    data = json.load(data_file)

timestamp = data["version"]["time"]
truncated = timestamp[:26] + timestamp[-1]

utc_dt = datetime.strptime(truncated, '%Y-%m-%dT%H:%M:%S.%fZ')
seconds = (utc_dt - datetime(1970, 1, 1)).total_seconds()
print int(round(seconds))
