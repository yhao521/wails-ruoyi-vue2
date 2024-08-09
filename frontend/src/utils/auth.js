// import Cookies from 'js-cookie'

// import VueCookies from 'vue-cookies'
// import store from './store'
const TokenKey = 'Admin-Token'

export function getToken() {
  // return Cookies.get(TokenKey)

  // return VueCookies.get(TokenKey)
  // return store.getters.token
  return sessionStorage.getItem(TokenKey);
}

export function setToken(token) {
  // return Cookies.set(TokenKey, token)
  // return VueCookies.set(TokenKey, token)
  return sessionStorage.setItem(TokenKey, token);
}

export function removeToken() {
  // return Cookies.remove(TokenKey)
  // return VueCookies.remove(TokenKey)
  return sessionStorage.removeItem(TokenKey);
}
