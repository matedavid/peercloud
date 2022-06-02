if [ ! -d ".peercloud" ] 
then
  mkdir .peercloud
fi

if [ ! -z "$1" ] 
then
  if [ ! -d ".peercloud/$1" ] 
  then
    mkdir .peercloud/$1
    mkdir .peercloud/$1/.shards
    mkdir .peercloud/$1/.storage
    mkdir .peercloud/$1/.tmp
    touch .peercloud/$1/hosts
  fi
else 
  mkdir .peercloud/.shards
  mkdir .peercloud/.storage
  mkdir .peercloud/.tmp
  touch .peercloud/hosts
fi
