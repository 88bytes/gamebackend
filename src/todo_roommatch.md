-------------------- ↓↓↓ developing ↓↓↓ -------------------- 
基本的SDKDemo 登陆功能已经测试完毕，现在要做的：

1.两个已知的Bug
    - 验证私人房的比赛还有没有问题
    - 头像加载有问题
        - 没有正确的缓存头像
        - 没有好看的默认头像框，黑色图片太丑了

2.存储和支付
    - 用户信息数据存储、金币数据存储
        - 服务器存了用户数据、存了金币数据，把数据库搭建好
        - 点击商城的按钮直接购买成功
        - 接入SDK支付，根据支付结果来增加金币

-------- 后期功能 -------- 
4.新版本私人房的界面

5.后期版本：   
    - 代理以及管理后台的功能，工具可以写在客户端里面

-------------------- ↓↓↓ developing ↓↓↓ -------------------- 

4.完整的多轮私人房
    - 这个我很想今天做完，这个一定做完了才能算完整了，但是我不打算管这么多

5.断线重连先不管，但是如果一个玩家的比赛房间没有解散(结算的道具没有完成结算)，不让他加入新的比赛

-------------------- ↓↓↓ done ↓↓↓ -------------------- 

--- done ---
- 听牌的时候，会自动跳过
    - 有Replay了，接下来把核心玩法的Bug修好
    ↑↑↑ 定位到了，不是听牌的问题，是因为网络的卡顿延迟，在PutCard阶段处理了AbandonChiPengGang的命令，导致摸牌玩家成了下一个玩家。

--- done ---
- 优化快播的逻辑，尝试一下，干脆不要快播，直接数据来多少就播放多少

--- done ---
- 加入不存在的房间房间，会把服务器干死
    --- done ---
    - 先把创建、加入房间的回调干掉 -> 这个只是因为我用了 logger.Fatal("messages.")
    --- done ---
    - 修复私人房的Bug，再打开

0.优化：
    - 排行榜的数字
    - 音效

2.加入新的UI，整合他们的功能：
    - 加入排行榜的UI，做个假的
    - 加入新的UI，把动画整合进来（快速开始、创建房间、加入房间）

1.整理重构代码，保证游戏核心代码可维护性
    - 整理客户端的代码，把蛋疼的 AfterSystem啥啥啥 的改掉
    - 其他看看还有没有想优化的代码 ...

3.Debug，保证私人房的内存回收稳定性
    - review、清理干净服务器的代码，保证服务器的稳定

2.把SDK登陆整合好
    - 自己的微信登陆已经好了
    - 能够设置正确的用户名字
    - 能够正确显示用户的头像
    - 能够在对战的时候，正确显示所有用户的头像

-------------------- ↓↓↓ developing ↓↓↓ -------------------- 
0.回放的文件需要优化，要不然写磁盘的时候太卡了
1.监控工具，可以看到服务器上所有的资源情况



---------------------- ILRuntime ---------------------- 
5.ILRuntime 开发新游戏
    - 开发麻将版本，把大厅做完整，再来搞
    - 我觉得我的UI框架不好用，太麻烦了，使用体验不好，需要重新搞一把
        - 每个UI都有 开、关 动画
    - 使用ILRuntime开发新版本

