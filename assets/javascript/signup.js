const signupButton = document.getElementById("signup-btn")

async function signup(username, password) {
    const response = fetch("/api/signup", {
        method: "POST",
        body: JSON.stringify({ Username: username, Password: password })
    })

    const resp = await response

    if (!resp.ok) {
        console.log("Error")
        return
    }

    const data = await resp.json()
    if (data.error) {
        console.log(data.error)
    } else {
        console.log(data.status)
    }
}

document.addEventListener("DOMContentLoaded", async () => {
    //signup("itsbailey444", "TestPassword123")
})

signupButton.addEventListener("click", () => {
    signup("itsbailey44444", "TestPassword123")
})