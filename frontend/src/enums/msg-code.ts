export enum EOperationCode {
  // 玩家加入
  ClientJoin = 1,

  // 廣播某玩家進入遊戲
  BroadcastJoin = 2,

  // 廣播某玩家離開遊戲
  BroadcastLeave = 3,

  // 玩家準備
  ClientReady = 4,

  // 廣播某玩家按下準備
  BroadcastReady = 5,

  // 廣播遊戲開始
  BroadcastGameStart = 6,

  // 玩家要牌
  ClientHit = 7,

  // 廣播某玩家要牌
  BroadcastHit = 8,

  // 廣播有人爆牌
  BroadcastBust = 9,

  // 玩家停止要牌
  ClientStand = 10,

  // 廣播某玩家停止要牌
  BroadcastStand = 11,

  // 廣播遊戲結束
  BroadcastGameOver = 12,

  // 廣播遊戲重新開始
  BroadcastReStart = 13,

  // 更新所有玩家資訊
  UpdatePlayersDetail = 14
}
