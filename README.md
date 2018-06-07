# qiita-trends-rss

[Qiita](https://qiita.com/) のトレンドのRSSフィードのサーバー

[Qiita](https://qiita.com/) trends RSS feed server


## install

```
go get github.com/Mushus/qiita-trends-rss
```

or

```
# for linux
wget https://github.com/Mushus/qiita-trends-rss/releases/download/v0.1.0/qiita-trends-rss-amd64-linux -O qiita-trends-rss
chmod 755 qiita-trends-rss
sudo mv qiita-trends-rss /usr/local/bin/qiita-trends-rss
```
サーバー用途なので win mac はないです

## run

run command as:
```
qiita-trends-rss -p 1234
```

operation check:
```
curl http://localhost:1234/
```
