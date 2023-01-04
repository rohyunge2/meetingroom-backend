#!/bin/bash
if [ $# -lt 1 ]; 
then
	VERSION="latest"
else
	VERSION="$1"
fi

git pull
docker build -t qcr.k8s.bns.co.kr/bnspace/meeting-backend:$VERSION .
docker login --username admin --password BNSoft2020@ qcr.k8s.bns.co.kr
docker push qcr.k8s.bns.co.kr/bnspace/meeting-backend:$VERSION
kubectl set image deployment/bnspace-meeting-backend golang=qcr.k8s.bns.co.kr/bnspace/meeting-backend:$VERSION --record
#kubectl delete pod -l "app=bnspace-meeting-backend"
#docker stop test
#docker rm test
#docker run -d -p 8080:8080 --name test qcr.k8s.bns.co.kr/bnspace/meeting-backend:$VERSION
