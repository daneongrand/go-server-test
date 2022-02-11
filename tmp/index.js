const status = document.getElementById('connection-status')
const input = document.getElementById('input')
const btnAddComment = document.getElementById('btn-comment')
const commentList = document.getElementById('comment-list')

const socket = new WebSocket('ws://localhost:8080')

btnAddComment.addEventListener('click', () => {
    if (!input.value) return
    socket.send(input.value)
    input.value = ''
})

socket.onopen = function () {
    status.classList.remove('status-disconnect')
    status.classList.add('status-connect')
    status.innerHTML = 'connect'
    input.disabled = false
    btnAddComment.disabled = false
}

socket.onmessage = function (e) {
    commentList.appendChild(createComment(e.data))
}

socket.onclose = function () {
    status.classList.remove('status-connect')
    status.classList.add('status-disconnect')
    status.innerHTML = 'disconnect'
    input.disabled = true
    btnAddComment.disabled = true
}