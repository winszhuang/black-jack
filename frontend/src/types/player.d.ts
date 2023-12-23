type PlayerDetail = {
  id: string
  name: string
  deck: Card[]
  state: import('@/enums/player-status').EPlayerState
}
