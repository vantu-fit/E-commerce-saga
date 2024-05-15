#!/bin/bash

# becase we can't connect to redis cluster with localhost when use docker
# I have tried to use docker-compose but can't connect to redis cluster with localhost (network_mode=host,HAProxy,..)
# So that I have to install redis cluster on WSL (unbuntu) and connect to it with localhost


#install bloom filter
mkdir ~/Redis
cd ~/Redis
apt-get update -y && apt-get upgrade -y
apt-get install -y wget make pkg-config build-essential
wget https://download.redis.io/redis-stable.tar.gz
tar -xzvf redis-stable.tar.gz
cd redis-stable
make distclean
make
make install


apt-get install -y git
cd ~/Redis
git clone --recursive https://github.com/RedisBloom/RedisBloom.git
cd RedisBloom
./sbin/setup
bash -l
make

# change user
make run -n /home/vantu/RedisBloom/bin/linux-x64-release/redis-server
#check location file redisbloom.so
find /home/vantu/RedisBloom -name redisbloom.so

#then run this script
rm -rf ~/redis-cluster
# Set base directory for Redis cluster
BASE_DIR=~/redis-cluster
PORTS=(7000 7001 7002 7003 7004 7005)
IP=0.0.0.0

# Đường dẫn đầy đủ đến tệp redisbloom.so
REDIS_BLOOM_SO=/home/vantu/RedisBloom/bin/linux-x64-release/redisbloom.so

# Create directories for each Redis instance
for PORT in "${PORTS[@]}"; do
  mkdir -p ${BASE_DIR}/${PORT}
done

# Create Redis configuration files for each instance
for PORT in "${PORTS[@]}"; do
  CONFIG_FILE=${BASE_DIR}/${PORT}/redis.conf
  cat <<EOF > ${CONFIG_FILE}
bind ${IP}
port ${PORT}
cluster-enabled yes
cluster-config-file nodes.conf
cluster-node-timeout 5000
protected-mode no
appendonly yes
daemonize yes
dir ${BASE_DIR}/${PORT}
loadmodule ${REDIS_BLOOM_SO}
EOF
done

# Start each Redis instance
PORTS=(7000 7001 7002 7003 7004 7005)
for PORT in "${PORTS[@]}"; do
  redis-server ${BASE_DIR}/${PORT}/redis.conf
done

# Wait for Redis instances to start
sleep 5

# Create Redis cluster
yes "yes" | redis-cli --cluster create 127.0.0.1:7000 127.0.0.1:7001 127.0.0.1:7002 127.0.0.1:7003 127.0.0.1:7004 127.0.0.1:7005 --cluster-replicas 1

# Check Redis cluster status
for PORT in "${PORTS[@]}"; do
  redis-cli -p ${PORT} cluster info
done

# shut down all instances
for PORT in {7000..7005}; do
  redis-cli -p ${PORT} shutdown
done

rm -rf ~/redis-cluster

