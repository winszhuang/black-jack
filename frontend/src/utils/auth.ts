export function isLogin() {
  return !!localStorage.getItem('access_token')
}
