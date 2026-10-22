const API_URL = process.env.REACT_APP_API_URL || '';

async function apiFetch(path, options = {}) {
  const { body, token, ...customConfig } = options;
  const headers = { 'Content-Type': 'application/json' };

  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }

  const config = {
    method: body ? 'POST' : 'GET',
    ...customConfig,
    headers: {
      ...headers,
      ...customConfig.headers,
    },
  };

  if (body) {
    config.body = JSON.stringify(body);
  }

  const res = await fetch(`${API_URL}${path}`, config);
  
  if (res.status === 204) {
    return null;
  }

  const data = await res.json();
  if (!res.ok) {
    throw new Error(data.error?.message || 'เกิดข้อผิดพลาดจากเซิร์ฟเวอร์');
  }
  return data;
}

export async function getMovies(q = '', genreId = 0) {
  const params = new URLSearchParams();
  if (q) params.set('q', q);
  if (genreId) params.set('genre', genreId);
  const query = params.toString() ? `?${params.toString()}` : '';
  const data = await apiFetch(`/api/movies${query}`);
  return data.items || data;
}

export function getMovie(id) {
  return apiFetch(`/api/movies/${id}`);
}

export function loginUser(email, password) {
  return apiFetch('/api/auth/login', { body: { email, password } });
}

export function registerUser(email, password, displayName) {
  return apiFetch('/api/auth/register', { body: { email, password, displayName } });
}

export function getMe(token) {
  return apiFetch('/api/me', { token });
}

// ---------- แนวหนัง (สาธารณะ ใช้ทำ dropdown ในฟอร์มหลังบ้าน) ----------
export async function getGenres() {
  const data = await apiFetch('/api/genres');
  return data.items; // [{ id, nameEn, nameTh }]
}

// ---------- หลังบ้าน (ทุกตัวต้องแนบ token ของ admin ถ้าเป็น member Go ตอบ 403) ----------
export function getAdminStats(token) {
  return apiFetch('/api/admin/stats', { token });
}
export async function getAdminMembers(token) {
  const data = await apiFetch('/api/admin/members', { token });
  return data.items;
}
export function setMemberRole(memberId, role, token) {
  return apiFetch(`/api/admin/members/${memberId}/role`, { method: 'PUT', body: { role }, token });
}
export function getAdminMovie(movieId, token) {
  return apiFetch(`/api/admin/movies/${movieId}`, { token });
}
export function createMovie(fields, token) {
  return apiFetch('/api/admin/movies', { method: 'POST', body: fields, token });
}
export function updateMovie(movieId, fields, token) {
  return apiFetch(`/api/admin/movies/${movieId}`, { method: 'PUT', body: fields, token });
}
export function deleteMovie(movieId, token) {
  return apiFetch(`/api/admin/movies/${movieId}`, { method: 'DELETE', token }); // 204 ได้ null
}

export async function getReviews(movieId) {
  const data = await apiFetch(`/api/movies/${movieId}/reviews`);
  return data;
}

export function postReview(movieId, text, token) {
  return apiFetch(`/api/movies/${movieId}/reviews`, {
    method: 'POST',
    body: { text },
    token,
  });
}
