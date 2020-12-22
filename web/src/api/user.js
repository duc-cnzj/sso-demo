import request from '@/utils/request'

export const index = (params) => {
  return request({
    method: 'get',
    url: '/api/admin/users',
    params: params
  })
}

export const store = (data) => {
  return request({
    method: 'post',
    url: `/api/admin/users`,
    data
  })
}

export const show = (id) => {
  return request({
    method: 'get',
    url: `/api/admin/users/${id}`
  })
}

export const update = (id, data) => {
  return request({
    method: 'put',
    url: `/api/admin/users/${id}`,
    data
  })
}

export const destroy = (id) => {
  return request({
    method: 'delete',
    url: `/api/admin/users/${id}`
  })
}

export const syncRoles = (id, data) => {
  return request({
    method: 'post',
    url: `/api/admin/users/${id}/sync_roles`,
    data
  })
}

export const changePwd = (id, password) => {
  return request({
    method: 'post',
    url: `/api/admin/users/${id}/change_password`,
    data: {
      password
    }
  })
}

export const forceLogout = (id) => {
  return request({
    method: 'post',
    url: `/api/admin/users/${id}/force_logout`
  })
}

export function login(data) {
  return request({
    url: '/vue-admin-template/user/login',
    method: 'post',
    data
  })
}

export function getInfo(token) {
  return request({
    url: '/vue-admin-template/user/info',
    method: 'get',
    params: { token }
  })
}

export function logout() {
  return request({
    url: '/vue-admin-template/user/logout',
    method: 'post'
  })
}
