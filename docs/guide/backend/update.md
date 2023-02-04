# 后端升级

## 更新变动

突破性变动，例如更新需要修改原有配置等，请注意查看版本更新日志中的提示：[CHANGELOG.md](https://github.com/ArtalkJS/Artalk/blob/master/CHANGELOG.md)

## 命令行一键升级

执行 `./artalk upgrade`

此操作会从自动从 GitHub Release 下载并升级程序，执行前需关闭 Artalk。

:::tip

执行 `./artalk upgrade -f` 携带参数 `-f` 来进行同版本号的补充更新。

:::

## Docker 升级

可参考：[“Docker · 升级”](./docker.md#升级)

## 普通方式

前往 [GitHub Release](https://github.com/ArtalkJS/Artalk/releases) 手动下载最新构建

替换掉旧版本文件即可。
