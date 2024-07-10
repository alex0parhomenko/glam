const apiBaseUrl = 'http://localhost:8080';
let selectedProfileId = null;
let ws = null;


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
        <p>Avatar: ${profile.avatar} </p>
        <p>Name: ${profile.name}</p>
        <p>Posts: ${profile.posts}</p>
        <p>Liked Posts: ${profile.liked_posts}</p>
        <p>Notifications: ${profile.notifications}</p>
    `;
}

async function loadUserPosts(id) {
    if (!id) {
        alert('Please select a profile first.');
        return;
    }
    const response = await fetch(`${apiBaseUrl}/posts/${id}`);
    const posts = await response.json();
    console.log(posts.posts)
    const contentDiv = document.getElementById('content');
    contentDiv.innerHTML = `
        <h1>Posts</h1>
        <ul>
            ${posts.posts.map(post => `
                <li>
                    <p>User Id: ${post.user_id}</p>
                    <p>Post Id: ${post.id}</p>
                    <p>Content: ${post.content}</p>
                    <p>Likes: ${post.likes_count}</p>                    
                </li>
            `).join('')}
        </ul>
       
    `;
}

async function loadAllPosts(id) {
    if (!id) {
        alert('Please select a profile first.');
        return;
    }

    const response = await fetch(`${apiBaseUrl}/all_posts`);
    const posts = await response.json();
    console.log(posts.posts)
    const contentDiv = document.getElementById('content');
    contentDiv.innerHTML = `
        <h1>Posts</h1>
        <ul>
            ${posts.posts.map(post => `
                <li>
                    <p>User Id: ${post.user_id}</p>
                    <p>Post Id: ${post.id}</p>
                    <p>Content: ${post.content}</p>
                    <p>Likes: ${post.likes_count}</p>
                    
                    <button onclick="likePost('${post.id}', '${id}')">Like</button>
                </li>
            `).join('')}
        </ul>
    `;

}


async function createPPost(id) {
    const contentDiv = document.getElementById('content');
    contentDiv.innerHTML = `<input type="text" id="newPostContent" placeholder="Write a new post"/>
    <button onClick="createPost('${id}')">Create Post</button>`
}

async function loadUserLikedPosts(id) {
    if (!id) {
        alert('Please select a profile first.');
        return;
    }
    const response = await fetch(`${apiBaseUrl}/posts/liked/${id}`);
    const posts = await response.json();
    const contentDiv = document.getElementById('content');
    contentDiv.innerHTML = `
        <h1>Posts</h1>
        <ul>
            ${posts.liked_posts.map(post => `
                <li>
                    <p>User Id: ${post.user_id}</p>
                    <p>Post Id: ${post.id}</p>
                    <p>Content: ${post.content}</p>
                    <p>Likes: ${post.likes_count}</p>
                </li>
            `).join('')}
        </ul>
    `;
}

async function createPost(userId) {
    const content = document.getElementById('newPostContent').value;

    const result = await fetch(`${apiBaseUrl}/posts`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({user_id: userId, content: content})
    });
}

async function likePost(postId, userId) {
    await fetch(`${apiBaseUrl}/posts/like/${userId}/${postId}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        }
    });
    loadAllPosts(userId)
}



async function loadNotifications(userId) {
    if (!userId) {
        alert('Please select a profile first.');
        return;
    }
    const contentDiv = document.getElementById('content');
    if (ws == null) {
        ws = new WebSocket(`ws://localhost:8080/notifications/${userId}`);
        console.log("connect ws")
    }

    contentDiv.innerHTML = ``
    ws.onmessage = function (event) {
        const notification = JSON.parse(event.data);
        contentDiv.innerHTML += `
            <div>
                <p>Notification Type: ${notification.fullDocument.type}</p>
                <p>User ID: ${notification.fullDocument.user_id}</p>
                <p>Post ID: ${notification.fullDocument.post_id}</p>
            </div>
        `;
    };

    ws.onerror = function (error) {
        console.error('WebSocket error:', error);
    };
}

async function stopNotifications(userId) {
    if (!userId) {
        alert('Please select a profile first.');
        return;
    }
    if (ws) {
        ws.close()
        ws = null;
        console.log("Close ws")
    }
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
        body: JSON.stringify({name: name, avatar: avatar})
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
    stopNotifications(selectedProfileId)
    selectedProfileId = id;


    alert(`Selected profile ID: ${id}`);
}
