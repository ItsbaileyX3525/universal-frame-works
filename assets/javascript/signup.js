const signupForm = document.getElementById("signup-form")

async function signup(username, password) {
    const response = await fetch("/api/signup", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ Username: username, Password: password })
    })

    if (!response.ok) {
        console.log("Signup failed. Please try again.")
        return
    }

    const data = await response.json()
    if (data.error) {
        console.log(data.error)
    } else if (data.status == "success") {
        if (data.username) {
            localStorage.setItem("username", data.username)
        }
        window.location.href = "/account"
    } else {
        console.log(data.message || "Signup failed")
    }
}

signupForm.addEventListener("submit", (e) => {
    e.preventDefault()
    const username = document.getElementById("username").value.trim()
    const password = document.getElementById("password").value
    const confirmPassword = document.getElementById("confirm-password").value

    if (!username || !password || !confirmPassword) {
        console.log("Please fill in all fields")
        return
    }

    if (password !== confirmPassword) {
        console.log("Passwords do not match")
        return
    }

    if (password.length < 6) {
        console.log("Password must be at least 6 characters")
        return
    }

    signup(username, password)
})