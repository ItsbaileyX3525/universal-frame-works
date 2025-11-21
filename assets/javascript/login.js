const loginButton = document.getElementById("login-btn")

async function login(username, password, confirmPassword) {
    if (password !== confirmPassword) {
        console.log("Passwords don't match")
        return
    }

    const resp = await fetch('/api/login', {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            "username" : username,
            "password" : password,
            "confirmPassword" : confirmPassword,
        })
    })

    if (!resp.ok) {
        console.log("Something wrong")
        return
    }

    const data = await resp.json()

    console.log(data)
    if (data.status == "success") {
        console.log(data.message)
        if (data.username) {
            localStorage.setItem("username", data.username)
        }
        window.location.href = "/account"
    }
}

loginButton.addEventListener("click", () => {
    login("itsbailey444", "TestPassword123", "TestPassword123")
})