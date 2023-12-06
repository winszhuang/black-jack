export enum EMsgCode {
  OneJoin = 1,
  // 某玩家進入遊戲
  BroadcastJoin = 2,
  // 某玩家離開遊戲
  BroadcastLeave = 3,

  OneReady = 4,
  // 某玩家按下準備
  BroadcastReady = 5,
  // 遊戲開始
  BroadcastGameStart = 6,

  OneHit = 7,
  // 玩家要牌
  BroadcastHit = 8,

  OneStand = 9,

  // 玩家停止要牌
  BroadcastStand = 10,
  // 遊戲結束
  BroadcastGameOver = 11,

  // 更新所有玩家資訊
  UpdatePlayersDetail = 12
}
