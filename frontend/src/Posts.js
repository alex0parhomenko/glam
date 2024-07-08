import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useParams } from 'react-router-dom';

function Posts() {
    const { id } = useParams();
    const [posts, setPosts] = useState([]);
    const [newPost, setNewPost] = useState('');

    useEffect(() => {
        axios.get(`/posts/${id}`)
            .then(response => setPosts(response.data))
            .catch(error => console.error('Error fetching posts:', error));
    }, [id]);

    const handleCreatePost = () => {
        axios.post('/posts', { user_id: id, content: newPost })
            .then(response => setPosts([...posts, response.data]))
            .catch(error => console.error('Error creating post:', error));
    };

    return (
        <div>
            <h1>Posts</h1>
            <ul>
                {posts.map(post => (
                    <li key={post.id}>
                        {post.content}
                        <p>Likes: {post.likes_count}</p>
                    </li>
                ))}
            </ul>
            <input
                type="text"
                value={newPost}
                onChange={e => setNewPost(e.target.value)}
                placeholder="Write a new post"
            />
            <button onClick={handleCreatePost}>Create Post</button>
        </div>
    );
}

export default Posts;
