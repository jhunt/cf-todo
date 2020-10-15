let Token = null
//const BASE = "http://127.0.0.1:3003"
const BASE = ""
const TokenKey = 'auth-token'

export default {
  url(rel) {
    return BASE+rel;
  },

  token(v) {
    Token = v
    if (Token) {
      localStorage.setItem(TokenKey, Token)
    } else {
      localStorage.removeItem(TokenKey)
    }
    return this
  },

  fetch(url, opts) {
    let token = localStorage.getItem(TokenKey)
    if (token) {
      if (!('headers' in opts)) {
        opts.headers = new Headers()
      }
      opts.headers.set('Authorization', 'Bearer '+token)
    }
    return fetch(BASE+url, opts).then(r => r.json())
  },

  fetchWithPayload(url, method, body, opts) {
    opts = Object.assign({}, opts || {})
    opts.method = method
    if (body) { opts.body = JSON.stringify(body); }
    return this.fetch(url, opts)

  },

  GET(url) {
    return this.fetch(url, { method: 'GET' })
  },

  DELETE(url) {
    return this.fetch(url, { method: 'DELETE' })
  },

  POST(url, body, opts) {
    return this.fetchWithPayload(url, 'POST', body, opts)

  },
  PUT(url, body, opts) {
    return this.fetchWithPayload(url, 'PUT', body, opts)

  },

  PATCH(url, body, opts) {
    return this.fetchWithPayload(url, 'PATCH', body, opts)
  }
}
