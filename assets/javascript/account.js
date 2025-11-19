async function Logout() {
    const resp = await fetch("/api/logout", {
        method: "POST",
        body: JSON.stringify({
            "sessionID" : getCookie("session_id")
        })
    })

    if (!resp.ok) {
        console.log("Error fetching data")
        return
    }

    const data = await resp.json()

    console.log(data)
}

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
    await checkAuth()
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