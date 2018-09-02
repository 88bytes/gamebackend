1.创建房间的时候，先发送消息，然后直接进入场景

2.发送 CreateRoom 的消息给服务器，进入 WaitRoomMatch 的状态
    - 房主的座位从4个位置随机一个分配一下吧
    - WaitRoomMatch 状态中，等待服务器的回复，等待回复的时候，把几个UI都关掉

3.服务器回复消息 OnJoinPlayer 
    - 收到了这个消息的时候，就会听服务器的吩咐，把位置坐好，把UI打开
    - 把UI打开的同时可以显示坐在位置上面的玩家信息

4.发送 JoinRoom 消息给服务器，进入 WaitRoomMatch 的状态
    - 收到了 OnJoinPlayer 的消息之后，刷新UI

5.主机可以添加AI来让AI加入比赛
    - AddAIPlayer 消息给服务器发通知

6.服务器同样通过 OnJoinPlayer 消息告诉每个客户端有AI玩家进入房间

7.服务器通过 StartRoomBattle 的消息启动私人房的比赛
    - 这个消息里面会包含：
        - 房间号
        - 座位信息
        - 4个玩家的信息
        - 比赛启动的随机数种子

8.接下里客户端的流程和快速匹配就是一模一样的