 go build -o WebManageSvr main.go
echo "build ok"
sshpass -p blct4@hq!2020 scp WebManageSvr ct4user@47.113.92.166:~/ct4/zybTest/WebManageSvr/WebManageSvr
sshpass -p blct4@hq!2020 scp WebManageSvr.toml ct4user@47.113.92.166:~/ct4/zybTest/WebManageSvr/WebManageSvr.toml
#sshpass -p blhangqing scp enmapcn.json ct4user@192.168.1.88:~/ct4/zybtest/quote_micro_stockInfos/test/quote_micro_stockInfoApi/enmapcn.json
#sshpass -p blhangqing scp quote_micro_stockInfoApi.toml ct4user@192.168.1.88:~/ct4/zybtest/quote_micro_stockInfos/test/quote_micro_stockInfoApi.toml
#echo "scp . ok"
#sshpass -p blct4hq2020 scp quote_micro_stockInfoApi.toml ct4user@47.107.156.18:~/ct4/quote_micro_stockInfoApi.toml
#sshpass -p blct4hq2020 scp quote_micro_stockInfoApi.toml ct4user@192.168.1.88:~/ct4/quote_micro_stockInfoApi.toml

#echo "scp toml ok"
#sshpass -p blct4hq2020 scp enmapcn.json ct4user@47.107.156.18:~/ct4/enmapcn.json

echo "scp .json ok"
