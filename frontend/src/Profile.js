import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useParams } from 'react-router-dom';

function Profile() {
    const { id } = useParams();
    const [profile, setProfile] = useState(null);

    useEffect(() => {
        axios.get(`/profile/${id}`)
            .then(response => setProfile(response.data))
            .catch(error => console.error('Error fetching profile:', error));
    }, [id]);

    if (!profile) {
        return <div>Loading...</div>;
    }

    return (
        <div>
            <h1>Profile</h1>
            <img src={profile.avatar} alt="Avatar" />
            <p>Name: {profile.name}</p>
            <p>Posts: {profile.posts.length}</p>
            <p>Liked Posts: {profile.liked_posts.length}</p>
            <p>Notifications: {profile.notifications.length}</p>
        </div>
    );
}

export default Profile;
