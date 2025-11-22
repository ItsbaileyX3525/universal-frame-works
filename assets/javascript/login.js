const loginForm = document.getElementById("login-form")

async function login(username, password) {
    const resp = await fetch('/api/login', {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            "username" : username,
            "password" : password,
        })
    })

    if (!resp.ok) {
        console.log("Login failed. Please check your credentials.")
        return
    }

    const data = await resp.json()

    if (data.status == "success") {
        if (data.username) {
            localStorage.setItem("username", data.username)
        }
        window.location.href = "/account"
    } else {
        console.log(data.message || "Login failed")
    }
}

loginForm.addEventListener("submit", (e) => {
    e.preventDefault()
    const username = document.getElementById("username").value.trim()
    const password = document.getElementById("password").value

    if (!username || !password) {
        console.log("Please fill in all fields")
        return
    }

    login(username, password)
})