Game: 负责游戏逻辑
    ActionChannel：服务端发送过来的用户操作
    ChangeChannel：返回给服务端的游戏状态，用于更新游戏状态并广播

Game 状态机：
    1. 等待玩家加入 Wait
    2. 当人数足够时，开始游戏 Start
        DealCards -> Play
        监听 ActionChannel，接收用户操作，并更新游戏状态
    3. 游戏结束 End

Server：负责将与 Client 通信
    监听 game.ChangeChannel，将游戏状态广播给所有 Client
    Stream： 负责与 Client 通信，并将用户操作发送给 Game
