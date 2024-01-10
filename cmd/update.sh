#!/bin/sh

cd /data/gocode/gin-gorm-base
git reset --hard
git pull
go mod tidy

echo '2. Building ...'

#删除旧的编译生成文件
if [ -f gin-gorm-base ]; then
        rm gin-gorm-base
fi

go build -o gin-gorm-base main.go

# 检查是否编译成功
if [ ! -f gin-gorm-base ]; then
        echo "ERR: build error, return now"
else
        # 编译成功了，这里开始备份当前正在运行的程序文件，一定要备份/data/goservice/gin-gorm-base/gin-gorm-base，而不是 /data/code/gin-gorm-base/gin-gorm-base。
        # 因为/data/goservice/gin-gorm-base/gin-gorm-base，备份它肯定没错。
        echo "3. Backup old version"
        if [ -f /data/goservice/gin-gorm-base/gin-gorm-base ]; then
                #mv /data/goservice/gin-gorm-base/gin-gorm-base /data/goservice/gin-gorm-base/gin-gorm-base.`date "+%Y-%m-%d_%H:%M:%S"`
                if [ -f /data/goservice/gin-gorm-base/gin-gorm-base.backup ]; then
                    rm /data/goservice/gin-gorm-base/gin-gorm-base.backup
                fi
                mv /data/goservice/gin-gorm-base/gin-gorm-base /data/goservice/gin-gorm-base/gin-gorm-base.backup
        fi
        mv gin-gorm-base /data/goservice/gin-gorm-base/gin-gorm-base

        echo "4. Restart service"
        systemctl restart gin-gorm-base
        echo "Service enable :"
        systemctl enable gin-gorm-base
        echo "Service status is :"
        systemctl status gin-gorm-base
        echo ""
        echo "Restart OK, visit http://127.0.0.1:8000/ping for test"
fi
