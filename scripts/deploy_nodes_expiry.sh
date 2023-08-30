#!/usr/bin/env bash
basedir=$(cd `dirname $0`; pwd)
workspace=${basedir}/..

# global var
validatorAddr=()
validatorSecretLoc=()

function exit_previous() {
	# stop client
	ps -ef  | grep ${workspace}/bin/geth | awk '{print $2}' | xargs kill
	ps -ef  | grep ${workspace}/bin/bootnode | awk '{print $2}' | xargs kill
}

function start_bootnode() {
    tool=${workspace}/bin/bootnode
    
    file=${workspace}/scripts/asset/boot.key
    if [ ! -f "$file" ]; then
        $tool -genkey $file
    fi

    logFile=${workspace}/clusterNode/bootnode.log
    nohup $tool -nodekey ${file} -addr :30305 > $logFile 2>&1 &
    sleep 1

    # echo the first line
    enodeId=$(head -n 1 "$logFile")
    echo $enodeId
}

function generate_static_peers() {
    tool=${workspace}/bin/bootnode
    num=$1
    target=$2
    staticPeers=""
    for ((i=1;i<=$num;i++)); do
        if [ $i -eq $target ]
        then
           continue
        fi

        file=${workspace}/scripts/asset/nodekey${i}
        if [ ! -f "$file" ]; then
            $tool -genkey $file
        fi
        port=$((30331 + i))
        if [ ! -z "$staticPeers" ]
        then
            staticPeers+="\\,"
        fi
        staticPeers+="\"enode\\:\\/\\/$($tool -nodekey $file -writeaddress)\\@127\\.0\\.0\\.1\\:$port\""
    done

    echo $staticPeers
}

function generate_nodekey() {
    tool=${workspace}/bin/bootnode
    num=$1
    for ((i=1;i<=$num;i++)); do
        file=${workspace}/scripts/asset/nodekey${i}
        if [ ! -f "$file" ]; then
            $tool -genkey $file
        fi
    done
}

function prepare() {
    if ! [[ -f ${workspace}/bin/geth ]];then
         echo "bin/geth do not exist!"
         exit 1
    fi
    if ! [[ -f ${workspace}/bin/bootnode ]];then
         echo "bin/bootnode do not exist!"
         exit 1
    fi
    rm -rf ${workspace}/clusterNode
    cd ${workspace}/genesis
    npm install
}

function prepareGethEnv(){
    num=$1
    for((i=1;i<=$num;i++)); do
        rm -rf ${workspace}/clusterNode/node${i}
        mkdir -p ${workspace}/clusterNode/node${i}
        echo 'password' >> ${workspace}/clusterNode/password.txt
        ${workspace}/bin/geth --datadir ${workspace}/clusterNode/node${i} account new --password ${workspace}/clusterNode/password.txt > ${workspace}/clusterNode/validator${i}Info
        validatorAddr=("${validatorAddr[@]}" `cat ${workspace}/clusterNode/validator${i}Info|grep 'Public address of the key'|awk '{print $6}'` )
        validatorSecretLoc=("${validatorSecretLoc[@]}" `cat ${workspace}/clusterNode/validator${i}Info|grep  'Path of the secret key file'|awk '{print $7}'`)
    done
}

function generateGenesis(){
    rm ${workspace}/genesis/validators.conf
    num=$1
    
    for i in "${validatorAddr[@]}"
    do
       echo "${i}" >> ${workspace}/genesis/validators.conf
    done

    sed "s/{{INIT_HOLDER_ADDR}}/${validatorAddr[0]}/g" ${workspace}/genesis/init_holders.template > ${workspace}/genesis/init_holders.js
    node generate-validator.js
    node generate-genesis.js

    for((i=1;i<=$num;i++)); do
      ${workspace}/bin/geth --datadir ${workspace}/clusterNode/node${i} init ${workspace}/genesis/genesis.json
      staticPeers=$(generate_static_peers $num $i)
      sed "s/{{StaticNodes}}/${staticPeers}/g" ${workspace}/scripts/asset/config-cluster.toml > ${workspace}/clusterNode/node${i}/config.toml

      cp ${workspace}/scripts/asset/nodekey${i} ${workspace}/clusterNode/node${i}/geth/nodekey
    done
}

function startFullNodeWithExpiry() {
    num=$1
    nodeNum=$2
    bootnode=$3
    for((i=1;i<=$num;i++)); do
        validatorIndex=$(($nodeNum-1))
        nohup ${workspace}/bin/geth -unlock ${validatorAddr[$validatorIndex]} --http --http.port "$((8501+$nodeNum))" --ws.port "$((8545+$nodeNum))" \
         --bootnodes ${bootnode} --state-expiry "true" --state-expiry.remote "127.0.0.1:$((8501+${#validatorAddr[@]}))" \
         --syncmode "full" --config ${workspace}/clusterNode/node${nodeNum}/config.toml \
         --port "$((30305+$nodeNum))" --authrpc.port "$((8550+$nodeNum))" --password "${workspace}/clusterNode/password.txt" \
         --mine --miner.etherbase ${validatorAddr[$validatorIndex]} --rpc.allow-unprotected-txs --allow-insecure-unlock --light.serve 50 \
         --gcmode full --ws --datadir ${workspace}/clusterNode/node${nodeNum} \
         --metrics --pprof --pprof.port "$((6060+$nodeNum))" --http.corsdomain "https://remix.ethereum.org" > ${workspace}/clusterNode/node${nodeNum}/geth.log 2>&1 &

        echo "start validator $nodeNum as full node"
        nodeNum=$(($nodeNum+1))

        sleep 1
    done
}

function startFullNodeNoExpiry() {
    num=$1
    nodeNum=$2
    bootnode=$3
    for((i=1;i<=$num;i++)); do
        validatorIndex=$(($nodeNum-1))
        nohup ${workspace}/bin/geth -unlock ${validatorAddr[$validatorIndex]} --http --http.port "$((8501+$nodeNum))" --ws.port "$((8545+$nodeNum))" \
         --bootnodes ${bootnode} \
         --syncmode "full" --config ${workspace}/clusterNode/node${nodeNum}/config.toml \
         --port "$((30305+$nodeNum))" --authrpc.port "$((8550+$nodeNum))" --password "${workspace}/clusterNode/password.txt" \
         --mine --miner.etherbase ${validatorAddr[$validatorIndex]} --rpc.allow-unprotected-txs --allow-insecure-unlock --light.serve 50 \
         --gcmode full --ws --datadir ${workspace}/clusterNode/node${nodeNum} \
         --metrics --pprof --pprof.port "$((6060+$nodeNum))" --http.corsdomain "https://remix.ethereum.org" > ${workspace}/clusterNode/node${nodeNum}/geth.log 2>&1 &

        echo "start validator $nodeNum as full node"
        nodeNum=$(($nodeNum+1))
        sleep 1
    done
}

CMD=$1

case ${CMD} in
start)
    source ~/.bash_profile
    exit_previous
    fullNumWithExpiry=2
    if [ ! -z $3 ] && [ "$3" -gt "0" ]; then
      fullNum=$3
    fi
    fullNumNoExpiry=1
    if [ ! -z $4 ] && [ "$4" -gt "0" ]; then
      archiveNum=$4
    fi
    validatorNum=fullNumWithExpiry+fullNumNoExpiry
    echo "===== generate node key ===="
    generate_nodekey $validatorNum
    echo "===== preparing ===="
    bnbHolderAddr=$2
    prepare $validatorNum
    prepareGethEnv $validatorNum
    generateGenesis $validatorNum
    echo "===== starting bootnode ===="
    bootnode=$(start_bootnode)
    echo "===== starting client ===="
    startFullNodeNoExpiry fullNumNoExpiry $((fullNumWithExpiry+1)) $bootnode # By default, last node is remoteDB
    startFullNodeWithExpiry fullNumWithExpiry 1 $bootnode
    echo "Finish deploy"
    ;;
stop)
    echo "===== stopping client ===="
    exit_previous
    echo "===== client stopped ===="
    ;;
*)
    echo "Usage: deploy_cluster_nodes.sh start | stop"
    ;;
esac