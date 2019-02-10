import { stringify } from 'qs';
import request, { v1API } from '../utils/request';

const router = 'roles';

export async function query(params) {
  return request(`${v1API}/${router}?${stringify(params)}`);
}

export async function get(params) {
  return request(`${v1API}/${router}/${params.record_id}`);
}

export async function create(params) {
  return request(`${v1API}/${router}`, {
    method: 'POST',
    body: params,
  });
}

export async function update(params) {
  return request(`${v1API}/${router}/${params.record_id}`, {
    method: 'PUT',
    body: params,
  });
}

export async function del(params) {
  return request(`${v1API}/${router}/${params.record_id}`, {
    method: 'DELETE',
  });
}

export async function delMany(params) {
  return request(`${v1API}/${router}?${stringify(params)}`, {
    method: 'DELETE',
  });
}
