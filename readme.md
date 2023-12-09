## 設計思路

> first version
### 客戶端主動推送
- ready(所有玩家都按下ready遊戲就開始)
- hit(玩家要牌)
- stand(玩家停牌)  

### 客戶端被動接收
- 更新玩家資訊和牌組
- 更新某玩家要牌、停牌資訊
- 更新遊戲狀態
  - 準備狀態
  - 開始遊戲
  - 遊戲結束
- 更新最終遊戲結果

## 代辦

> 20231209
### 後端 
- [ ] 修正玩家獲勝判斷
- [ ] 玩家註冊和登入(簡單版 不用email)
- [ ] 玩家斷線或者中途離場判斷
  - [ ] 剩一人 遊戲跳回等待模式
  - [ ] 一人都沒 遊戲跳回等待模式
  - [X] 剩一人以上 遊戲邏輯直接不判斷離開玩家
- [ ] 多房間
  - [ ] 增加多房間功能
  - [ ] 房間滿四人就不可再進入
  - [ ] 房間不滿四人 進入後為等待玩家 可觀戰但不可遊戲 須等下一輪後才可開始
  - [ ] 創建新房間
  - [ ] 刪除房間
- [ ] 倒數機制(每個階段都有限定時間)

## 前端
- [ ] 註冊登入頁面
- [ ] Lobby頁面 (可選擇房間)
