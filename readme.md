# 目录结构
参考：https://github.com/golang-standards/project-layout/blob/master/README_zh.md

/cmd 每个应用程序的目录名应该与你想要的可执行文件的名称相匹配(例如，/cmd/myapp)

/internal 私有应用程序和库代码

/pkg 外部应用程序可以使用的库代码(例如 /pkg/mypubliclib)

/vendor 应用程序依赖项(手动管理或使用你喜欢的依赖项管理工具，如新的内置 Go Modules 功能)。go mod vendor 命令将为你创建 /vendor 目录

/api 

/web 特定于 Web 应用程序的组件:静态 Web 资源、服务器端模板和 SPAs

/configs 

/init

/scripts

/build

/deployments IaaS、PaaS、系统和容器编排部署配置和模板

/test 

/docs

/tools

/examples

/assets
