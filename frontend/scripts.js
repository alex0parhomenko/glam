const apiBaseUrl = 'http://localhost:8080';
let selectedProfileId = null;

async function loadProfile(id) {
    if (!id) {
        alert('Please select a profile first.');
        return;
    }
    const response = await fetch(`${apiBaseUrl}/profile/${id}`);
    const profile = await response.json();
    const contentDiv = document.getElementById('content');
    contentDiv.innerHTML = `
        <h1>Profile</h1>
        <img src="${profile.avatar}" alt="Avatar" />
        <p>Name: ${profile.name}</p>
        <p>Posts: ${profile.posts.length}</p>
        <p>Liked Posts: ${profile.liked_posts.length}</p>
        <p>Notifications: ${profile.notifications.length}</p>
    `;
}

async function loadPosts(id) {
    if (!id) {
        alert('Please select a profile first.');
        return;
    }
    const response = await fetch(`${apiBaseUrl}/posts`);
    const posts = await response.json();
    const contentDiv = document.getElementById('content');
    contentDiv.innerHTML = `
        <h1>Posts</h1>
        <ul>
            ${posts.map(post => `
                <li>
                    ${post.content}
                    <p>Likes: ${post.likes_count}</p>
                    <button onclick="likePost('${post.id}', ${id})">Like</button>
                </li>
            `).join('')}
        </ul>
        <input type="text" id="newPostContent" placeholder="Write a new post" />
        <button onclick="createPost(${id})">Create Post</button>
    `;
}

async function createPost(userId) {
    const content = document.getElementById('newPostContent').value;
    await fetch(`${apiBaseUrl}/posts`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ user_id: userId, content: content })
    });
    loadPosts(userId);
}

async function likePost(postId, userId) {
    await fetch(`${apiBaseUrl}/posts/${postId}/like`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ user_id: userId })
    });
    loadPosts(userId);
}

async function loadNotifications(userId) {
    if (!userId) {
        alert('Please select a profile first.');
        return;
    }
    const contentDiv = document.getElementById('content');
    const ws = new WebSocket(`ws://localhost:8080/notifications/${userId}`);

    ws.onmessage = function(event) {
        const notification = JSON.parse(event.data);
        contentDiv.innerHTML += `
            <div>
                <p>Notification Type: ${notification.type}</p>
                <p>Post ID: ${notification.post_id}</p>
            </div>
        `;
    };

    ws.onclose = function() {
        console.log('WebSocket connection closed');
    };

    ws.onerror = function(error) {
        console.error('WebSocket error:', error);
    };
}

async function showCreateProfileForm() {
    const contentDiv = document.getElementById('content');
    contentDiv.innerHTML = `
        <h1>Create Profile</h1>
        <input type="text" id="newProfileName" placeholder="Name" /><br />
        <input type="text" id="newProfileAvatar" placeholder="Avatar URL" /><br />
        <button onclick="createProfile()">Create Profile</button>
    `;
}

async function createProfile() {
    const name = document.getElementById('newProfileName').value;
    const avatar = document.getElementById('newProfileAvatar').value;
    const resp = await fetch(`${apiBaseUrl}/profile`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ name: name, avatar: avatar })
    });
    loadProfiles();
}

async function loadProfiles() {
    const response = await fetch(`${apiBaseUrl}/profiles`);
    const profiles = await response.json();
    const contentDiv = document.getElementById('content');
    contentDiv.innerHTML = `
        <h1>Profiles</h1>
        <ul>
            ${profiles.map(profile => `
                <li>
                    <a href="#" onclick="selectProfile('${profile.id}')">${profile.name}</a>
                </li>
            `).join('')}
        </ul>
    `;
}

function selectProfile(id) {
    selectedProfileId = id;
    alert(`Selected profile ID: ${id}`);
}
