import request from '@/utils/request'

export const index = (params) => {
  console.log(params)
  return request({
    method: 'get',
    url: '/api/roles',
    params
  })
}

export const store = (data) => {
  return request({
    method: 'post',
    url: `/api/roles`,
    data
  })
}

export const show = (id) => {
  return request({
    method: 'get',
    url: `/api/roles/${id}`
  })
}

export const update = (id, data) => {
  return request({
    method: 'put',
    url: `/api/roles/${id}`,
    data
  })
}

export const destroy = (id) => {
  return request({
    method: 'delete',
    url: `/api/roles/${id}`
  })
}

export const allRoles = (id) => {
  return request({
    method: 'get',
    url: `/api/all_roles`
  })
}
