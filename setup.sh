if [ ! -d ".peercloud" ] 
then
  mkdir .peercloud
fi

if [ ! -z "$0" ] 
then
  mkdir .peercloud/$1
  mkdir .peercloud/$1/.shards
  mkdir .peercloud/$1/.storage
  mkdir .peercloud/$1/.tmp
  touch .peercloud/$1/hosts
else 
  mkdir .peercloud/.shards
  mkdir .peercloud/.storage
  mkdir .peercloud/.tmp
  touch .peercloud/hosts
fi
