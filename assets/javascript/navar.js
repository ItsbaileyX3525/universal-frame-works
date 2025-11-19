import { getAuthPromise } from '/assets/javascript/session.js'

const loginNavBarButton = document.getElementById("navbar-login")
const signupNavBarButton = document.getElementById("navbar-signup")
const accountNavBarButton = document.getElementById("navbar-account")

document.addEventListener("DOMContentLoaded", async () => {
    const auth = await getAuthPromise()
    console.log("User ID: " + auth)
    if (auth == null || !auth) {
        console.log("Not logged in")
    } else {
        accountNavBarButton.classList.remove("navbar-disabled")
        signupNavBarButton.classList.add("navbar-disabled")
        loginNavBarButton.classList.add("navbar-disabled")
        console.log("Logged in, " + auth)
    }
})