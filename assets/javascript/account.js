import { getAuthPromise, logout } from '/assets/javascript/session.js'

const logoutButton = document.getElementById("logout-btn")
let username

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
        return
    }

    username = localStorage.getItem("username") || "None"
    var usernameAreas = document.querySelectorAll(".username-area")
    for (let e of usernameAreas) {
        e.innerText = username
    }
})