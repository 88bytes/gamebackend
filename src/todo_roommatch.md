public class HelloWorld
{
    public string Message { get; set; }
}

// Serialize object to JSON
var toObject = new HelloWorld{ Message = "Hello world!" };
var toJson = LitJson.JsonMapper.ToJson(toObject);
Console.WriteLine("To JSON: {0}", toJson);

// Serialize JSON to object
var fromJson = "{\"Message\":\"Hello world!\"}";
var fromObject = LitJson.JsonMapper.ToObject<HelloWorld>(fromJson);	
Console.WriteLine("From json: {0}", fromObject.Message);

--------- ↑↑↑ LitJson用例 ↑↑↑ --------- 

服务器代码，new 对象的时候，我没有处理好，写的不统一。
需要重构整理代码，回头把代码重新写一遍。

1.OnJoinPlayer 刷新UI
2.开始比赛
3.完成比赛
4.各种情况-回收资源

-------------------- ↓↓↓ developing ↓↓↓ -------------------- 
    - 等待服务器返回的消息
    - 受到OnJoinPlayer的消息之后，刷新UI显示上桌的玩家信息
        - 要知道每个玩家坐在哪里

-------------------- ↓↓↓ developing ↓↓↓ -------------------- 
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

-------------------- ↓↓↓ developing ↓↓↓ -------------------- 
0.如果中途有人退出了，也一定要把数据处理好！
1.房间对战完毕一定要回收好！
2.如果玩家已经在一个对战中了，不能再创建新的对战，不能加入新的对战
3.第2点的内容要把自由匹配和私人房都确定好
