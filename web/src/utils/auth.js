import Cookies from 'js-cookie'

export function getToken() {
  return Cookies.get('sso_token')
}

export function setToken(token, lifetime) {
  return Cookies.set('sso_token', 'Bearer ' + token, { expires: 365 })
}

export function removeToken() {
  return Cookies.remove('sso_token')
}
