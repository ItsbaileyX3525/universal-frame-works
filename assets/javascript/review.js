import { getAuthPromise, sessionState } from '/assets/javascript/session.js'

const selectionContainer = document.getElementById("selection-container")
const modal = document.getElementById("review-modal")
const closeBtn = document.getElementById("modal-close")
const submitBtn = document.getElementById("submit-btn")
const starRating = document.getElementById("star-rating")
const commentInput = document.getElementById("comment-input")
const authMessage = document.getElementById("auth-message")

let tabOpen = "none"
let onPage = 1
let currentItemUUID = null
let currentRating = 0

document.addEventListener("DOMContentLoaded", async () => {
    await getAuthPromise()
    registerDivs()
    setupModalListeners()
    updateButtonState()
    setupStarRating()
})

function setupStarRating() {
    const stars = starRating.querySelectorAll(".star")
    
    stars.forEach(star => {
        star.addEventListener("mouseover", () => {
            const value = parseInt(star.getAttribute("data-value"))
            highlightStars(value)
        })
        
        star.addEventListener("click", () => {
            const value = parseInt(star.getAttribute("data-value"))
            currentRating = value
            highlightStars(value)
            updateButtonState()
        })
    })
    
    starRating.addEventListener("mouseleave", () => {
        highlightStars(currentRating)
    })
}

function highlightStars(value) {
    const stars = starRating.querySelectorAll(".star")
    stars.forEach((star, index) => {
        if (index < value) {
            star.classList.add("active")
        } else {
            star.classList.remove("active")
        }
    })
    document.getElementById("rating-value").textContent = value
}

function setupModalListeners() {
    closeBtn.addEventListener("click", closeModal)
    
    modal.addEventListener("click", (e) => {
        if (e.target === modal) {
            closeModal()
        }
    })
    
    submitBtn.addEventListener("click", submitReviewAndComment)
}

function updateButtonState() {
    if (sessionState.isSessionValid) {
        submitBtn.disabled = false
        authMessage.style.display = "none"
        commentInput.disabled = false
    } else {
        submitBtn.disabled = true
        authMessage.style.display = "block"
        commentInput.disabled = true
    }
}

function handleForm(leForm) {
    console.log(leForm)
}

/*document.addEventListener("DOMContentLoaded", () => {
    let children = document.body.children
    for (let e of children) {
        if (e.tagName == "FORM") {
            e.addEventListener("submit", (ev) => {
                ev.preventDefault()
                handleForm(e)
            })
        }
        
    }
})*/

async function fetchItems(category, onPage) {
    const url = "/api/items?category=" + category + "&page=" + onPage
    const resp = await fetch(url)

    if (!resp.ok) {
        console.log("fetch failed")
        return ["fail"]
    }

    const data = await resp.json()

    if (data.status == "success") {
        var items = data.items
        if (typeof items != "object" /*Will always be object but yk*/ || items == null) {
            console.log("No items.")
            return ["fail"]
        }

        /*for (let e of items) { //Shud always be 8 max
            console.log("Item name: " + e.Name)
            console.log("Item ID: " + e.ID)
        }*/

        return items
    }
}

async function submitRating(uuid, score) {
    const resp = await fetch('/api/submitRating', {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            "uuid" : uuid,
            "rating" : score
        })
    })

    if (!resp.ok) {
        console.log("Something wrong with the request")
        return false
    }

    const data = await resp.json()

    if (data.status == "success") {
        console.log(data.message)
        return true
    } else {
        console.log(data)
        return false
    }
}

async function submitComment(uuid, content) {
    const resp = await fetch('/api/submitMessage', {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            "uuid" : uuid,
            "message" : content
        })
    })

    if (!resp.ok) {
        console.log("Something wrong with the request")
        return false
    }

    const data = await resp.json()

    if (data.status == "success") {
        console.log(data.message)
        return true
    } else {
        console.log(data)
        return false
    }
}

async function submitReviewAndComment() {
    if (!sessionState.isSessionValid) {
        console.log("User not authenticated")
        return
    }
    
    if (currentRating === 0) {
        console.log("Please select a star rating")
        return
    }
    
    const comment = commentInput.value.trim()
    const ratingSuccess = await submitRating(currentItemUUID, currentRating)
    let commentSuccess = true

    if (comment) {
        commentSuccess = await submitComment(currentItemUUID, comment)
    }

    if (ratingSuccess && commentSuccess) {
        console.log("Thank you for your review!")
        await fetchAndDisplayComments(currentItemUUID)
        closeModal()
    } else {
        console.log("Failed to submit rating. Please try again.")
    }
}

async function fetchAndDisplayComments(uuid) {
    const url = "/api/comments?uuid=" + uuid + "&page=1"
    const resp = await fetch(url)
    
    if (!resp.ok) {
        console.log("Failed to fetch comments")
        displayNoComments()
        return
    }
    
    const text = await resp.text()
    
    if (!text || text.trim() === "") {
        console.log("Empty response from comments endpoint")
        displayNoComments()
        return
    }
    
    const data = JSON.parse(text)
    
    if (data.status === "success" && data.comments && data.comments.length > 0) {
        displayComments(data.comments)
    } else {
        displayNoComments()
    }
}

function displayComments(comments) {
    const commentsList = document.getElementById("comments-list")
    commentsList.innerHTML = ""
    
    comments.forEach(comment => {
        const commentItem = document.createElement("div")
        commentItem.className = "comment-item"
        
        const author = document.createElement("div")
        author.className = "comment-author"
        author.textContent = comment.Username || "Anonymous"
        
        const rating = document.createElement("div")
        rating.className = "comment-rating"
        rating.textContent = "★".repeat(comment.Rating || 0)
        
        const text = document.createElement("div")
        text.className = "comment-text"
        text.textContent = comment.Comment || comment.Text || ""
        
        commentItem.appendChild(author)
        commentItem.appendChild(rating)
        commentItem.appendChild(text)
        commentsList.appendChild(commentItem)
    })
}

function displayNoComments() {
    const commentsList = document.getElementById("comments-list")
    commentsList.innerHTML = '<div class="comments-empty">No comments yet. Be the first!</div>'
}

async function openModal(uuid, itemData) {
    currentItemUUID = uuid
    currentRating = 0
    commentInput.value = ""
    
    document.getElementById("modal-title").textContent = itemData.Name || "Item"
    document.getElementById("modal-description").textContent = itemData.Description || "No description available"
    document.getElementById("modal-image").src = itemData.ImagePath || ""
    
    const ratingResp = await fetch("/api/averageRating?uuid=" + uuid)
    if (ratingResp.ok) {
        const ratingData = await ratingResp.json()
        if (ratingData.status === "success") {
            const avgRating = (Math.round(ratingData.avgRating * 100) / 100).toFixed(2)
            document.getElementById("modal-description").textContent = (itemData.Description || "No description available") + " | Avg Rating: " + avgRating + "/5"
        }
    }
    
    highlightStars(0)
    updateButtonState()
    
    await fetchAndDisplayComments(uuid)
    
    modal.classList.add("show")
    document.body.style.overflow = "hidden"
}

function closeModal() {
    modal.classList.remove("show")
    document.body.style.overflow = "auto"
    currentItemUUID = null
    currentRating = 0
}

function createCards(items) {
    for (let e of items) {
        const div = document.createElement("div")
        div.id = "frame-card"
        div.style.backgroundImage = `url('${e.ImagePath}')`
        div.setAttribute("uuid", e.ID)

        if (e.AvgRating !== undefined && e.AvgRating > 0) {
            const ratingDiv = document.createElement("div")
            ratingDiv.className = "card-rating"
            ratingDiv.textContent = "★ " + (Math.round(e.AvgRating * 100) / 100).toFixed(2) + "/5"
            div.appendChild(ratingDiv)
        }

        div.addEventListener("click", () => {
            const uuid = div.getAttribute("uuid")
            openModal(uuid, e)
        })
        selectionContainer.appendChild(div)
    }
}

async function switchTab(tabName) {
    selectionContainer.innerHTML = ""
    let items
    console.log(`Switching to ${tabName} tab`)
    tabOpen = tabName
    items = await fetchItems(tabOpen, onPage)
    createCards(items)
}

function registerDivs() {
    var divs = selectionContainer.children
    for (let e of divs) {
        if (e.tagName == "DIV") {
            e.addEventListener("click", (ev) => {
                ev.preventDefault() //Just in case
                var href = e.getAttribute("href")
                switchTab(href)
            })
        }
    }
}