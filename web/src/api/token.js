import request from '@/utils/request'

export const index = (params) => {
  return request({
    method: 'get',
    url: '/api/admin/api_tokens',
    params: params
  })
}
