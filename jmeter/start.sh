#！ /bin/bash
rm -rf queue-report.jtl
rm -rf queue-report-final.jtl
rm -rf report/*
sh apache-jmeter-5.2/bin/jmeter.sh -n -t queue-service.jmx -e -o report/ -l queue-report-final.jtl -Jthread-num=2000 -Jloop-times=100 -Jreport-file=queue-report.jtl 

