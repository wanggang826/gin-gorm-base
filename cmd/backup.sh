#!/bin/sh

cd /data/goservice/gin-gorm-base

mv gin-gorm-base gin-gorm-base.0

#mv $1 main
mv gin-gorm-base.backup gin-gorm-base

systemctl daemon-reload

systemctl restart gin-gorm-base.service

#echo "gin-gorm-base backup to "$1
echo "gin-gorm-base backup to gin-gorm-base.backup"