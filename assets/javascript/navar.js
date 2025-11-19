const loginNavBarButton = document.getElementById("navbar-login")
const signupNavBarButton = document.getElementById("navbar-signup")
const accountNavBarButton = document.getElementById("navbar-account")

async function checkAuth() {
    const resp = await fetch("/api/requireLogin", {
        method: "POST",
    })

    if (!resp.ok) {
        console.log("Error with fetch request")
        return
    }

    const data = await resp.json()

    console.log(data)
    return data.userID
}

document.addEventListener("DOMContentLoaded", async () => {
    const userID = await checkAuth()
    console.log("User ID: " + userID)
    if (userID == null || !userID) {
        console.log("Not logged in")
    } else {
        accountNavBarButton.classList.remove("navbar-disabled")
        signupNavBarButton.classList.add("navbar-disabled")
        loginNavBarButton.classList.add("navbar-disabled")
        console.log("Logged in, " + userID)
    }
})