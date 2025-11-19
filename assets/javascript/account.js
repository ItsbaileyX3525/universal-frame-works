import { getAuthPromise, logout } from '/assets/javascript/session.js'

const logoutButton = document.getElementById("logout-btn")

logoutButton.addEventListener("click", async () => {
    await logout()
    window.location.href = "/login"
})

document.addEventListener("DOMContentLoaded", async () => {
    const auth = await getAuthPromise()
    console.log("User ID: " + auth)
    if (auth == null || !auth) {
        console.log("Not logged in")
        window.location.href = "/login"
    } else {
       console.log("Logged in, " + auth)
    }
})