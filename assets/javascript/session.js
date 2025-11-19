export const sessionState = {
  isSessionValid: false,
  userID: 0
}

let authPromise = null

export function getAuthPromise() {
  if (!authPromise) {
    authPromise = checkAuth()
  }
  return authPromise
}

export async function checkAuth() {
    const resp = await fetch("/api/requireLogin", {
        method: "POST",
    })

    if (!resp.ok) {
        console.log("Error with fetch request")
        return
    }

    const data = await resp.json()

    console.log(data)
    if (data.status == "authenticated") {
        sessionState.isSessionValid = true
        sessionState.userID = data.userID
        return data.userID
    }
    return null
}

export async function logout() {
    const resp = await fetch("/api/logout", {
        method: "POST",
    })

    if (!resp.ok) {
        console.log("Error fetching data")
        return
    }

    const data = await resp.json()

    console.log(data)
    sessionState.isSessionValid = false
    sessionState.userID = 0
    return data
}
