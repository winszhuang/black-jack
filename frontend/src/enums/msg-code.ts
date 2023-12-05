export enum EMsgCode {
  SomeOneJoin = 1,

  // 某玩家離開遊戲
  SomeOneLeave = 2,

  // 某玩家按下準備
  SomeOneReady = 3,

  // 遊戲開始
  GameStart = 4,

  // 玩家要牌
  SomeOneHit = 5,

  // 玩家停止要牌
  SomeOneStand = 6,

  // 遊戲結束
  GameOver = 7,

  // 更新所有玩家資訊
  UpdatePlayersDetail = 8
}
